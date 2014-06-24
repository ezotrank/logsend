package main

import (
	"./logsend"
	"flag"
	logpkg "log"
	"os"
)

var (
	logDir     = flag.String("log-dir", "./tmp", "log directories")
	config     = flag.String("config", "config.json", "path to config.json file")
	dbhost     = flag.String("db-host", "localhost:8086", "db host")
	dbuser     = flag.String("db-user", "root", "db user")
	dbpassword = flag.String("db-password", "root", "db-password")
	database   = flag.String("database", "test1", "database")
	debug      = flag.Bool("debug", false, "turn on debug messages")
	sendBuffer = flag.Int("send-buffer", 25, "send buffer")
)

var (
	log = logpkg.New(os.Stderr, "", logpkg.Lmicroseconds)
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	logsend.SendBuffer = *sendBuffer
	logsend.Debug = *debug
	logsend.Conf.DBHost = *dbhost
	logsend.Conf.DBUser = *dbuser
	logsend.Conf.DBPassword = *dbpassword
	logsend.Conf.DBName = *database

	logsend.WatchLogs(*logDir, *config)
}
