package logsend

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func GetValue(vtype string) (result interface{}, err error) {
	switch vtype {
	case "GetHostname":
		result, err = getHostname()
	default:
		result = vtype
	}
	return
}

func getHostname() (interface{}, error) {
	return os.Hostname()
}

func ConvertToPoint(val, convert string) (result interface{}, err error) {
	switch convert {
	case "DurationToMillisecond":
		result, err = convertDurationToMillisecond(val)
	}
	return
}

func convertDurationToMillisecond(val string) (result interface{}, err error) {
	duration, err := time.ParseDuration(val)
	if err != nil {
		return
	}
	result = duration.Nanoseconds() / 1000 / 1000
	return
}

func LeadToType(val, valType string) (result interface{}, err error) {
	switch valType {
	case "int":
		result, err = strconv.ParseInt(val, 0, 64)
	case "float":
		result, err = strconv.ParseFloat(val, 64)
	default:
		result = val
	}
	return
}

func toFloat(val string) (float64, error) {
	return strconv.ParseFloat(val, 64)
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
	}
	return
}
