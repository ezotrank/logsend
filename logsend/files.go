package logsend

import (
	"github.com/ActiveState/tail"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

func WatchFiles(dir, configFile string) {
	groups, err := LoadConfig(configFile)
	if err != nil {
		log.Fatalf("can't load config %+v", err)
	}

	for _, group := range groups {
		files, err := getFilesByGroup(&dir, group)
		if err != nil {
			log.Fatalf("can't get file %+v", err)
		}
		for _, file := range files {
			go file.tail(group)
		}
	}
	select {}
}

func getFilesByGroup(dir *string, group *Group) ([]*File, error) {
	files := make([]*File, 0)
	dfiles, err := ioutil.ReadDir(*dir)
	if err != nil {
		return files, err
	}
	regex := regexp.MustCompile(*group.Mask)
	for _, f := range dfiles {
		if !regex.MatchString(f.Name()) {
			continue
		}
		filepath := filepath.Join(*dir, f.Name())
		file, err := NewFile(filepath)
		if err != nil {
			return files, err
		}
		files = append(files, file)
	}
	return files, err
}

func NewFile(fpath string) (*File, error) {
	seekInfo := &tail.SeekInfo{Offset: 0, Whence: 2}
	file := &File{}
	var err error
	file.Tail, err = tail.TailFile(fpath, tail.Config{Follow: true, ReOpen: true, Location: seekInfo})
	return file, err
}

type File struct {
	Tail *tail.Tail
}

func (self *File) tail(group *Group) {
	log.Printf("start tailing %+v", self.Tail.Filename)
	for line := range self.Tail.Lines {
		go checkLine(&line.Text, group.Rules)
	}
}

func checkLine(line *string, rules []*Rule) error {
	for _, rule := range rules {
		match := rule.Match(line)
		if len(match) == 0 {
			continue
		}
		colums, values, err := GetValues(match, rule.Columns)
		if err != nil {
			log.Printf("GetValues err %+v", err)
			return err
		}
		series := GetSeries(rule, colums, values)
		SendSeries(series)
	}
	return nil
}
