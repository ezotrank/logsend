package logsend

import (
	"github.com/influxdb/influxdb-go"
	"net/http"
)

var (
	SenderCh = make(chan *influxdb.Series)
)

func NewDBClient(host, user, password, database string) (*influxdb.Client, error) {
	config := &influxdb.ClientConfig{
		Host:       host,
		Username:   user,
		Password:   password,
		Database:   database,
		HttpClient: http.DefaultClient,
	}
	client, err := influxdb.NewClient(config)
	go func() {
		buf := make([]*influxdb.Series, 0)
		for series := range SenderCh {
			buf = append(buf, series)
			if len(buf) >= SendBuffer {
				debug("send series: ", buf)
				go client.WriteSeries(buf)
				buf = make([]*influxdb.Series, 0)
			}

		}
	}()
	return client, err
}

func GetSeries(rule *Rule, columns []string, values []interface{}) *influxdb.Series {
	series := influxdb.Series{}
	series.Name = rule.Name
	series.Columns = columns
	points := [][]interface{}{values}
	series.Points = points
	return &series
}

func SendSeries(series *influxdb.Series, client *influxdb.Client) {
	SenderCh <- series
	return
}
