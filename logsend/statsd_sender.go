package logsend

import (
	"flag"
	"fmt"
	"github.com/quipo/statsd"
	"strings"
	"time"
)

var (
	statsdCh        = make(chan *map[string]map[string]int64, 0)
	statsdHost      = flag.String("statsd-host", "", "statsd host")
	statsdPrefix    = flag.String("statsd-prefix", "", "statsd prefix")
	statsdInterval  = flag.String("statsd-interval", "1s", "statsd prefix")
	statsdIncrement = flag.String("statsd-increment", "", "Example: 'test1incr,test2incr'")
	statsdTiming    = flag.String("statsd-timing", "", "Example: 'testtiming1,val testtiming2,val'")
	statsdGauge     = flag.String("statsd-gauge", "", "Example: 'gauge1,val gauge2,val'")
)

func init() {
	RegisterNewSender("statsd", InitStatsd, NewStatsdSender)
}

func InitStatsd(conf interface{}) {
	host := conf.(map[string]interface{})["host"].(string)
	prefix := ""
	_interval := "1s"
	if val, ok := conf.(map[string]interface{})["prefix"].(string); ok {
		prefix = val
	}
	if val, ok := conf.(map[string]interface{})["interval"].(string); ok {
		_interval = val
	}
	statsdclient := statsd.NewStatsdClient(host, prefix)

	statsdclient.CreateSocket()
	interval, err := time.ParseDuration(_interval)
	if err != nil {
		fmt.Printf("can't parse interval %+v", err)
	}
	stats := statsd.NewStatsdBuffer(interval, statsdclient)
	go func() {
		defer stats.Close()
		fmt.Println("Statsd queue is starts")
		for data := range statsdCh {
			for op, values := range *data {
				for key, val := range values {
					switch op {
					case "increment":
						go stats.Incr(key, val)
					case "timing":
						go stats.Timing(key, val)
					case "gauge":
						go stats.Gauge(key, val)
					}
				}

			}
		}

	}()
	return
}

func NewStatsdSender() Sender {
	sender := &StatsdSender{
		sendCh: statsdCh,
	}
	return Sender(sender)
}

type StatsdSender struct {
	timing    [][]string
	gauge     [][]string
	increment []string
	sendCh    chan *map[string]map[string]int64
}

func (self *StatsdSender) SetConfig(rawConfig interface{}) error {
	if val, ok := rawConfig.(map[string]interface{})["timing"]; ok {
		switch val.(type) {
		default:
			for _, vals := range val.([]interface{}) {
				self.timing = append(self.timing, []string{vals.([]interface{})[0].(string), vals.([]interface{})[1].(string)})
			}
		case string:
			if val.(string) != "" {
				for _, keys := range strings.Split(val.(string), " ") {
					key := strings.Split(keys, ",")
					self.timing = append(self.timing, []string{key[0], key[1]})
				}
			}
		}
	}

	if val, ok := rawConfig.(map[string]interface{})["increment"]; ok {
		switch val.(type) {
		default:
			for _, vals := range val.([]interface{}) {
				self.increment = append(self.increment, vals.(string))
			}
		case string:
			if val.(string) != "" {
				for _, key := range strings.Split(val.(string), ",") {
					self.increment = append(self.increment, key)
				}
			}
		}
	}

	if val, ok := rawConfig.(map[string]interface{})["gauge"]; ok {
		switch val.(type) {
		default:
			for _, vals := range val.([]interface{}) {
				self.gauge = append(self.gauge, []string{vals.([]interface{})[0].(string), vals.([]interface{})[1].(string)})
			}
		case string:
			if val.(string) != "" {
				for _, keys := range strings.Split(val.(string), " ") {
					key := strings.Split(keys, ",")
					self.gauge = append(self.gauge, []string{key[0], key[1]})
				}
			}
		}
	}
	return nil
}

func (self *StatsdSender) Name() string {
	return "statsd"
}

func interfaceToInt64(i interface{}) (val int64, err error) {
	switch i.(type) {
	case int64:
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
				fmt.Printf("can't convert to int64 %+v", err)
			}
			if err != nil {
				fmt.Println(err)
			}
			statsdCh <- &map[string]map[string]int64{"timing": {values[0]: intval}}
		}
	}

	for _, values := range self.gauge {
		if val, ok := data.(map[string]interface{})[values[1]]; ok {
			intval, err := interfaceToInt64(val)
			if err != nil {
				fmt.Printf("can't convert to int64 %+v", err)
			}
			if err != nil {
				fmt.Println(err)
			}
			self.sendCh <- &map[string]map[string]int64{"gauge": {values[0]: intval}}
		}
	}
}
