package logsend

import (
	"testing"
)

func benchmarkCheckLineRule(text *string, rule *Rule, b *testing.B) {
        for n := 0; n < b.N; n++ {
			CheckLineRule(text, rule)
        }
}

func BenchmarkCheckLineRuleNotMatch(b *testing.B)  {
	str := `test word one`
	regex := ""
	rule := &Rule{
		Regexp: &regex,
	}
	rule.loadRegexp()
	benchmarkCheckLineRule(&str, rule, b)
}

func BenchmarkCheckLineRuleSimpleMatch(b *testing.B)  {
	str := `test word one`
	regex := `test word one`
	rule := &Rule{
		Regexp: &regex,
	}
	rule.loadRegexp()
	benchmarkCheckLineRule(&str, rule, b)
}

func BenchmarkCheckLineRuleSimpleMatchGroup(b *testing.B)  {
	str := `test word one`
	regex := `test (?P<test_STRING>\w+) one`
	rule := &Rule{
		Regexp: &regex,
	}
	rule.loadRegexp()
	benchmarkCheckLineRule(&str, rule, b)
}

func BenchmarkCheckLineRuleSimpleMultiMatchGroup(b *testing.B)  {
	str := `test word one`
	regex := `(?P<test_STRING>\w+) (?P<test1_STRING>\w+) (?P<test2_STRING>\w+)`
	rule := &Rule{
		Regexp: &regex,
	}
	rule.loadRegexp()
	benchmarkCheckLineRule(&str, rule, b)
}