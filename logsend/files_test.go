package logsend

import (
	"regexp"
	"testing"
)

func benchmarkCheckLineRule(text *string, rule *Rule, b *testing.B) {
	for n := 0; n < b.N; n++ {
		CheckLineRule(text, rule)
	}
}

func BenchmarkCheckLineRuleNotMatch(b *testing.B) {
	str := `test word one`
	regex := regexp.MustCompile(``)
	rule := &Rule{
		regexp: regex,
	}
	benchmarkCheckLineRule(&str, rule, b)
}

func BenchmarkCheckLineRuleSimpleMatch(b *testing.B) {
	str := `test word one`
	regex := regexp.MustCompile(`test word one`)
	rule := &Rule{
		regexp: regex,
	}
	benchmarkCheckLineRule(&str, rule, b)
}

func BenchmarkCheckLineRuleSimpleMatchGroup(b *testing.B) {
	str := `test word one`
	regex := regexp.MustCompile(`test (?P<test_STRING>\w+) one`)
	rule := &Rule{
		regexp: regex,
	}
	benchmarkCheckLineRule(&str, rule, b)
}

func BenchmarkCheckLineRuleSimpleMultiMatchGroup(b *testing.B) {
	str := `test word one`
	regex := regexp.MustCompile(`(?P<test_STRING>\w+) (?P<test1_STRING>\w+) (?P<test2_STRING>\w+)`)
	rule := &Rule{
		regexp: regex,
	}
	benchmarkCheckLineRule(&str, rule, b)
}
