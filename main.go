package main

import (
	"flag"
	"fmt"
	"github.com/ezotrank/logsend/logsend"
	logpkg "log"
	"os"
	"runtime"
)

var (
	watchDir     = flag.String("watch-dir", "./tmp", "log directories")
	config       = flag.String("config", "config.json", "path to config.json file")
	check        = flag.Bool("check", false, "check config.json")
	debug        = flag.Bool("debug", false, "turn on debug messages")
	stopContinue = flag.Bool("stop-continue", false, "watching folder for new files")
	logFile      = flag.String("log", "", "log file")
	dryRun       = flag.Bool("dry-run", false, "not send data")
	memprofile   = flag.String("memprofile", "", "memory profiler")
	maxprocs     = flag.Int("maxprocs", 0, "max count of cpu")
)

func main() {
	flag.Parse()

	if *logFile != "" {
		file, err := os.OpenFile(*logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Errorf("Failed to open log file: %+v\n", err)
		}
		logsend.Conf.Logger = logpkg.New(file, "", logpkg.Ldate|logpkg.Ltime|logpkg.Lshortfile)
	}

	if *maxprocs <= 0 {
		*maxprocs = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(*maxprocs)
	fmt.Printf("set GOMAXPROCS to %v\n", *maxprocs)

	logsend.Conf.Debug = *debug
	logsend.Conf.ContinueWatch = !*stopContinue
	logsend.Conf.WatchDir = *watchDir
	logsend.Conf.Memprofile = *memprofile
	logsend.Conf.DryRun = *dryRun

	if *check {
		_, err := logsend.LoadConfig(*config)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("ok")
		os.Exit(0)
	}

	logsend.WatchFiles(*watchDir, *config)
}
