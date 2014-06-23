package logsend

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"regexp"
)

func LoadConfig(fileName string) (config *Config, err error) {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	brules, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(brules, &config)
	return
}

type Config struct {
	Groups []Group `json:"groups"`
}

type Group struct {
	Mask string `json:"mask"`
	Rules []Rule `json:"rules"`
	regex *regexp.Regexp
}

func (group *Group) Match(line string) bool {
	if group.regex == nil {
		regex,err := regexp.Compile(group.Mask)
		if err != nil {
			log.Printf("group match err %+v", err)
			return false
		}
		group.regex = regex
	}
	return group.regex.MatchString(line)
}

type Rule struct {
	Name string `json:"name"`
	Regexp string `json:"regexp"`
	Columns [][]string `json:"columns"`
	regex *regexp.Regexp
}

func GetValues(svals []string, rculumns [][]string) (columns []string, points []interface{}, err error) {
	for index,val := range svals {
		columns = append(columns, rculumns[index][0])
		var ival interface{}
		if len(rculumns[index]) == 1 {
			points = append(points, val)
		} else if len(rculumns[index]) == 2 {
			ival, err := LeadToType(val, rculumns[index][1])
			if err != nil {
				log.Fatalf("GetValues %+v", err)
			}
			points = append(points, ival)
		} else {
			ival, err = ConvertToPoint(val, rculumns[index][2])
			points = append(points, ival)
		}
	}
	return
}

func (rule *Rule) Match(line string) (matches []string) {
	if rule.regex == nil {
		regex,err := regexp.Compile(rule.Regexp)
		if err != nil {
			log.Printf("rule match err %+v", err)
		}
		rule.regex = regex
	}
	if !rule.regex.MatchString(line) {
		return
	}
	matches = rule.regex.FindStringSubmatch(line)
	matches = append(matches[1:])
	return
}

