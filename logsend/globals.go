package logsend

import (
	logpkg "log"
	"os"
)

var (
	log        = logpkg.New(os.Stderr, "", logpkg.Lmicroseconds)
	debugState = true
)

func debug(msg ...interface{}) {
	if !debugState {
		return
	}
	log.Printf("debug: %+v", msg)
}
