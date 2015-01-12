package logsend

import (
	"github.com/ActiveState/tail"
	"github.com/golang/glog"
	"github.com/howeyc/fsnotify"
	"os"
	"path/filepath"
)

func walkLogDir(dir string) (files []string, err error) {
	if string(dir[len(dir)-1]) != "/" {
		dir = dir + "/"
	}
	visit := func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		abs, err := filepath.Abs(path)
		if err != nil {
			glog.Fatalln(err)
		}
		files = append(files, abs)
		return nil
	}
	err = filepath.Walk(dir, visit)
	return
}

func WatchFiles(dirs []string, configFile string) {
	// load config
	groups, err := LoadConfigFromFile(configFile)
	if err != nil {
		glog.Fatalln("can't load config", err)
	}

	// get list of all files in watch dir
	files := make([]string, 0)
	for _, dir := range dirs {
		fs, err := walkLogDir(dir)
		if err != nil {
			panic(err)
		}
		for _, f := range fs {
			files = append(files, f)
		}
	}

	// assign file per group
	assignedFiles, err := assignFiles(files, groups)
	if err != nil {
		glog.Fatalln("can't assign file per group", err)
	}

	doneCh := make(chan string)
	assignedFilesCount := len(assignedFiles)

	for _, file := range assignedFiles {
		file.doneCh = doneCh
		go file.tail()
	}

	if Conf.ContinueWatch {
		for _, dir := range dirs {
			go continueWatch(&dir, groups)
		}
	}

	for {
		select {
		case fpath := <-doneCh:
			assignedFilesCount = assignedFilesCount - 1
			if assignedFilesCount == 0 {
				glog.Infof("finished reading file %+v", fpath)
				if Conf.ReadOnce {
					return
				}
			}

		}
	}

}

func assignFiles(files []string, groups []*Group) (outFiles []*File, err error) {
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

func getFilesByGroup(allFiles []string, group *Group) ([]*File, error) {
	files := make([]*File, 0)
	regex := *group.Mask
	for _, f := range allFiles {
		if !regex.MatchString(filepath.Base(f)) {
			continue
		}
		file, err := NewFile(f)
		if err != nil {
			return files, err
		}
		file.group = group
		files = append(files, file)
	}
	return files, nil
}

func continueWatch(dir *string, groups []*Group) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		glog.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() {
					files := make([]string, 0)
					file, err := filepath.Abs(ev.Name)
					if err != nil {
						glog.Infof("can't get file %+v", err)
						continue
					}
					files = append(files, file)
					assignFiles(files, groups)
				}
			case err := <-watcher.Error:
				glog.Infoln("error:", err)
			}
		}
	}()

	err = watcher.Watch(*dir)
	if err != nil {
		glog.Fatal(err)
	}

	<-done

	/* ... do stuff ... */
	watcher.Close()
}

func NewFile(fpath string) (*File, error) {
	file := &File{}
	var err error
	if Conf.ReadWholeLog && Conf.ReadOnce {
		glog.Infof("read whole file once %+v", fpath)
		file.Tail, err = tail.TailFile(fpath, tail.Config{})
	} else if Conf.ReadWholeLog {
		glog.Infof("read whole file and continue %+v", fpath)
		file.Tail, err = tail.TailFile(fpath, tail.Config{Follow: true, ReOpen: true})
	} else {
		seekInfo := &tail.SeekInfo{Offset: 0, Whence: 2}
		file.Tail, err = tail.TailFile(fpath, tail.Config{Follow: true, ReOpen: true, Location: seekInfo})
	}
	return file, err
}

type File struct {
	Tail   *tail.Tail
	group  *Group
	doneCh chan string
}

func (self *File) tail() {
	glog.Infof("start tailing %+v", self.Tail.Filename)
	defer func() { self.doneCh <- self.Tail.Filename }()
	for line := range self.Tail.Lines {
		checkLineRules(&line.Text, self.group.Rules)
	}
}

func checkLineRule(line *string, rule *Rule) {
	if glog.V(2) {
		glog.Infof("get line: %s\n", *line)
	}
	match := rule.Match(line)
	if match != nil {
		if glog.V(2) {
			glog.Infof("match regexp and line %q  %q", rule.regexp, *line)
		}
		rule.send(match)
	}
}

func checkLineRules(line *string, rules []*Rule) {
	for _, rule := range rules {
		checkLineRule(line, rule)
	}
}
