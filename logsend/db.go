package logsend

import (
	"github.com/ezotrank/influxdb-go"
	"net/http"
)

var (
	SenderCh = make(chan *influxdb.Series)
)

func NewDBClient() error {
	config := &influxdb.ClientConfig{
		Host:       Conf.DBHost,
		Username:   Conf.DBUser,
		Password:   Conf.DBPassword,
		Database:   Conf.DBName,
		HttpClient: http.DefaultClient,
		IsUDP:      Conf.UDP,
	}
	client, err := influxdb.NewClient(config)
	client.DisableCompression()
	go func() {
		buf := make([]*influxdb.Series, 0)
		for series := range SenderCh {
			buf = append(buf, series)
			if len(buf) >= SendBuffer {
				debug("buf: ", buf)
				if Conf.UDP {
					debug("send series over udp")
					go client.WriteSeriesOverUDP(buf)
				} else {
					debug("send series over http")
					go client.WriteSeries(buf)
				}

				buf = make([]*influxdb.Series, 0)
			}

		}
	}()
	return err
}

func GetSeries(rule *Rule, columns []string, values []interface{}) *influxdb.Series {
	series := influxdb.Series{}
	series.Name = rule.Name
	series.Columns = columns
	points := [][]interface{}{values}
	series.Points = points
	return &series
}

func SendSeries(series *influxdb.Series) {
	SenderCh <- series
	return
}
