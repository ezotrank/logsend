package logsend

import (
	"strconv"
)

func RegisterNewSender(name string, init func(interface{}), get func() Sender) {
	sender := &SenderRegister{
		init: init,
		get:  get,
	}
	Conf.registeredSenders[name] = sender
	return
}

type SenderRegister struct {
	init        func(interface{})
	get         func() Sender
	initialized bool
}

func (self *SenderRegister) Init(val interface{}) {
	self.init(val)
	self.initialized = true
}

type Configuration struct {
	ContinueWatch     bool
	DryRun            bool
	ReadWholeLog      bool
	ReadOnce          bool
	registeredSenders map[string]*SenderRegister
}

var Conf = &Configuration{
	registeredSenders: make(map[string]*SenderRegister),
}

var (
	rawConfig = make(map[string]interface{}, 0)
)

func i2float64(i interface{}) float64 {
	switch i.(type) {
	case string:
		val, _ := strconv.ParseFloat(i.(string), 32)
		return val
	case int:
		return float64(i.(int))
	case float64:
		return i.(float64)
	}
	panic(i)
}

func i2int(i interface{}) int {
	switch i.(type) {
	case string:
		val, _ := strconv.ParseFloat(i.(string), 32)
		return int(val)
	case int:
		return i.(int)
	case float64:
		return int(i.(float64))
	}
	panic(i)
}
