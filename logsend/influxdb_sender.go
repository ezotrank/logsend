package logsend

import (
	"flag"
	log "github.com/ezotrank/logger"
	influxdb "github.com/influxdb/influxdb/client"
	"strings"
)

// need remove this global variable on all senders
var (
	influxdbCh          = make(chan *influxdb.Series, 0)
	influxdbHost        = flag.String("influxdb-host", "localhost:8086", "influxdb host")
	influxdbUser        = flag.String("influxdb-user", "root", "influxdb user")
	influxdbPassword    = flag.String("influxdb-password", "root", "influxdb password")
	influxdbDatabase    = flag.String("influxdb-database", "", "influxdb database")
	influxdbUdp         = flag.Bool("influxdb-udp", false, "influxdb send via UDP")
	influxdbSendBuffer  = flag.Int("influxdb-send_buffer", 1, "influxdb UDP buffer size")
	influxdbSeriesName  = flag.String("influxdb-name", "", "influxdb series name")
	influxdbExtraFields = flag.String("influxdb-extra_fields", "", "Example: 'host,HOST service,www' ")
)

func init() {
	RegisterNewSender("influxdb", InitInfluxdb, NewInfluxdbSender)
}

func InitInfluxdb(conf interface{}) {
	config := &influxdb.ClientConfig{
		Host: conf.(map[string]interface{})["host"].(string),
	}

	if val, ok := conf.(map[string]interface{})["udp"]; ok {
		config.IsUDP = val.(bool)
	}

	if val, ok := conf.(map[string]interface{})["user"]; ok {
		config.Username = val.(string)
	} else if !config.IsUDP {
		log.Infoln("you must set `user`")
	}

	if val, ok := conf.(map[string]interface{})["password"]; ok {
		config.Password = val.(string)
	} else if !config.IsUDP {
		log.Infoln("you must set `password`")
	}

	if val, ok := conf.(map[string]interface{})["database"]; ok {
		config.Database = val.(string)
	} else if !config.IsUDP {
		log.Infoln("you must set `database`")
	}

	sendBuffer := 0
	if val, ok := conf.(map[string]interface{})["send_buffer"]; ok {
		sendBuffer = i2int(val)
	}
	client, err := influxdb.NewClient(config)
	if err != nil {
		log.Infof("can't create influxdb client err: %s\n", err)
	}
	client.DisableCompression()

	go func() {
		log.Infoln("Influxdb queue is starts")
		buf := make([]*influxdb.Series, 0)
		for series := range influxdbCh {
			if false {
				log.Infof("influxdb recieve series from influxdbCh: %+v", *series)
			}
			buf = append(buf, series)
			if len(buf) >= sendBuffer {
				if Conf.DryRun {
					log.Infof("influxdb dry-run send series: %+v", buf)
					buf = make([]*influxdb.Series, 0)
				} else {
					go writeSeries(client, config, buf)
				}
				buf = make([]*influxdb.Series, 0)
			}
		}
	}()
	return
}

func writeSeries(client *influxdb.Client, config *influxdb.ClientConfig, buf []*influxdb.Series) {
	log.Debugf("influxdb send series: %+v", buf)
	if config.IsUDP {
		client.WriteSeriesOverUDP(buf)
		return
	}
	if err := client.WriteSeries(buf); err != nil {
		log.Warnf("can't send series to db err: %s\n", err)
	}
	return
}

func NewInfluxdbSender() Sender {
	influxSender := &InfluxdbSender{
		sendCh: influxdbCh,
	}
	return Sender(influxSender)
}

type InfluxdbSender struct {
	name        string
	extraFields [][]*string
	sendCh      chan *influxdb.Series
}

// TODO: write test
func (self *InfluxdbSender) SetConfig(rawConfig interface{}) error {
	self.name = rawConfig.(map[string]interface{})["name"].(string)
	if extraFields, ok := rawConfig.(map[string]interface{})["extra_fields"]; ok {
		switch extraFields.(type) {
		case []interface{}:
			for _, pair := range extraFields.([]interface{}) {
				field := pair.([]interface{})[0].(string)
				value := pair.([]interface{})[1].(string)
				self.extraFields = append(self.extraFields, []*string{&field, &value})
			}
		case string:
			if extraFields.(string) != "" {
				for _, keys := range strings.Split(extraFields.(string), " ") {
					key := strings.Split(keys, ",")
					self.extraFields = append(self.extraFields, []*string{&key[0], &key[1]})
				}
			}
		}

	}
	return nil
}

func (self *InfluxdbSender) Name() string {
	return "influxdb"
}

func (self *InfluxdbSender) Send(data interface{}) {
	series := &influxdb.Series{
		Name: self.name,
	}
	columns := make([]string, 0)
	points := make([]interface{}, 0)
	switch data.(type) {
	case map[string]interface{}:
		for key, value := range data.(map[string]interface{}) {
			columns = append(columns, key)
			points = append(points, value)
		}
	}

	for _, extraField := range self.extraFields {
		if val, err := ExtendValue(extraField[1]); err == nil {
			columns = append(columns, *extraField[0])
			points = append(points, val)
		}
	}
	series.Columns = columns
	series.Points = [][]interface{}{points}

	self.sendCh <- series
}
