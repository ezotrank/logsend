package logsend

import (
	"github.com/ActiveState/tail"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"os"
	"path/filepath"
)

func WatchFiles(dir, configFile string) {
	// load config
	groups, err := LoadConfigFromFile(configFile)
	if err != nil {
		Conf.Logger.Fatalln("can't load config", err)
	}

	// get list of all files in watch dir
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		Conf.Logger.Fatalln("can't read logs dir", err)
	}

	// assign file per group
	assignedFiles, err := assignFiles(files, groups)
	if err != nil {
		Conf.Logger.Fatalln("can't assign file per group", err)
	}

	doneCh := make(chan bool)
	assignedFilesCount := len(assignedFiles)

	for _, file := range assignedFiles {
		file.doneCh = doneCh
		go file.tail()
	}

	if Conf.ContinueWatch {
		go continueWatch(&dir, groups)
	}

	select {
	case done := <-doneCh:
		assignedFilesCount = -1
		if assignedFilesCount == 0 {
			debug(done)
			Conf.Logger.Println("done")
		}

	}
}

func assignFiles(files []os.FileInfo, groups []*Group) (outFiles []*File, err error) {
	for _, group := range groups {
		var assignedFiles []*File
		if assignedFiles, err = getFilesByGroup(files, group); err == nil {
			for _, assignedFile := range assignedFiles {
				outFiles = append(outFiles, assignedFile)
			}
		} else {
			return
		}
	}
	return
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
						continue
					}
					files = append(files, file)
					assignFiles(files, groups)
				}
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
	regex := *group.Mask
	for _, f := range allFiles {
		if !regex.MatchString(f.Name()) {
			continue
		}
		filepath := filepath.Join(Conf.WatchDir, f.Name())
		file, err := NewFile(filepath)
		if err != nil {
			return files, err
		}
		file.group = group
		files = append(files, file)
	}
	return files, nil
}

func NewFile(fpath string) (*File, error) {
	file := &File{}
	var err error
	if Conf.ReadWholeLog {
		Conf.Logger.Println("read whole logs")
		file.Tail, err = tail.TailFile(fpath, tail.Config{})
	} else {
		seekInfo := &tail.SeekInfo{Offset: 0, Whence: 2}
		file.Tail, err = tail.TailFile(fpath, tail.Config{Follow: true, ReOpen: true, Location: seekInfo})
	}
	return file, err
}

type File struct {
	Tail   *tail.Tail
	group  *Group
	doneCh chan bool
}

func (self *File) tail() {
	Conf.Logger.Printf("start tailing %+v", self.Tail.Filename)
	defer func() { self.doneCh <- true }()
	for line := range self.Tail.Lines {
		checkLineRules(&line.Text, self.group.Rules)
	}
}

func checkLineRule(line *string, rule *Rule) {
	match := rule.Match(line)
	if match != nil {
		rule.send(match)
	}
}

func checkLineRules(line *string, rules []*Rule) {
	for _, rule := range rules {
		checkLineRule(line, rule)
	}
}
