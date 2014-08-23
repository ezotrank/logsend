package logsend

import (
	influxdb "github.com/influxdb/influxdb/client"
	"net/http"
)

var (
	influxdbCh = make(chan *influxdb.Series)
)

func init() {
	config := &influxdb.ClientConfig{
		Host:       Conf.DBHost,
		Username:   Conf.DBUser,
		Password:   Conf.DBPassword,
		Database:   Conf.DBName,
		HttpClient: http.DefaultClient,
		IsUDP:      Conf.UDP,
	}
	client, err := influxdb.NewClient(config)
	if err != nil {
		log.Fatalln(err)
	}
	client.DisableCompression()

	go func() {
		buf := make([]*influxdb.Series, 0)
		for series := range influxdbCh {
			debug("go func", *series)
			buf = append(buf, series)
			if len(buf) >= Conf.SendBuffer {
				if Conf.UDP {
					go client.WriteSeriesOverUDP(buf)
				} else {
					go client.WriteSeries(buf)
				}
				// clean buffer
				buf = make([]*influxdb.Series, 0)
			}
		}
	}()
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

// func getSeries(rule *Rule, columns []string, values []interface{}) *influxdb.Series {
// 	series := &influxdb.Series{
// 		Name:    *rule.Name,
// 		Columns: columns,
// 		Points:  [][]interface{}{values},
// 	}
// 	return series
// }

// func SendSeries(series *influxdb.Series) {
// 	influxCh <- series
// 	return
// }
