package logsend

import (
	"github.com/ActiveState/tail"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func WatchFiles(dir, configFile string) {
	groups, err := LoadConfig(configFile)
	if err != nil {
		Conf.Logger.Fatalf("can't load config %+v", err)
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		Conf.Logger.Fatalf("can't read config dir: %+v", err)
	}
	assignFiles(files, groups)
	if Conf.ContinueWatch {
		go continueWatch(&dir, groups)
	}
	select {}
}

func assignFiles(files []os.FileInfo, groups []*Group) {
	for _, group := range groups {
		files, err := getFilesByGroup(files, group)
		if err != nil {
			Conf.Logger.Fatalf("can't get file %+v", err)
		}
		for _, file := range files {
			go file.tail(group)
		}
	}
}

func continueWatch(dir *string, groups []*Group) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		Conf.Logger.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() {
					files := make([]os.FileInfo, 0)
					file, err := os.Stat(ev.Name)
					if err != nil {
						Conf.Logger.Printf("can't get file %+v", err)
					}
					files = append(files, file)
					assignFiles(files, groups)
				}
				// debug(ev.IsCreate())
			case err := <-watcher.Error:
				Conf.Logger.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(*dir)
	if err != nil {
		Conf.Logger.Fatal(err)
	}

	<-done

	/* ... do stuff ... */
	watcher.Close()
}

func getFilesByGroup(allFiles []os.FileInfo, group *Group) ([]*File, error) {
	files := make([]*File, 0)
	regex := regexp.MustCompile(*group.Mask)
	for _, f := range allFiles {
		if !regex.MatchString(f.Name()) {
			continue
		}
		filepath := filepath.Join(Conf.WatchDir, f.Name())
		file, err := NewFile(filepath)
		if err != nil {
			return files, err
		}
		files = append(files, file)
	}
	return files, nil
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
	Conf.Logger.Printf("start tailing %+v", self.Tail.Filename)
	for line := range self.Tail.Lines {
		checkLine(&line.Text, group.Rules)
	}
}

func checkLine(line *string, rules []*Rule) error {
	for _, rule := range rules {
		match := rule.Match(line)
		if match == nil {
			continue
		}
		if !Conf.DryRun {
			rule.send(match)
		}
	}
	return nil
}
