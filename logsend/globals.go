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
	memprofile    *os.File
	Cpuprofile    string
	cpuprofile    *os.File
}

var (
	log     = logpkg.New(os.Stderr, "", logpkg.Lmicroseconds)
	senders = []Sender{}
)

var Conf = &Configuration{
	WatchDir:      "",
	ContinueWatch: true,
	Debug:         false,
	Memprofile:    "",
	Cpuprofile:    "",
}

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
	log.Printf("debug: %+v", msg)
}
