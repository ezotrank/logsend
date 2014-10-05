package logsend

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type RegisteredConvertes map[string]func(interface{}) (interface{}, error)
type RegisteredExtendValues map[string]func() (interface{}, error)

var (
	registeredConvertes             = make(RegisteredConvertes, 0)
	registeredConvertesExtendValues = make(RegisteredExtendValues, 0)
)

func init() {
	registeredConvertes["STRING"] = ci2string
	registeredConvertes["FLOAT"] = ci2float
	registeredConvertes["INT"] = ci2int
	registeredConvertes["DurationToMillisecond"] = ci2DurationToMillisecond

	registeredConvertesExtendValues["HOST"] = revHost
}

func ci2string(i interface{}) (o interface{}, err error) {
	switch i.(type) {
	default:
		err = errors.New("interface not a string")
	case string:
		o = i
	}
	return
}

func ci2float(i interface{}) (o interface{}, err error) {
	switch i.(type) {
	default:
		err = errors.New("interface not a float")
	case string:
		var fl float64
		fl, err = strconv.ParseFloat(i.(string), 64)
		o = fl
	case float64:
		o = i
	}
	return
}

func ci2int(i interface{}) (o interface{}, err error) {
	switch i.(type) {
	default:
		err = errors.New("interface not a int")
	case string:
		var fl float64
		fl, err = strconv.ParseFloat(i.(string), 64)
		o = int64(fl)
	case float64:
		o = int64(i.(float64))
	}
	return
}

func ci2DurationToMillisecond(i interface{}) (o interface{}, err error) {
	switch i.(type) {
	default:
		err = errors.New("interface not a string")
	case string:
		var duration time.Duration
		duration, err = time.ParseDuration(i.(string))
		if err != nil {
			return
		}
		o = duration.Nanoseconds() / 1000 / 1000
	}
	return
}

func revHost() (interface{}, error) {
	return os.Hostname()
}

func ExtendValue(name *string) (val interface{}, err error) {
	if fn, ok := registeredConvertesExtendValues[*name]; ok {
		val, err = fn()
	} else {
		val = *name
	}
	return
}

func PrepareValue(source, data string) (key string, val interface{}, err error) {
	tSource := strings.Split(source, "_")
	key = strings.Join(tSource[:len(tSource)-1], "_")
	if key == "" {
		key = source
	}
	keyType := tSource[len(tSource)-1]
	if fn, ok := registeredConvertes[keyType]; ok {
		val, err = fn(data)
	} else {
		val = data
	}
	return
}
