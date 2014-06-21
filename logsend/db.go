package logsend

import (
	"net/http"
	"github.com/influxdb/influxdb-go"
)

func NewDBClient(host, user, password, database string) (*influxdb.Client, error) {
	config :=  &influxdb.ClientConfig{
		Host:       host,
		Username:   user,
		Password:   password,
		Database:   database,
		HttpClient: http.DefaultClient,
	}
	client, err := influxdb.NewClient(config)
	return client, err
}

func GetSeries(rule *Rule, data []interface{}) (seriesT []*influxdb.Series) {
	series := &influxdb.Series{}
	out := [][]interface{}{data}
	series.Name = rule.Name
	series.Columns = rule.Columns
	series.Points = out
	seriesT = append(seriesT, series)
	return
}

func SendSeries(series []*influxdb.Series, client *influxdb.Client) {
	debug("write series")
	go client.WriteSeries(series)
}