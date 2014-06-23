package logsend

import (
	"strconv"
	"time"
)

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
	default:
		result = val
	}
	return
}
