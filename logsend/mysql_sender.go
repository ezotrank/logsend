package logsend

import (
	"bytes"
	"database/sql"
	"flag"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
	"strings"
	"text/template"
)

// DOCS
// https://github.com/go-sql-driver/mysql#dsn-data-source-name
var (
	mysqlCh    = make(chan *string, 0)
	mysqlHost  = flag.String("mysql-host", "", "Example: user:password@/database?timeout=30s&strict=true")
	mysqlQuery = flag.String("mysql-query", "", "Example: insert into test1(teststring, testfloat) values('{{.gate}}', {{.exec_time}});")
)

func init() {
	RegisterNewSender("mysql", InitMysql, NewMysqlSender)
}

func InitMysql(conf interface{}) {
	host := conf.(map[string]interface{})["host"].(string)
	db, err := sql.Open("mysql", host)
	if err != nil {
		panic(err.Error())
	}

	go func() {
		defer db.Close()
		glog.Infoln("mysql queue is starts")
		for query := range mysqlCh {

			err = db.Ping()
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}

			// TODO: exec query with transaction support
			if Conf.DryRun {
				continue
			}

			trx, err := db.Begin()
			if err != nil {
				glog.Infoln("mysql init transaction ", err, *query)
			}

			for _, q := range strings.Split(*query, ";") {
				if q == "" {
					continue
				}
				if _, err = db.Exec(q + ";"); err != nil {
					break
				}

			}
			if err != nil {
				trx.Rollback()
				glog.Infoln("rollback ", err, *query)
				continue
			}
			trx.Commit()
		}
	}()
	return
}

func NewMysqlSender() Sender {
	mysqlSender := &MysqlSender{}
	mysqlSender.sendCh = mysqlCh
	return Sender(mysqlSender)
}

type MysqlSender struct {
	sendCh chan *string
	tmpl   *template.Template
}

func (self *MysqlSender) Name() string {
	return "mysql"
}

func (self *MysqlSender) SetConfig(rawConfig interface{}) error {
	var query string
	switch rawConfig.(map[string]interface{})["query"].(type) {
	case []interface{}:
		for _, s := range rawConfig.(map[string]interface{})["query"].([]interface{}) {
			query = query + s.(string)
		}
	case string:
		query = rawConfig.(map[string]interface{})["query"].(string)
	}
	self.tmpl, _ = template.New("query").Parse(query)
	return nil
}

func (self *MysqlSender) Send(data interface{}) {
	buf := new(bytes.Buffer)
	err := self.tmpl.Execute(buf, data)
	if err != nil {
		glog.Infoln("mysql template error ", err, data)
	}
	str := buf.String()
	self.sendCh <- &str
	return
}
