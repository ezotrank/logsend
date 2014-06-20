package logsend

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"regexp"
	logpkg "log"
	"strconv"
	"fmt"
	"github.com/influxdb/influxdb-go"
	"net/http"
)

var (
	Log        = logpkg.New(os.Stdout, "", logpkg.Lmicroseconds)
)

var DBClientConfig = &influxdb.ClientConfig{
		Host:       "localhost:8086",
		Username:   "root",
		Password:   "root",
		Database:   "test1",
		HttpClient: http.DefaultClient,
	}

	var DBClient, err = influxdb.NewClient(DBClientConfig)

type Config struct {
	Groups []Group `json:"groups"`
}

func LoadConfig(fileName string) (c *Config, err error) {
	config := &Config{}
	var file *os.File
	file, err = os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return
	}
	defer file.Close()
	brules, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(brules, config)
	if err != nil {
		return
	}
	c = config
	return
}

type Group struct {
	Mask string `json:"mask"`
	Rules []Rule `json:"rules"`
	regex *regexp.Regexp
}

func (group *Group) MatchLogLine(str string) bool {
	if group.regex == nil {
		Log.Println("assign regexp")
		r,err := regexp.Compile(group.Mask)
		group.regex = r
		if err != nil {
			Log.Printf("%+v, %+v", err, str)
		}
	}
	if group.regex.MatchString(str) {
		return true
	}
	return false
}

type Rule struct {
	Name string `json:"name"`
	PointersRegexp string `json:"points_regexp"`
	Types []string `json:types`
	Columns []string `json:columns`
	regex *regexp.Regexp
}

func leadToType(obj string, objt string) (result interface{}, err error) {
	switch objt {
	case "int":
		result, err = strconv.ParseInt(obj, 0, 64)
	}
	return
}

func (rule *Rule) MatchLogLine(str string) (matches []interface{}, err error) {
	if rule.regex == nil {
		Log.Println("assign regexp")
		r,err := regexp.Compile(rule.PointersRegexp)
		rule.regex = r
		if err != nil {
			Log.Printf("%+v, %+v", err, str)
		}
	}
	out := make([]interface{},0)
	if rule.regex.MatchString(str) {
		finded := rule.regex.FindStringSubmatch(str)
		finded = append(finded[1:])
		for i,obj := range finded {
			val, err := leadToType(obj, rule.Types[i])
			if err != nil {
				Log.Fatalf("MatchLogLine %+v %+v", val, err)
			}
			out = append(out, val)
		}
		matches = out
		Log.Println(out)
		return
	}
	err = fmt.Errorf("cant find")
	return
}

func (rule *Rule) MakeJSON(str []interface{}) (err error) {
	series := &influxdb.Series{}
	out := [][]interface{}{str}
	series.Name = rule.Name
	series.Columns = rule.Columns
	series.Points = out
	seriesT := []*influxdb.Series{series}
	go DBClient.WriteSeries(seriesT)
	return
}