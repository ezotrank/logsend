package logsend

import (
	logpkg "log"
	"os"
)

var (
	log        = logpkg.New(os.Stderr, "", logpkg.Lmicroseconds)
	Debug      = true
	SendBuffer = 50
)

func debug(msg ...interface{}) {
	if !Debug {
		return
	}
	log.Printf("debug: %+v", msg)
}
