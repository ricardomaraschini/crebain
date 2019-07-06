package match

import (
	"regexp"
	"strings"
)

// Multi is a custom type complying to flag.Value interface. In every subsequent call of `Set` it will add the rule to the previous one by an OR clause.
type Multi struct {
	rules   []string
	matcher *regexp.Regexp
}

func (mm Multi) String() string {
	return strings.Join(mm.rules, "\n")
}

func (mm *Multi) Set(value string) error {
	if err := mm.addRule(value); err != nil {
		return err
	}
	mm.rules = append(mm.rules, value)
	return nil
}

func (mm *Multi) addRule(rule string) error {
	rulesTotal := ""
	if mm.matcher != nil {
		rulesTotal = mm.matcher.String() + "|"
	}
	rulesTotal += rule

	exp, err := regexp.Compile(rulesTotal)
	if err != nil {
		return err
	}
	mm.matcher = exp
	return nil
}

func (mm Multi) Match(value string) bool {
	if mm.matcher == nil {
		return false
	}
	return mm.matcher.MatchString(value)
}
