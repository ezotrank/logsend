package logsend

import (
	logpkg "log"
	"os"
	"runtime/pprof"
)

type Configuration struct {
	WatchDir      string
	ContinueWatch bool
	Debug         bool
	Memprofile    string
	Logger        *logpkg.Logger
	DryRun        bool
	ReadWholeLog  bool
	ReadOnce      bool
	memprofile    *os.File
	Cpuprofile    string
	cpuprofile    *os.File
}

var Conf = &Configuration{
	WatchDir:   "",
	Memprofile: "",
	Cpuprofile: "",
	Logger:     logpkg.New(os.Stderr, "", logpkg.Ldate|logpkg.Ltime|logpkg.Lshortfile),
}

var (
	senders   = []Sender{}
	rawConfig = make(map[string]interface{}, 0)
)

func mempprof() {
	if Conf.memprofile == nil {
		Conf.memprofile, _ = os.Create(Conf.Memprofile)
	}
	pprof.WriteHeapProfile(Conf.memprofile)
}

func debug(msg ...interface{}) {
	if !Conf.Debug {
		return
	}
	Conf.Logger.Printf("debug: %+v", msg)
}
