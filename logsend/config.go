package logsend

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"regexp"
	"strconv"
)

func LoadConfig(fileName string) (config *Config, err error) {
	var file *os.File
	var tmpConfig *Config
	file, err = os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	brules, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(brules, &tmpConfig)
	config = tmpConfig
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
	PointersRegexp string `json:"points_regexp"`
	Types []string `json:types`
	Columns []string `json:columns`
	regex *regexp.Regexp
}

func leadToType(val string, valType string) (result interface{}, err error) {
	switch valType {
	case "int":
		result, err = strconv.ParseInt(val, 0, 64)
	}
	return
}

func (rule *Rule) Match(line string) (matches []interface{}) {
	if rule.regex == nil {
		regex,err := regexp.Compile(rule.PointersRegexp)
		if err != nil {
			log.Printf("rule match err %+v", err)
		}
		rule.regex = regex
	}
	if !rule.regex.MatchString(line) {
		return
	}
	submatches := rule.regex.FindStringSubmatch(line)
	submatches = append(submatches[1:])
	for i,obj := range submatches {
		val, err := leadToType(obj, rule.Types[i])
		if err != nil {
			log.Fatalf("MatchLogLine %+v %+v", val, err)
		}
		matches = append(matches, val)
	}
	return
}

