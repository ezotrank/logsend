package logsend

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type ConfigFile struct {
	Influxdb *InfluxDBConfig
	Groups   []*Group
}

func LoadConfig(fileName string) ([]*Group, error) {
	configFile := &ConfigFile{}
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	rawConfig, _ := ioutil.ReadAll(file)
	if err := json.Unmarshal(rawConfig, configFile); err != nil {
		return nil, err
	}

	if configFile.Influxdb != nil {
		InitInfluxdb(influxdbCh, configFile.Influxdb)
	}

	for _, group := range configFile.Groups {
		if err := group.loadRules(); err != nil {
			return nil, err
		}
	}
	return configFile.Groups, nil
}
