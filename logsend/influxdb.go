package logsend

import (
	influxdb "github.com/influxdb/influxdb/client"
	"net/http"
)

type InfluxDBConfig struct {
	Host       string
	User       string
	Password   string
	Database   string
	Udp        bool
	SendBuffer int
}

var (
	influxdbCh = make(chan *influxdb.Series, 0)
)

func InitInfluxdb(ch chan *influxdb.Series, conf *InfluxDBConfig) error {
	config := &influxdb.ClientConfig{
		Host:       conf.Host,
		Username:   conf.User,
		Password:   conf.Password,
		Database:   conf.Database,
		IsUDP:      conf.Udp,
		HttpClient: http.DefaultClient,
	}
	client, err := influxdb.NewClient(config)
	if err != nil {
		Conf.Logger.Fatalln(err)
	}
	client.DisableCompression()

	go func() {
		Conf.Logger.Println("Influxdb queue is starts")
		buf := make([]*influxdb.Series, 0)
		for series := range ch {
			debug("go func", *series)
			buf = append(buf, series)
			if len(buf) >= conf.SendBuffer {
				if conf.Udp {
					if !Conf.DryRun {
						go client.WriteSeriesOverUDP(buf)
					}
				} else {
					go client.WriteSeries(buf)
				}
				// clean buffer
				buf = make([]*influxdb.Series, 0)
			}
		}
	}()
	return nil
}

type InfluxdbSender struct {
	name        string
	extraFields [][]*string
}

func (self *InfluxdbSender) SetConfig(rawConfig interface{}) error {
	self.name = rawConfig.(map[string]interface{})["name"].(string)
	if extraFields, ok := rawConfig.(map[string]interface{})["extra_fields"]; ok {
		for _, pair := range extraFields.([]interface{}) {
			field := pair.([]interface{})[0].(string)
			value := pair.([]interface{})[1].(string)
			self.extraFields = append(self.extraFields, []*string{&field, &value})
		}
	}
	return nil
}

func (self *InfluxdbSender) Name() string {
	return "InfluxdbSender"
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
		if val, err := extendValue(extraField[1]); err == nil {
			columns = append(columns, *extraField[0])
			points = append(points, val)
		}
	}
	series.Columns = columns
	series.Points = [][]interface{}{points}

	influxdbCh <- series
}
