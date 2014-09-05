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
		log.Fatalln(err)
	}
	client.DisableCompression()

	go func() {
		log.Println("Influxdb queue is starts")
		buf := make([]*influxdb.Series, 0)
		for series := range ch {
			debug("go func", *series)
			buf = append(buf, series)
			if len(buf) >= conf.SendBuffer {
				if conf.Udp {
					go client.WriteSeriesOverUDP(buf)
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
	name string
}

func (self *InfluxdbSender) SetConfig(rawConfig interface{}) error {
	self.name = rawConfig.(map[string]interface{})["name"].(string)
	return nil
}

func (self *InfluxdbSender) Name() string {
	return "InfluxdbSender"
}

func (self *InfluxdbSender) Send(data interface{}) {
	series := &influxdb.Series{
		Name: self.name,
	}
	switch data.(type) {
	case map[string]interface{}:
		columns := make([]string, 0)
		points := make([]interface{}, 0)
		for key, value := range data.(map[string]interface{}) {
			columns = append(columns, key)
			points = append(points, value)
		}
		series.Columns = columns
		series.Points = [][]interface{}{points}
	}
	influxdbCh <- series
}
