package main

import (
	"github.com/ActiveState/tail"
	"io/ioutil"
	"flag"
	"os"
	"./logsend"
	"fmt"
	logpkg "log"
)

var (
	Log        = logpkg.New(os.Stdout, "", logpkg.Lmicroseconds)
	LogDir     = flag.String("log-dir", "./tmp", "log directories")
)

type LogFilesGroup struct {
	LogFiles []*LogFile
	Config *logsend.Group
}

func (lfg *LogFilesGroup) Tailing() {
	for _, log := range lfg.LogFiles {
		go func(lf *LogFile){
			for line := range lf.Tail.Lines {
				for _,rule := range lfg.Config.Rules {
					match,err := rule.MatchLogLine(line.Text)
					if err != nil { continue }
					err = rule.MakeJSON(match)
					if err != nil {
						Log.Println(err)
					}
				}
			}
		}(log)
	}
}

func (lfg *LogFilesGroup) AppendLog(log *LogFile) {
	lfg.LogFiles = append(lfg.LogFiles, log)
	return
}

type LogFile struct {
	Tail *tail.Tail
}

func NewLogFile(filename string) (*LogFile, error) {
	var err error
	logfile := &LogFile{}
	logfile.Tail, err = tail.TailFile(filename, tail.Config{Follow: true, ReOpen: true})
	return logfile, err
}

func GetFilesFromDir(dir string) ([]os.FileInfo, error) {
	files, err := ioutil.ReadDir(dir)
	return files, err
}

func main() {
	flag.Parse()
	Log.SetFlags(0)

	config, err := logsend.LoadConfig("config.json")
	if err != nil {
		Log.Fatalf("can't load config")
	}

	LogFilesGroups := make([]*LogFilesGroup, 0)
	for _, group := range config.Groups {
		lfg := &LogFilesGroup{Config: &group}
		LogFilesGroups = append(LogFilesGroups, lfg)
	}

	files, err := GetFilesFromDir(*LogDir)
	if err != nil {
		Log.Fatalln(err)
	}
	for _,f := range files {
		for _,lfg := range LogFilesGroups {
			fmt.Printf("%s/%s", *LogDir, f.Name())
			if lfg.Config.MatchLogLine(f.Name()) {
				Log.Println("match")
				logFile, err := NewLogFile(fmt.Sprintf("%s/%s", *LogDir, f.Name()))
				if err != nil {
					Log.Fatalln(err)
				}
				lfg.AppendLog(logFile)
			}
		}
	}

	for _,lfg := range LogFilesGroups {
		go lfg.Tailing()
	}


	select{}
}
