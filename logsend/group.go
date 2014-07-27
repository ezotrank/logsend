package logsend

import (
	"regexp"
)

type Group struct {
	Mask  *string `json:"mask"`
	Rules []*Rule `json:"rules"`
}

func (self *Group) loadRulesRegexp() (err error) {
	for _, rule := range self.Rules {
		if err = rule.loadRegexp(); err != nil {
			return err
		}
	}
	return
}

type Rule struct {
	Name    *string    `json:"name"`
	Regexp  *string    `json:"regexp"`
	Columns [][]string `json:"columns"`
	regexp  *regexp.Regexp
}

func (self *Rule) loadRegexp() (err error) {
	self.regexp, err = regexp.Compile(*self.Regexp)
	return
}

func (rule *Rule) Match(line *string) (matches []string) {
	if matches = rule.regexp.FindStringSubmatch(*line); len(matches) != 0 {
		return matches[1:]
	}
	return
}
