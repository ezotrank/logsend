package logsend

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func LoadConfig(fileName string) ([]*Group, error) {
	groups := make([]*Group, 0)
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	brules, _ := ioutil.ReadAll(file)
	if err := json.Unmarshal(brules, &groups); err != nil {
		return nil, err
	}
	for _, group := range groups {
		if err := group.loadRulesRegexp(); err != nil {
			return nil, err
		}
	}
	return groups, nil
}
