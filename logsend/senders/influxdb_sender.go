package logsend

import (
	"flag"
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
	influxdbSendBuffer  = flag.Int("influxdb-send_buffer", 8, "influxdb UDP buffer size")
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

	config.IsUDP = conf.(map[string]interface{})["udp"].(bool)
	config.Username = conf.(map[string]interface{})["user"].(string)
	config.Password = conf.(map[string]interface{})["password"].(string)
	config.Database = conf.(map[string]interface{})["database"].(string)

	sendBuffer := 0
	if val, ok := conf.(map[string]interface{})["send_buffer"]; ok {
		sendBuffer = i2int(val)
	}
	client, err := influxdb.NewClient(config)
	if err != nil {
		Conf.Logger.Fatalln(err)
	}
	client.DisableCompression()

	go func() {
		Conf.Logger.Println("Influxdb queue is starts")
		buf := make([]*influxdb.Series, 0)
		for series := range influxdbCh {
			debug("go func", *series)
			buf = append(buf, series)
			if len(buf) >= sendBuffer {
				if Conf.DryRun {
					continue
				}
				if config.IsUDP {
					go client.WriteSeriesOverUDP(buf)
				} else {
					go writeSeries(client, buf)
				}
				// clean buffer
				buf = make([]*influxdb.Series, 0)
			}
		}
	}()
	return
}

func writeSeries(client *influxdb.Client, buf []*influxdb.Series) {
	err := client.WriteSeries(buf)
	if err != nil {
		Conf.Logger.Printf("influxdb can't write series %+v", err)
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
