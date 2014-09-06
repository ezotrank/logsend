package logsend

import (
	"regexp"
)

type Group struct {
	Mask  *string `json:"mask"`
	Rules []*Rule `json:"rules"`
}

func (self *Group) loadRules() (err error) {
	for _, rule := range self.Rules {
		if err = rule.loadRegexp(); err != nil {
			return err
		}
		if rule.Influxdb != nil {
			influxsender := &InfluxdbSender{}
			influxsender.SetConfig(rule.Influxdb)
			rule.senders = append(rule.senders, influxsender)
		}

		if rule.Statsd != nil {
			statsdsender := &StatsdSender{}
			statsdsender.SetConfig(rule.Statsd)
			rule.senders = append(rule.senders, statsdsender)
		}
	}
	return
}

type Rule struct {
	Name     *string `json:"name"`
	Regexp   *string `json:"regexp"`
	regexp   *regexp.Regexp
	Influxdb interface{} `json:"influxdb"`
	Statsd   interface{} `json:"statsd"`
	senders  []Sender
}

func (self *Rule) loadRegexp() (err error) {
	self.regexp, err = regexp.Compile(*self.Regexp)
	return
}

func (rule *Rule) Match(line *string) interface{} {
	matches := rule.regexp.FindStringSubmatch(*line)

	if len(matches) == 0 {
		return nil
	}

	if len(matches) <= 1 {
		return true
	}

	// TODO: cache subexnames
	out := make(map[string]interface{})
	for i, value := range matches[1:] {
		key, value, err := prepareValue(rule.regexp.SubexpNames()[i+1], value)
		if err != nil {
			Conf.Logger.Printf("can't prepareValue with %+v and %+v have err %+v", rule.regexp.SubexpNames()[i+1], value, err)
			return nil
		}
		out[key] = value
	}
	if len(out) > 0 {
		return out
	}
	return nil
}

func (rule *Rule) send(data interface{}) {
	for _, sender := range rule.senders {
		sender.Send(data)
	}
}
