package logsend

import (
	"flag"
	log "github.com/ezotrank/logger"
	client "github.com/influxdb/influxdb/client"
	"net/url"
	"strings"
)

// need remove this global variable on all senders
var (
	influxCh              = make(chan client.Write, 0)
	influxHost            = flag.String("influx-host", "http://localhost:8086", "client host")
	influxUser            = flag.String("influx-user", "root", "client user")
	influxPassword        = flag.String("influx-password", "root", "client password")
	influxDatabase        = flag.String("influx-database", "", "client database")
	influxSendBuffer      = flag.Int("influx-send_buffer", 1, "buffer size")
	influxSeriesName      = flag.String("influx-name", "", "client series name")
	influxRetentionPolicy = "180d"
	influxExtraFields     = flag.String("influx-extra_fields", "", "Example: 'host,HOST service,www' ")
)

func init() {
	RegisterNewSender("influxdb", InitInfluxdb, NewInfluxdbSender)
}

func InitInfluxdb(conf interface{}) {
	config := &client.Config{}

	if val, ok := conf.(map[string]interface{})["host"]; ok {
		surl := val.(string)
		if surl[:7] != "http://" || surl[:7] != "htts://" {
			surl = "http://" + surl
		}
		u, err := url.Parse(surl)
		if err != nil {
			log.Fatalln(err)
		}
		config.URL = *u
	} else {
		log.Fatalln("you must set `host`")
	}

	if val, ok := conf.(map[string]interface{})["user"]; ok {
		config.Username = val.(string)
	} else {
		log.Fatalln("you must set `user`")
	}

	if val, ok := conf.(map[string]interface{})["password"]; ok {
		config.Password = val.(string)
	} else {
		log.Fatalln("you must set `password`")
	}

	influxClient, err := client.NewClient(*config)
	if err != nil {
		log.Fatalln("can't create client client err: %s\n", err)
	}

	if val, ok := conf.(map[string]interface{})["database"]; ok {
		*influxDatabase = val.(string)
	}

	if val, ok := conf.(map[string]interface{})["send_buffer"]; ok {
		*influxSendBuffer = i2int(val)
	}

	go func() {
		log.Infoln("Influxdb queue is starts")
		for write := range influxCh {
			if Conf.DryRun {
				log.Infof("client dry-run send series: %+v", write)
			} else {
				go writeSeries(influxClient, write)
			}
		}
	}()
	return
}

func writeSeries(influxClient *client.Client, write client.Write) {
	log.Debugf("client send series: %+v", write)
	if results, err := influxClient.Write(write); err != nil {
		log.Warnf("can't send series to db err: %s\n", err)
		log.Warnln(results)
	}
	return
}

func NewInfluxdbSender() Sender {
	influxSender := &InfluxdbSender{
		sendCh: influxCh,
	}
	return Sender(influxSender)
}

type InfluxdbSender struct {
	name            string
	database        string
	retentionPolicy string
	extraFields     [][]*string
	sendCh          chan client.Write
}

// TODO: write test
func (self *InfluxdbSender) SetConfig(rawConfig interface{}) error {
	if val, ok := rawConfig.(map[string]interface{})["name"]; ok {
		self.name = val.(string)
	} else {
		log.Fatalf("clientsender set config field %s should be present", "name")
	}

	if val, ok := rawConfig.(map[string]interface{})["database"]; ok {
		self.database = val.(string)
	} else {
		self.database = *influxDatabase
	}

	if val, ok := rawConfig.(map[string]interface{})["retantion_policy"]; ok {
		self.retentionPolicy = val.(string)
	} else {
		self.retentionPolicy = influxRetentionPolicy
	}

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
	write := client.Write{
		Database:        self.database,
		RetentionPolicy: self.retentionPolicy,
	}
	point := client.Point{
		Name:   self.name,
		Values: data.(map[string]interface{}),
	}

	for _, extraField := range self.extraFields {
		if val, err := ExtendValue(extraField[1]); err == nil {
			point.Values[*extraField[0]] = val
		}
	}
	write.Points = []client.Point{point}
	self.sendCh <- write
}
