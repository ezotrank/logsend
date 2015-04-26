package main

import (
	"flag"
	"fmt"
	"github.com/ezotrank/logsend/logsend"
	"os"
)

const (
	VERSION = "2.0"
)

var (
	config        = flag.String("config", "", "path to config.json file")
	check         = flag.Bool("check", false, "check config.json")
	continueWatch = flag.Bool("continue-watch", false, "watching folder for new files")
	dryRun        = flag.Bool("dry-run", false, "not send data")
	readWholeLog  = flag.Bool("read-whole-log", false, "read whole logs")
	readOnce      = flag.Bool("read-once", false, "read logs once and exit")
	regex         = flag.String("regex", "", "regex rule")
	version       = flag.Bool("version", false, "show version number")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Printf("logsend version %s\n", VERSION)
		os.Exit(0)
	}

	logsend.Conf.ContinueWatch = *continueWatch
	logsend.Conf.DryRun = *dryRun
	logsend.Conf.ReadWholeLog = *readWholeLog
	logsend.Conf.ReadOnce = *readOnce

	if *check {
		_, err := logsend.LoadConfigFromFile(*config)
		if err != nil {
			fmt.Printf("config check err: %s", err)
			os.Exit(1)
		}
		fmt.Println("ok")
		os.Exit(0)
	}

	fi, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if fi.Mode()&os.ModeNamedPipe == 0 {
		if len(flag.Args()) < 1 {
			fmt.Printf("you forget specify watch directories")
		}
		logsend.WatchFiles(flag.Args(), *config)
	} else {
		flag.VisitAll(logsend.LoadRawConfig)
		logsend.ProcessStdin()
	}
	os.Exit(0)
}
