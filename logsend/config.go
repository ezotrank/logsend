package logsend

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"
)

func LoadConfig(fileName string) (*Config, error) {
	config := &Config{}
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	brules, _ := ioutil.ReadAll(file)
	if err := json.Unmarshal(brules, &config); err != nil {
		return nil, err
	}
	for _, group := range config.Groups {
		group.loadRulesRegexp()
	}
	return config, nil
}

type Config struct {
	Groups []*Group `json:"groups"`
}

type Group struct {
	Mask  string  `json:"mask"`
	Rules []*Rule `json:"rules"`
	regex *regexp.Regexp
}

func (self *Group) loadRulesRegexp() {
	for _, rule := range self.Rules {
		if err := rule.loadRegexp(); err != nil {
			log.Panicf("can't load regexp %+v", err)
		}
	}
}

func (group *Group) Match(line string) bool {
	if group.regex == nil {
		regex, err := regexp.Compile(group.Mask)
		if err != nil {
			log.Printf("group match err %+v", err)
			return false
		}
		group.regex = regex
	}
	return group.regex.MatchString(line)
}

type Rule struct {
	Name    string     `json:"name"`
	Regexp  string     `json:"regexp"`
	Columns [][]string `json:"columns"`
	regex   *regexp.Regexp
}

func (self *Rule) loadRegexp() (err error) {
	self.regex, err = regexp.Compile(self.Regexp)
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

func (rule *Rule) Match(line string) (matches []string) {
	if !rule.regex.MatchString(line) {
		return
	}
	if matches = rule.regex.FindStringSubmatch(line); len(matches) != 0 {
		return matches[1:]
	}
	return
}
