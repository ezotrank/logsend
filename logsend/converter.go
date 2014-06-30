package logsend

import (
	"os"
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