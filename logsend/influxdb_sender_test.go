package logsend

import (
	"reflect"
	"testing"
)

func TestName(t *testing.T) {
	var want, result string
	want = "influxdb"
	sender := &InfluxdbSender{}
	result = sender.Name()
	if !reflect.DeepEqual(result, want) {
		t.Errorf("Name() returned %+v, want %+v", result, want)
	}
}
