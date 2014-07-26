package logsend

import (
	"github.com/ActiveState/tail"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

func WatchLogs(logDir, configFile string) {
	config, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("can't load config %+v", err)
	}

	logsScopes := make([]*LogScope, 0)

	for _, group := range config.Groups {
		lsc := NewLogScope(group)
		logsScopes = append(logsScopes, lsc)
	}

	AssociatedLogPerFile(logDir, &logsScopes)
	err = NewDBClient()
	if err != nil {
		log.Fatalf("NewDBClient %+v", err)
	}

	for _, lsc := range logsScopes {
		go lsc.Tailing()
	}

	select {}
}

func NewLogScope(group Group) *LogScope {
	lsc := &LogScope{}
	lsc.ConfigGroup = &group
	return lsc
}

type LogScope struct {
	Logs        []*Log
	ConfigGroup *Group
}

func checkLine(line *string, rules *[]Rule) error {
	for _, rule := range *rules {
		match := rule.Match(*line)
		if len(match) == 0 {
			continue
		}
		colums, values, err := GetValues(match, rule.Columns)
		if err != nil {
			log.Printf("GetValues err %+v", err)
			return err
		}
		series := GetSeries(&rule, colums, values)
		SendSeries(series)
	}
	return nil
}

func (lsc *LogScope) Tailing() {
	for _, logf := range lsc.Logs {
		go func(lg *Log) {
			log.Printf("start tailing %+v", lg.Tail.Filename)
			for line := range lg.Tail.Lines {
				go checkLine(&line.Text, &lsc.ConfigGroup.Rules)
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
	seekInfo := &tail.SeekInfo{Offset: 0, Whence: 2}
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
	for _, f := range dfiles {
		if !regex.MatchString(f.Name()) {
			continue
		}
		files = append(files, filepath.Join(dir, f.Name()))
	}
	return
}

func AssociatedLogPerFile(dir string, logsScopes *[]*LogScope) {
	for _, lsc := range *logsScopes {
		files, err := GetFilesByName(dir, lsc.ConfigGroup.Mask)
		if err != nil {
			log.Fatalf("GetFilsByName %+v", err)
		}
		for _, file := range files {
			logf, err := NewLogFile(file)
			if err != nil {
				log.Fatalf("NewLogFile %+v", err)
			}
			lsc.Append(logf)
		}
	}
}
