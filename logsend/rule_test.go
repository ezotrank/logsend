package logsend

import (
	"reflect"
	"regexp"
	"testing"
)

func TestNewRule(t *testing.T) {
	reStr := "\\w+ (?P<test>.*)"
	rule, err := NewRule(reStr)
	if err != nil {
		t.Errorf("%+v", err)
	}
	if rule.regexp.String() != regexp.MustCompile(reStr).String() {
		t.Errorf("TestNewRule regexp not match %+v %+v", rule.regexp.String(), regexp.MustCompile(reStr).String())
	}
	if !reflect.DeepEqual(rule.subexpNames, []string{"", "test"}) {
		t.Errorf("subxpNames not match %+v %+v", rule.subexpNames, []string{"", "test"})
	}
}
