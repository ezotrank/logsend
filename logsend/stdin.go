package logsend

import (
	"bufio"
	"os"
)

func ProcessStdin(conf *InfluxDBConfig, regex, seriesName string) error {
	if err := InitInfluxdb(influxdbCh, conf); err != nil {
		Conf.Logger.Fatalln(err)
	}

	sender := &InfluxdbSender{name: seriesName}

	rule := &Rule{
		Name:   &seriesName,
		Regexp: &regex,
	}
	rule.senders = []Sender{sender}
	if err := rule.loadRegexp(); err != nil {
		Conf.Logger.Fatalln(err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			// You may check here if err == io.EOF
			break
		}

		checkLineRule(&line, rule)
	}
	return nil
}
