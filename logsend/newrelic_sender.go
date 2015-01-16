package logsend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"text/template"
	"time"
)

var (
	newRelicConfig = &NewRelicConfig{}
)

func init() {
	RegisterNewSender("newrelic", InitNewRelic, NewNewRelicSender)
}

type NewRelicConfig struct {
	Host     string
	Key      string
	Name     string
	Duration uint
	Version  string
}

type NewRelicReport struct {
	Agent      NewRelicReportAgent       `json:"agent"`
	Components []NewRelicReportComponent `json:"components"`
}

type NewRelicReportAgent struct {
	Host    string `json:"host"`
	Pid     int    `json:"pid"`
	Version string `json:"version"`
}

type NewRelicReportComponent struct {
	Name     string                        `json:"name"`
	Guid     string                        `json:"guid"`
	Duration uint                          `json:"duration"`
	Metrics  map[string]map[string]float64 `json:"metrics"`
}

func InitNewRelic(conf interface{}) {
	if val, ok := conf.(map[string]interface{})["host"]; !ok {
		panic("`host` must be present in config")
	} else {
		newRelicConfig.Host = val.(string)
	}

	if val, ok := conf.(map[string]interface{})["key"]; !ok {
		panic("`key` must be present in config")
	} else {
		newRelicConfig.Key = val.(string)
	}

	if val, ok := conf.(map[string]interface{})["name"]; !ok {
		panic("`name` must be present in config")
	} else {
		newRelicConfig.Name = val.(string)
	}

	newRelicConfig.Version = "0.1"
	return
}

func NewNewRelicSender() Sender {
	sender := &NewRelicSender{
		metrics: map[string][]float64{},
		config:  newRelicConfig,
	}
	return Sender(sender)
}

type NewRelicSender struct {
	sync.Mutex
	tmpl     *template.Template
	duration uint
	config   *NewRelicConfig
	metrics  map[string][]float64
}

func maxMetrics(metrics []float64) (max float64) {
	for _, i := range metrics {
		if i > max {
			max = i
		}
	}
	return
}

func minMetrics(metrics []float64) (min float64) {
	for i, val := range metrics {
		if i == 0 {
			min = val
			continue
		}
		if val < min {
			min = val
		}
	}
	return
}

func sumMetrics(metrics []float64) (sum float64) {
	for _, i := range metrics {
		sum = sum + i
	}
	return
}

func (self *NewRelicSender) SetConfig(rawConfig interface{}) error {
	var err error
	if val, ok := rawConfig.(map[string]interface{})["tmpl"]; !ok {
		panic("newrelic SetConfig `tmpl` not present in config")
	} else {
		if self.tmpl, err = template.New("query").Parse(val.(string)); err != nil {
			fmt.Printf("newrelic can't parse template %+v err: %s", val, err)
		}
	}
	if val, ok := rawConfig.(map[string]interface{})["duration"]; !ok {
		panic("newrelic SetConfig `duration` not present in config")
	} else {
		self.duration = uint(val.(float64))
	}

	go func() {
		for {
			time.Sleep(time.Duration(self.duration) * time.Second)
			if len(self.metrics) < 1 {
				continue
			}
			self.Lock()
			metrics := self.metrics
			self.metrics = map[string][]float64{}
			self.Unlock()
			go self.send(metrics)
		}
	}()
	return nil
}

func (self *NewRelicSender) send(metrics map[string][]float64) {
	report := &NewRelicReport{
		Agent: NewRelicReportAgent{
			Host:    self.config.Host,
			Pid:     0,
			Version: self.config.Version,
		},
		Components: []NewRelicReportComponent{
			NewRelicReportComponent{
				Name:     self.config.Name,
				Guid:     "logsend",
				Duration: self.duration,
			},
		},
	}
	for name, val := range metrics {
		metric := map[string]float64{
			"min":   minMetrics(val),
			"max":   maxMetrics(val),
			"total": sumMetrics(val),
			"count": float64(len(val)),
		}
		if report.Components[0].Metrics == nil {
			report.Components[0].Metrics = make(map[string]map[string]float64, 0)
		}
		report.Components[0].Metrics[name] = metric
	}
	breport, err := json.Marshal(report)
	if err != nil {
		fmt.Printf("newrelic can't marshal report %v", err)
		return
	}
	if Conf.DryRun {
		return
	}

	req, err := http.NewRequest("POST", "https://platform-api.newrelic.com/platform/v1/metrics", bytes.NewBuffer(breport))
	if err != nil {
		fmt.Printf("newrelic can't make request %v", err)
		return
	}
	req.Header.Set("X-License-Key", self.config.Key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("newrelic can't send payload to NewRelic %v", err)
	}
	defer resp.Body.Close()
}

func (self *NewRelicSender) Name() string {
	return "NewRelic"
}

func (self *NewRelicSender) Send(data interface{}) {
	timeGroupNames := make([]string, 0)
	switch data.(type) {
	case map[string]interface{}:
		for key, val := range data.(map[string]interface{}) {
			switch val.(type) {
			case float64, int:
				timeGroupNames = append(timeGroupNames, key)
			}
		}
	}

	if len(timeGroupNames) > 0 {
		for _, name := range timeGroupNames {
			buf := new(bytes.Buffer)
			if err := self.tmpl.Execute(buf, map[string]string{"NAME": name}); err != nil {
				fmt.Println("newrelic template error ", err, data)
				return
			}
			str := buf.String()
			if val, ok := self.metrics[str]; ok {
				self.metrics[str] = append(val, data.(map[string]interface{})[name].(float64))
			} else {
				self.metrics[str] = []float64{data.(map[string]interface{})[name].(float64)}
			}
		}

	} else {
		buf := new(bytes.Buffer)
		if err := self.tmpl.Execute(buf, data); err != nil {
			fmt.Println("newrelic template error ", err, data)
			return
		}
		str := buf.String()
		if val, ok := self.metrics[str]; ok {
			self.metrics[str] = append(val, 1)
		} else {
			self.metrics[str] = []float64{1}
		}

	}
	return
}
