package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/ActiveState/tail"
	log "github.com/ezotrank/logger"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
	"github.com/golang/snappy"
)

var (
	server = flag.String("server", "http://localhost:8000/config", "server host")
)

type Config struct {
	PushAddr  string
	FilesMask []string
	ChunkSize int
	Compress bool
}

type LogLine struct {
	Ts   int64
	Line []byte
}

func MarshaLogLines(loglines []*LogLine) []byte {
	b, err := json.Marshal(loglines)
	if err != nil {
		panic(err)
	}
	return b
}

func hostname() string {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return name
}

func getConfig(server string) (*Config, error) {
	resp, err := http.Get(fmt.Sprintf("%v?host=%v", server, hostname()))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	config := &Config{}
	err = json.Unmarshal(body, config)
	return config, err
}

func getFiles(path string) []*tail.Tail {
	files := make([]*tail.Tail, 0)
	dir, mask := filepath.Split(path)
	if err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		matched, err := filepath.Match(mask, f.Name())
		if err == nil && matched {
			file, err := tail.TailFile(path, tail.Config{Follow: true, ReOpen: true, Location: &tail.SeekInfo{Offset: 0, Whence: 2}})
			if err != nil {
				log.Infoln(err)
				return err
			}
			log.Infof("add file %q for watching", file.Filename)
			files = append(files, file)
		}
		return err
	}); err != nil {
		panic(err)
	}
	return files
}

func process(config *Config) {
	conn, err := net.Dial("tcp", config.PushAddr)
	if err != nil {
		panic(err)
	}
	files := make([]*tail.Tail, 0)
	defer func() {
		conn.Close()
		for _, f := range files {
			log.Infof("Stop listen files %q", f.Filename)
			f.Stop()
		}
	}()

	for _, mask := range config.FilesMask {
		files = append(files, getFiles(mask)...)
	}
	ch := make(chan *LogLine, 0)
	fn := func(t *tail.Tail) {
		log.Infof("start tailing %q", t.Filename)
		defer func() { t.Stop() }()
		for line := range t.Lines {
			ch <- &LogLine{
				Ts:   time.Now().UTC().UnixNano(),
				Line: []byte(line.Text),
			}
		}
	}
	for _, file := range files {
		go fn(file)
	}

	chunk := make([]*LogLine, 0)
	for {
		select {
		case logLine := <-ch:
			chunk = append(chunk, logLine)
			if len(chunk) >= config.ChunkSize {
				data := MarshaLogLines(chunk)
				if config.Compress {
					data = snappy.Encode(make([]byte,0), data)
				}
				conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
				bytesCount, err := conn.Write(data)
				if err != nil {
					log.Errorf("can't write data to socket %v", err)
					panic(err)
				}
				log.Debugf("write %d bytes to socket", bytesCount)
				chunk = make([]*LogLine, 0)
			}
		}
	}
}

func run(server string) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorln("Recovered in f", r)
			time.Sleep(100 * time.Millisecond)
			run(server)
		}
	}()
	config, err := getConfig(server)
	if err != nil {
		panic(err)
	}
	process(config)
}

func main() {
	flag.Parse()
	log.Infoln("Start agent")
	run(*server)
}
