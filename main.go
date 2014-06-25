package main

import (
	"./logsend"
	"flag"
)

var (
	logDir     = flag.String("log-dir", "./tmp", "log directories")
	config     = flag.String("config", "config.json", "path to config.json file")
	dbhost     = flag.String("db-host", "localhost:8086", "db host")
	dbuser     = flag.String("db-user", "root", "db user")
	dbpassword = flag.String("db-password", "root", "db-password")
	database   = flag.String("database", "test1", "database")
	udp        = flag.Bool("udp", false, "send series over udp")
	debug      = flag.Bool("debug", false, "turn on debug messages")
	sendBuffer = flag.Int("send-buffer", 25, "send buffer")
)

func main() {
	flag.Parse()

	logsend.SendBuffer = *sendBuffer
	logsend.Debug = *debug
	logsend.Conf.DBHost = *dbhost
	logsend.Conf.DBUser = *dbuser
	logsend.Conf.DBPassword = *dbpassword
	logsend.Conf.DBName = *database
	logsend.Conf.UDP = *udp

	logsend.WatchLogs(*logDir, *config)
}
