package logsend

import (
	"encoding/json"
	"flag"
	"fmt"
	log "github.com/ezotrank/logger"
	client "github.com/influxdb/influxdb/client"
	"net/url"
	"strings"
	"time"
)

// need remove this global variables on all senders
var (
	influxCh              = make(chan *SendPoint, 0)
	influxHost            = flag.String("influx-host", "http://localhost:8086", "client host")
	influxUser            = flag.String("influx-user", "root", "client user")
	influxPassword        = flag.String("influx-password", "root", "client password")
	influxDatabase        = flag.String("influx-database", "", "client database")
	influxSendBuffer      = flag.Int("influx-buffer", 1, "buffer size")
	influxSeriesName      = flag.String("influx-name", "", "client series name")
	influxRetentionPolicy = flag.String("influx-retention-policy", "default", "retention policy")
	influxExtraFields     = flag.String("influx-extra_fields", "", "Example: 'host,HOST service,www'")
)

func init() {
	RegisterNewSender("influxdb", InitInfluxdb, NewInfluxdbSender)
}

func InitInfluxdb(conf interface{}) {
	config := &client.Config{}

	if val, ok := conf.(map[string]interface{})["host"]; ok {
		surl := val.(string)
		if surl[:3] != "http" {
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

	go func() {
		log.Infoln("Influxdb queue is starts")
		buf := make(map[string][]SendPoint, 0)
		for lastPoint := range influxCh {
			key := lastPoint.sender.key()
			if val, ok := buf[lastPoint.sender.key()]; ok {
				buf[key] = append(val, *lastPoint)
			} else {
				buf[key] = []SendPoint{*lastPoint}
			}
			if len(buf[key]) >= lastPoint.sender.bufLen {
				go writeSeries(influxClient, buf[key])
				buf[key] = make([]SendPoint, 0)
			}
		}
	}()
	return
}

func writeSeries(influxClient *client.Client, sendPoints []SendPoint) {
	if len(sendPoints) < 1 {
		log.Errorln("send points less then one")
		return
	}
	firstSendPoint := sendPoints[0]
	points := make([]client.Point, 0)
	for _, sp := range sendPoints {
		points = append(points, *sp.point)
	}
	batchPoints := client.BatchPoints{
		Points:          points,
		Database:        firstSendPoint.sender.database,
		RetentionPolicy: firstSendPoint.sender.retentionPolicy,
		Tags:            firstSendPoint.sender.tags,
		Precision:       firstSendPoint.sender.precision,
	}
	if log.AvailableForLevel(log.DEBUGLV) {
		b, err := json.Marshal(&batchPoints)
		if err != nil {
			panic(err)
		}
		log.Debugf("query for send to db: %v", string(b))
	}
	if Conf.DryRun {
		log.Infof("client dry-run send series: count %d", len(batchPoints.Points))
		return
	}
	resp, err := influxClient.Write(batchPoints)
	if err != nil {
		log.Warnf("can't sent points to database err: %s\n", err)
		log.Warnf("trying again")
		for _, sendPoint := range sendPoints {
			sendPoint.sender.sendCh <- &sendPoint
		}
	} else {
		log.Infof("series sent %+v", len(batchPoints.Points))
	}
	log.Debugln(resp)
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
	tags            map[string]string
	precision       string
	extraFields     [][]*string
	bufLen          int
	sendCh          chan *SendPoint
}

func (self *InfluxdbSender) Name() string {
	return "influxdb"
}

func (sender *InfluxdbSender) key() string {
	return sender.database + sender.name + fmt.Sprintln(sender.tags)
}

type SendPoint struct {
	sender *InfluxdbSender
	point  *client.Point
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

	if val, ok := rawConfig.(map[string]interface{})["retention_policy"]; ok {
		self.retentionPolicy = val.(string)
	} else {
		self.retentionPolicy = *influxRetentionPolicy
	}

	if val, ok := rawConfig.(map[string]interface{})["buffer"]; ok {
		self.bufLen = int(val.(float64))
	} else {
		self.bufLen = *influxSendBuffer
	}

	self.tags = make(map[string]string, 0)
	if val, ok := rawConfig.(map[string]interface{})["tags"]; ok {
		for k, v := range val.(map[string]interface{}) {
			sVal := v.(string)
			if val, err := ExtendValue(&sVal); err == nil {
				self.tags[k] = val.(string)
			}
		}
	}
	if _, ok := self.tags["host"]; !ok {
		sVal := "HOST"
		host, _ := ExtendValue(&sVal)
		self.tags["host"] = host.(string)
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

func (self *InfluxdbSender) Send(data interface{}) {
	fields := make(map[string]interface{}, 0)

	for k, v := range data.(map[string]interface{}) {
		fields[k] = v
	}

	for _, extraField := range self.extraFields {
		if val, err := ExtendValue(extraField[1]); err == nil {
			fields[*extraField[0]] = val
		}
	}

	point := &client.Point{
		Measurement: self.name,
		Time:        time.Now().UTC(),
		Fields:      fields,
		Tags:        self.tags,
	}
	sendPoint := &SendPoint{sender: self, point: point}
	self.sendCh <- sendPoint
}
