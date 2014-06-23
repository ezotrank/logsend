package logsend

import (
	"path/filepath"
	"github.com/ActiveState/tail"
	"github.com/influxdb/influxdb-go"
	"regexp"
	"io/ioutil"
)

func NewLogScope(group *Group) *LogScope {
	lsc := &LogScope{}
	lsc.ConfigGroup = group
	return lsc
}

type LogScope struct {
	Logs []*Log
	ConfigGroup *Group
}

func (lsc *LogScope) Tailing(client *influxdb.Client) {
	for _, logf := range lsc.Logs {
		go func(lg *Log){
			log.Printf("start tailing %+v", lg.Tail.Filename)
			for line := range lg.Tail.Lines {
				for _,rule := range lsc.ConfigGroup.Rules {
					match := rule.Match(line.Text)
					if len(match) == 0 {
						continue
					}
					colums, values, err := GetValues(match, rule.Columns)
					if err != nil {
						log.Printf("GetValues err %+v", err)
					}
					series := GetSeries(&rule, colums, values)
					go SendSeries(series, client)
				}
			}
		}(logf)
	}
}

func (ls *LogScope) Append(log *Log) {
	ls.Logs = append(ls.Logs, log)
}

type Log struct {
	Tail *tail.Tail
}

func NewLogFile(filename string) (*Log, error) {
	var err error
	seekInfo := &tail.SeekInfo{0, 2}
	logfile := &Log{}
	logfile.Tail, err = tail.TailFile(filename, tail.Config{Follow: true, ReOpen: true, Location: seekInfo})
	return logfile, err
}

func GetFilesByName(dir, mask string) (files []string, err error) {
	dfiles, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	regex := regexp.MustCompile(mask)
	for _,f := range dfiles {
		if !regex.MatchString(f.Name()) {
			continue
		}
		files = append(files, filepath.Join(dir, f.Name()))
	}
	return
}

func AssociatedLogPerFile(dir string, logsScopes []*LogScope) {
	for _, lsc := range logsScopes {
		files,err := GetFilesByName(dir, lsc.ConfigGroup.Mask)
		if err != nil {
			log.Fatalf("GetFilsByName %+v", err)
		}
		for _,file := range files {
			logf, err := NewLogFile(file)
			if err != nil {
				log.Fatalf("NewLogFile %+v", err)
			}
			lsc.Append(logf)
		}
	}
}