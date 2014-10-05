package logsend

import (
	"reflect"
	"testing"
)

func TestPrepareValue(t *testing.T) {
	var key string
	var val interface{}
	var err error

	if key, val, err = PrepareValue("test", "word"); err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(key, "test") {
		t.Errorf("key does't equel")
	}
	if !reflect.DeepEqual(val, "word") {
		t.Errorf("val doesn't equal %+v %+v", val, "word")
	}

	if key, val, err = PrepareValue("test_STRING", "word"); err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(key, "test") {
		t.Errorf("key does't equel")
	}
	if !reflect.DeepEqual(val, "word") {
		t.Errorf("val doesn't equal %+v %+v", val, "word")
	}

	if key, val, err = PrepareValue("test_FLOAT", "23.11"); err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(key, "test") {
		t.Errorf("key does't equel")
	}
	if !reflect.DeepEqual(val, 23.11) {
		t.Errorf("val doesn't equal %+v %+v", val, "word")
	}

	if key, val, err = PrepareValue("test_INT", "23.11"); err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(key, "test") {
		t.Errorf("key does't equel")
	}
	if !reflect.DeepEqual(val.(int64), int64(23)) {
		t.Errorf("val doesn't equal %+v %+v", val, 23)
	}

	if key, val, err = PrepareValue("test_DurationToMillisecond", "1s"); err != nil {
		panic(err)
	}
	if !reflect.DeepEqual(key, "test") {
		t.Errorf("key does't equel")
	}
	if !reflect.DeepEqual(val.(int64), int64(1000)) {
		t.Errorf("val doesn't equal %+v %+v", val, 1000)
	}
}

func TestExtendValue(t *testing.T) {
	str := "test"
	data, err := ExtendValue(&str)
	if err != nil {
		t.Errorf("ExtendValue() err %+v", err)
	}
	if !reflect.DeepEqual(data, "test") {
		t.Errorf("ExtendValue() not equal %+v %+v", data, "test")
	}
}

func BenchmarkPrepareValueString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PrepareValue("test", "word")
	}
}

func BenchmarkPrepareValueString2(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PrepareValue("test_STRING", "word")
	}
}

func BenchmarkPrepareValueInt(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PrepareValue("test_INT", "12")
	}
}

func BenchmarkPrepareValueFloat(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PrepareValue("test_FLOAT", "12.11")
	}
}

func BenchmarkPrepareValueDurationToMillisecond(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PrepareValue("test_DurationToMillisecond", "1s")
	}
}
