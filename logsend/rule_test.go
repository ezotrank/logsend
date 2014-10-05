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

func TestMatch(t *testing.T) {
	line := `string name: a1 with float 133.33 and int 313`
	var rule *Rule
	var result interface{}

	rule, _ = NewRule("match nothing")
	result = rule.Match(&line)
	if !reflect.DeepEqual(result, nil) {
		t.Errorf("rule.Match failed %+v %+v", result, nil)
	}

	rule, _ = NewRule("string name")
	result = rule.Match(&line)
	if !reflect.DeepEqual(result, true) {
		t.Errorf("rule.Match failed %+v %+v", result, true)
	}

	rule, _ = NewRule(`name: (?P<str>\w+).* float (?P<float_FLOAT>\d+.\d+).* int (?P<int_INT>.+)`)
	result = rule.Match(&line)
	want := map[string]interface{}{"str": "a1", "int": int64(313), "float": float64(133.33)}
	if !reflect.DeepEqual(result, want) {
		t.Errorf("rule.Match failed %+v %+v", result, want)
	}
}

func BenchmarkMatchNotMath(b *testing.B) {
	line := `string name: a1 with float 133.33 and int 313`
	rule, _ := NewRule(`match nothing`)
	for n := 0; n < b.N; n++ {
		rule.Match(&line)
	}
}

func BenchmarkMatchBoolMatch(b *testing.B) {
	line := `string name: a1 with float 133.33 and int 313`
	rule, _ := NewRule(`string name`)
	for n := 0; n < b.N; n++ {
		rule.Match(&line)
	}
}

func BenchmarkMatchGroupMatch(b *testing.B) {
	line := `string name: a1 with float 133.33 and int 313`
	rule, _ := NewRule(`name: (?P<str>\w+).* float (?P<float_FLOAT>\d+.\d+).* int (?P<int_INT>.+)`)
	for n := 0; n < b.N; n++ {
		rule.Match(&line)
	}
}
