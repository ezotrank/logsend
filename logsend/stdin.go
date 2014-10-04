package logsend

import (
	"bufio"
	"flag"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func ProcessStdin() error {
	rules := make([]*Rule, 0)
	if rawConfig["config"].(flag.Value).String() != "" {
		groups, err := LoadConfigFromFile(rawConfig["config"].(flag.Value).String())
		if err != nil {
			Conf.Logger.Fatalf("can't load config %+v", err)
		}
		for _, group := range groups {
			for _, rule := range group.Rules {
				rules = append(rules, rule)
			}
		}
	} else {
		// TODO: move this to separate method
		matchSender := regexp.MustCompile(`(\w+)-host`)
		var sender Sender
		for key, val := range rawConfig {
			match := matchSender.FindStringSubmatch(key)
			if len(match) == 0 || val.(flag.Value).String() == "" {
				continue
			}
			if register, ok := Conf.registeredSenders[match[1]]; ok {
				conf := make(map[string]interface{})
				for key, val := range rawConfig {
					newKey := key
					if ok, _ := regexp.MatchString(match[1], key); ok {
						newKey = strings.Split(key, match[1]+"-")[1]
					}
					switch val.(flag.Value).String() {
					default:
						conf[newKey] = interface{}(val.(flag.Value).String())
					case "true", "false":
						b, err := strconv.ParseBool(val.(flag.Value).String())
						if err != nil {
							Conf.Logger.Fatalln(err)
						}
						conf[newKey] = interface{}(b)
					}
				}
				register.Init(conf)
				sender = register.get()
				sender.SetConfig(conf)
				break
			}
		}
		rule := &Rule{
			regexp:  regexp.MustCompile(rawConfig["regex"].(flag.Value).String()),
			senders: []Sender{sender},
		}
		rules = append(rules, rule)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			break
		}
		checkLineRules(&line, rules)
	}
	return nil
}
