package logsend

import (
	logpkg "log"
	"os"
	"runtime/pprof"
	"strconv"
)

type SenderRegister struct {
	Init func(interface{})
	Get  func() Sender
}

type Configuration struct {
	WatchDir          string
	ContinueWatch     bool
	Debug             bool
	Memprofile        string
	Logger            *logpkg.Logger
	DryRun            bool
	ReadWholeLog      bool
	ReadOnce          bool
	memprofile        *os.File
	Cpuprofile        string
	cpuprofile        *os.File
	registeredSenders map[string]*SenderRegister
}

var Conf = &Configuration{
	WatchDir:          "",
	Memprofile:        "",
	Cpuprofile:        "",
	Logger:            logpkg.New(os.Stderr, "", logpkg.Ldate|logpkg.Ltime|logpkg.Lshortfile),
	registeredSenders: make(map[string]*SenderRegister),
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

func i2float64(i interface{}) float64 {
	switch i.(type) {
	case string:
		val, _ := strconv.ParseFloat(i.(string), 32)
		return val
	}
	panic(i)
}

func i2int(i interface{}) int {
	switch i.(type) {
	case string:
		val, _ := strconv.ParseFloat(i.(string), 32)
		return int(val)
	}
	panic(i)
}
