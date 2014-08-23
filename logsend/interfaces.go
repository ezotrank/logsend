package logsend

type Sender interface {
	Send(interface{})
	SetConfig(interface{}) error
	Name() string
}
