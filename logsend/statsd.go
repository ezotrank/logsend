package logsend

import (
	"github.com/quipo/statsd"
	"strings"
	"time"
)

var (
	statsdCh = make(chan map[string]map[string]int64)
)

func init() {
	go func() {
		prefix := "test."
		statsdclient := statsd.NewStatsdClient("10.70.120.213:8125", prefix)
		statsdclient.CreateSocket()
		interval, _ := time.ParseDuration("1s")
		stats := statsd.NewStatsdBuffer(interval, statsdclient)
		defer stats.Close()
		for data := range statsdCh {
			for op, values := range data {
				for key, val := range values {
					switch op {
					case "increment":
						debug("send incr", key, val)
						stats.Incr(key, val)
					case "timing":
						debug("send timing", key, val)
						stats.Timing(key, val)
					}
				}

			}
		}
	}()
}

type StatsdSender struct {
	timing    [][]string
	increment []string
	gauge     []string
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
	return nil
}

func (self *StatsdSender) Name() string {
	return "StatsdSender"
}

func interfaceToInt64(i interface{}) (val int64, err error) {
	switch i.(type) {
	case float64:
		val = int64(i.(float64))
	}
	return
}

func replaceKey(str string) string {
	hostname, _ := getHostname()
	return strings.Replace(str, `%HOST%`, hostname.(string), 1)
}

func (self *StatsdSender) Send(data interface{}) {
	for _, name := range self.increment {
		name = replaceKey(name)
		statsdCh <- map[string]map[string]int64{"increment": {name: 1}}
	}

	for _, values := range self.timing {
		if val, ok := data.(map[string]interface{})[values[1]]; ok {
			intval, err := interfaceToInt64(val)
			if err != nil {
				log.Fatalln(err)
			}
			key := replaceKey(values[0])
			statsdCh <- map[string]map[string]int64{"timing": {key: intval}}
		}
	}
}
