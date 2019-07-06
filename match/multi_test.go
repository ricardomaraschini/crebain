package match

import (
	"fmt"
	"testing"
)

func TestString(t *testing.T) {
	rules := []string{`aa`, `bb`}
	expected := "aa\nbb"

	mm := generateMulti(t, rules)
	got := mm.String()
	if got != expected {
		t.Fatal("Unexpected string result: ", got)
	}
}

func TestOneMatch(t *testing.T) {
	testStrings := []string{
		"I change my mind, I take it back",
		"/var/logs/something/",
		"./relativeDir/file.ext",
	}

	subtests := []struct {
		name          string
		rules         []string
		expectedMatch []bool
	}{
		{
			name:          "no rules",
			rules:         []string{},
			expectedMatch: []bool{false, false, false},
		},
		{
			name:          "simple",
			rules:         []string{`change`},
			expectedMatch: []bool{true, false, false},
		},
		{
			name:          "multimaching:morerules",
			rules:         []string{`^\./`, `^/`},
			expectedMatch: []bool{false, true, true},
		},
		{
			name:          "multimaching:OR",
			rules:         []string{`^\./|^/`},
			expectedMatch: []bool{false, true, true},
		},
		{
			name:          "selected",
			rules:         []string{`something/`, `relativeDir/`},
			expectedMatch: []bool{false, true, true},
		},
	}

	for _, st := range subtests {
		t.Run(st.name, func(t *testing.T) {
			multiMatcher := generateMulti(t, st.rules)
			for i, str := range testStrings {
				t.Run(fmt.Sprintf("%d: %s", i, str), func(t *testing.T) {
					if got := multiMatcher.Match(str); got != st.expectedMatch[i] {
						t.Fatalf("Unexpected match: should be %t but %t returned", st.expectedMatch[i], got)
					}
				})
			}
		})
	}

	t.Run("invalid input", func(t *testing.T) {
		// In this test only the last rule is broken.
		subtests := []struct {
			name  string
			rules []string
		}{
			{
				"1", []string{`gollum\`},
			},
			{
				"2", []string{`gollum`, `\`},
			},
		}
		for _, st := range subtests {
			t.Run(st.name, func(t *testing.T) {
				lastElPos := len(st.rules) - 1
				mm := generateMulti(t, st.rules[:lastElPos])

				if err := mm.Set(st.rules[lastElPos]); err == nil {
					t.Fatal("No error returned")
				}
			})
		}
	})
}

func generateMulti(t *testing.T, rules []string) Multi {
	t.Helper()
	regexs := Multi{}
	for _, r := range rules {
		if err := regexs.Set(r); err != nil {
			t.Fatal(err)
		}
	}
	return regexs
}
