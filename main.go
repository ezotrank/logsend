package main

import (
	"flag"
	"fmt"
	"github.com/andrew-d/go-termutil"
	"github.com/ezotrank/logsend/logsend"
	logpkg "log"
	"os"
)

var (
	watchDir      = flag.String("watch-dir", "./tmp", "log directories")
	config        = flag.String("config", "", "path to config.json file")
	check         = flag.Bool("check", false, "check config.json")
	debug         = flag.Bool("debug", false, "turn on debug messages")
	continueWatch = flag.Bool("continue-watch", false, "watching folder for new files")
	logFile       = flag.String("log", "", "log file")
	dryRun        = flag.Bool("dry-run", false, "not send data")
	memprofile    = flag.String("memprofile", "", "memory profiler")
	maxprocs      = flag.Int("maxprocs", 0, "max count of cpu")
	readWholeLog  = flag.Bool("read-whole-log", false, "read whole logs")
	readOnce      = flag.Bool("read-once", false, "read logs once and exit")
	regex         = flag.String("regex", "", "regex rule")
)

func main() {
	flag.Parse()

	if *logFile != "" {
		file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Errorf("Failed to open log file: %+v\n", err)
		}
		defer file.Close()
		logsend.Conf.Logger = logpkg.New(file, "", logpkg.Ldate|logpkg.Ltime|logpkg.Lshortfile)
	}

	logsend.Conf.Debug = *debug
	logsend.Conf.ContinueWatch = *continueWatch
	logsend.Conf.WatchDir = *watchDir
	logsend.Conf.Memprofile = *memprofile
	logsend.Conf.DryRun = *dryRun
	logsend.Conf.ReadWholeLog = *readWholeLog
	logsend.Conf.ReadOnce = *readOnce

	if *check {
		_, err := logsend.LoadConfigFromFile(*config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("ok")
		os.Exit(0)
	}

	if termutil.Isatty(os.Stdin.Fd()) {
		logsend.WatchFiles(*watchDir, *config)
	} else {
		flag.VisitAll(logsend.LoadRawConfig)
		logsend.ProcessStdin()
	}
	os.Exit(0)
}
