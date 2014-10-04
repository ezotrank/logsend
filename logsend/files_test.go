package logsend

import (
	"testing"
)

func benchmarkCheckLineRule(text *string, rule *Rule, b *testing.B) {
	for n := 0; n < b.N; n++ {
		CheckLineRule(text, rule)
	}
}

func BenchmarkCheckLineRuleNotMatch(b *testing.B) {
	str := `test word one`
	rule, err := NewRule(``)
	if err != nil {
		panic("failed")
	}
	benchmarkCheckLineRule(&str, rule, b)
}

func BenchmarkCheckLineRuleSimpleMatch(b *testing.B) {
	str := `test word one`
	rule, err := NewRule(`test word one`)
	if err != nil {
		panic("failed")
	}
	benchmarkCheckLineRule(&str, rule, b)
}

func BenchmarkCheckLineRuleSimpleMatchGroup(b *testing.B) {
	str := `test word one`
	rule, err := NewRule(`test (?P<test_STRING>\w+) one`)
	if err != nil {
		panic("failed")
	}
	benchmarkCheckLineRule(&str, rule, b)
}

func BenchmarkCheckLineRuleSimpleMultiMatchGroup(b *testing.B) {
	str := `test word one`
	rule, err := NewRule(`(?P<test_STRING>\w+) (?P<test1_STRING>\w+) (?P<test2_STRING>\w+)`)
	if err != nil {
		panic("failed")
	}
	benchmarkCheckLineRule(&str, rule, b)
}
