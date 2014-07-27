package logsend

import (
	influxdb "github.com/influxdb/influxdb/client"
	logpkg "log"
	"os"
)

var (
	log        = logpkg.New(os.Stderr, "", logpkg.Lmicroseconds)
	Debug      = true
	SendBuffer = 50
	SenderCh   = make(chan *influxdb.Series)
)

var Conf = struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	UDP        bool
	WatchDir   string
}{
	"localhost:8086",
	"root",
	"root",
	"test1",
	false,
	"",
}

func debug(msg ...interface{}) {
	if !Debug {
		return
	}
	log.Printf("debug: %+v", msg)
}
