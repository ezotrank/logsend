package logsend

import (
	"github.com/influxdb/influxdb-go"
	"net/http"
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
	return client, err
}

func GetSeries(rule *Rule, columns []string, values []interface{}) (seriesT []*influxdb.Series) {
	series := &influxdb.Series{}
	out := [][]interface{}{values}
	series.Name = rule.Name
	series.Columns = columns
	series.Points = out
	seriesT = append(seriesT, series)
	return
}

func SendSeries(series []*influxdb.Series, client *influxdb.Client) {
	debug("write series", series[0].Name, series[0].Columns, series[0].Points)
	go client.WriteSeries(series)
}
