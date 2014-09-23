package logsend

import (
	"bufio"
	"flag"
	"os"
	"strconv"
)

func ProcessStdin() error {
	rules := make([]*Rule, 0)
	if rawConfig["config"].(flag.Value).String() != "" {
		groups, err := LoadConfig(rawConfig["config"].(flag.Value).String())
		if err != nil {
			Conf.Logger.Fatalf("can't load config %+v", err)
		}
		for _, group := range groups {
			for _, rule := range group.Rules {
				rules = append(rules, rule)
			}
		}
	} else {
		influxdbConfg := &InfluxDBConfig{
			Host:     rawConfig["influx-host"].(flag.Value).String(),
			User:     rawConfig["influx-user"].(flag.Value).String(),
			Password: rawConfig["influx-password"].(flag.Value).String(),
			Database: rawConfig["influx-dbname"].(flag.Value).String(),
		}
		if rawConfig["influx-udp"].(flag.Value).String() == "true" {
			influxdbConfg.Udp = true
		} else {
			influxdbConfg.Udp = false
		}
		sendBuffer, _ := strconv.ParseInt(rawConfig["influx-udp-buffer"].(flag.Value).String(), 0, 16)
		influxdbConfg.SendBuffer = int(sendBuffer)
		if err := InitInfluxdb(influxdbCh, influxdbConfg); err != nil {
			Conf.Logger.Fatalln(err)
		}
		seriesName := rawConfig["influx-series-name"].(flag.Value).String()
		regex := rawConfig["regex"].(flag.Value).String()
		sender := &InfluxdbSender{name: seriesName}

		rule := &Rule{
			Name:   &seriesName,
			Regexp: &regex,
		}
		rule.senders = []Sender{sender}
		if err := rule.loadRegexp(); err != nil {
			Conf.Logger.Fatalln(err)
		}
		rules = append(rules, rule)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			// You may check here if err == io.EOF
			break
		}
		checkLineRules(&line, rules)
	}
	return nil
}
