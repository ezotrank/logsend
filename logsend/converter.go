package logsend

import (
	"os"
	"strconv"
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
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
			log.Printf("can't convert %+v and %+v", val, convert)
		}
	}()
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
