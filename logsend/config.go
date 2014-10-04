package logsend

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"regexp"
)

func LoadRawConfig(f *flag.Flag) {
	rawConfig[f.Name] = f.Value
}

func LoadConfigFromFile(fileName string) (groups []*Group, err error) {
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	rawConfig, err := ioutil.ReadAll(file)
	if err != nil {
		Conf.Logger.Fatalln(err)
	}
	return LoadConfig(rawConfig)
}

// TODO: test
func LoadConfig(rawConfig []byte) (groups []*Group, err error) {
	config := make(map[string]interface{})
	if err := json.Unmarshal(rawConfig, &config); err != nil {
		return nil, err
	}

	for sender, register := range Conf.registeredSenders {
		if val, ok := config[sender]; ok {
			register.Init(val)
		}
	}

	for _, groupConfig := range config["groups"].([]interface{}) {
		group := &Group{}
		if group.Mask, err = regexp.Compile(groupConfig.(map[string]interface{})["mask"].(string)); err != nil {
			Conf.Logger.Fatalln(err)
		}
		for _, groupRule := range groupConfig.(map[string]interface{})["rules"].([]interface{}) {
			// regex, err := regexp.Compile()
			if err != nil {
				Conf.Logger.Fatalln(err)
			}
			senders := make([]Sender, 0)
			for senderName, register := range Conf.registeredSenders {
				// not load rules to not initilized senders
				if register.initialized != true {
					continue
				}
				if val, ok := groupRule.(map[string]interface{})[senderName].(interface{}); ok {
					sender := register.get()
					if err = sender.SetConfig(val); err != nil {
						Conf.Logger.Fatalln(err)
					}
					senders = append(senders, sender)
				}
			}
			rule, err := NewRule(groupRule.(map[string]interface{})["regexp"].(string))
			if err != nil {
				panic(err)
			}
			rule.senders = senders
			group.Rules = append(group.Rules, rule)
		}
		groups = append(groups, group)
	}
	return
}
