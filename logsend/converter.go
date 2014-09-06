package logsend

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func getHostname() (interface{}, error) {
	return os.Hostname()
}

func extendValue(name *string) (val interface{}, err error) {
	switch *name {
	default:
		val, err = *name, nil
	case "getHostName":
		val, err = getHostname()
	}
	return
}

func toFloat(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
}

func toInt(val string) (int64, error) {
	return strconv.ParseInt(val, 0, 64)
}

func durationToMillisecond(val *string) (result interface{}, err error) {
	duration, err := time.ParseDuration(*val)
	if err != nil {
		return
	}
	result = duration.Nanoseconds() / 1000 / 1000
	return
}

func prepareValue(source, data string) (key string, val interface{}, err error) {
	tSource := strings.Split(source, "_")
	key = strings.Join(tSource[:len(tSource)-1], "_")
	keyType := tSource[len(tSource)-1]
	switch keyType {
	default:
		val = data
	case "FLOAT":
		val, err = toFloat(data)
	case "INT":
		val, err = toInt(data)
	case "DurationToMillisecond":
		val, err = durationToMillisecond(&data)
	}
	return
}
