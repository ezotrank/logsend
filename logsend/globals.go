package logsend

import (
	influxdb "github.com/influxdb/influxdb/client"
	logpkg "log"
	"os"
	"runtime/pprof"
)

type Configuration struct {
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	UDP           bool
	WatchDir      string
	ContinueWatch bool
	SendBuffer    int
	Debug         bool
	Memprofile    string
	memprofile    *os.File
	Cpuprofile    string
	cpuprofile    *os.File
}

var (
	log      = logpkg.New(os.Stderr, "", logpkg.Lmicroseconds)
	SenderCh = make(chan *influxdb.Series)
)

var Conf = &Configuration{
	DBHost:        "localhost:8086",
	DBUser:        "root",
	DBPassword:    "root",
	DBName:        "test1",
	UDP:           false,
	WatchDir:      "",
	ContinueWatch: true,
	Debug:         false,
	Memprofile:    "",
	Cpuprofile:    "",
	SendBuffer:    8,
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
