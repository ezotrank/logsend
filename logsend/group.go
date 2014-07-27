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

func GetValues(svals []string, rculumns [][]string) (columns []string, points []interface{}, err error) {
	for index, col := range rculumns {
		columns = append(columns, col[0])
		if index <= len(svals)-1 {
			if len(col) == 1 {
				points = append(points, svals[index])
			} else if len(col) == 2 {
				ival, err := LeadToType(svals[index], col[1])
				if err != nil {
					log.Fatalf("GetValues %+v", err)
				}
				points = append(points, ival)
			} else {
				ival, err := ConvertToPoint(svals[index], col[2])
				if err != nil {
					log.Fatalf("GetValues %+v", err)
				}
				points = append(points, ival)
			}
		} else {
			if len(col) == 1 {
				points = append(points, "")
			} else {
				ival, err := GetValue(col[1])
				if err != nil {
					log.Fatalf("GetValues %+v", err)
				}
				points = append(points, ival)
			}
		}
	}
	return
}

func (rule *Rule) Match(line *string) (matches []string) {
	if matches = rule.regexp.FindStringSubmatch(*line); len(matches) != 0 {
		return matches[1:]
	}
	return
}
