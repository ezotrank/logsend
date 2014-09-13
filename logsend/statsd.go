package logsend

import (
	"github.com/quipo/statsd"
	// "strings"
	"time"
)

type StatsdConfig struct {
	Host     string
	Prefix   string
	Interval string
}

var (
	statsdCh = make(chan *map[string]map[string]int64, 0)
)

func InitStatsd(ch chan *map[string]map[string]int64, conf *StatsdConfig) error {
	go func() {
		statsdclient := statsd.NewStatsdClient(conf.Host, conf.Prefix)
		statsdclient.CreateSocket()
		interval, err := time.ParseDuration(conf.Interval)
		if err != nil {
			Conf.Logger.Fatalf("can't parse interval %+v", err)
		}
		stats := statsd.NewStatsdBuffer(interval, statsdclient)
		defer stats.Close()
		for data := range statsdCh {
			for op, values := range *data {
				for key, val := range values {
					switch op {
					case "increment":
						debug("send incr", key, val)
						stats.Incr(key, val)
					case "timing":
						debug("send timing", key, val)
						stats.Timing(key, val)
					case "gauge":
						debug("send gauge", key, val)
						stats.Gauge(key, val)
					}
				}

			}
		}
	}()
	return nil
}

type StatsdSender struct {
	timing    [][]string
	gauge     [][]string
	increment []string
}

func (self *StatsdSender) SetConfig(rawConfig interface{}) error {
	if val, ok := rawConfig.(map[string]interface{})["timing"]; ok {
		for _, vals := range val.([]interface{}) {
			self.timing = append(self.timing, []string{vals.([]interface{})[0].(string), vals.([]interface{})[1].(string)})
		}
		debug(self.timing)
	}

	if val, ok := rawConfig.(map[string]interface{})["increment"]; ok {
		for _, vals := range val.([]interface{}) {
			self.increment = append(self.increment, vals.(string))
		}
		debug(self.increment)
	}

	if val, ok := rawConfig.(map[string]interface{})["gauge"]; ok {
		for _, vals := range val.([]interface{}) {
			self.gauge = append(self.gauge, []string{vals.([]interface{})[0].(string), vals.([]interface{})[1].(string)})
		}
		debug(self.gauge)
	}
	return nil
}

func (self *StatsdSender) Name() string {
	return "StatsdSender"
}

func interfaceToInt64(i interface{}) (val int64, err error) {
	switch i.(type) {
	default:
		val = i.(int64)
	case float64:
		val = int64(i.(float64))
	}
	return
}

func (self *StatsdSender) Send(data interface{}) {
	for _, name := range self.increment {
		statsdCh <- &map[string]map[string]int64{"increment": {name: 1}}
	}

	for _, values := range self.timing {
		if val, ok := data.(map[string]interface{})[values[1]]; ok {
			intval, err := interfaceToInt64(val)
			if err != nil {
				Conf.Logger.Printf("can't convert to int64 %+v", err)
			}
			if err != nil {
				Conf.Logger.Fatalln(err)
			}
			statsdCh <- &map[string]map[string]int64{"timing": {values[0]: intval}}
		}
	}

	for _, values := range self.gauge {
		if val, ok := data.(map[string]interface{})[values[1]]; ok {
			intval, err := interfaceToInt64(val)
			if err != nil {
				Conf.Logger.Printf("can't convert to int64 %+v", err)
			}
			if err != nil {
				Conf.Logger.Fatalln(err)
			}
			statsdCh <- &map[string]map[string]int64{"gauge": {values[0]: intval}}
		}
	}
}
