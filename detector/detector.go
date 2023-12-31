package detector

import (
	"strconv"
	"strings"
)

var initialismsRep *strings.Replacer = func() *strings.Replacer {
	var r []string
	for k := range initialisms {
		r = append(r, k, "")
	}
	return strings.NewReplacer(r...)
}()

var numRep *strings.Replacer = func() *strings.Replacer {
	var r []string
	for i := 0; i <= 9; i++ {
		r = append(r, strconv.Itoa(i), "")
	}
	return strings.NewReplacer(r...)
}()

func IsMixedCaps(s string) bool {
	s = strings.TrimPrefix(s, "_")
	if strings.Contains(s, "_") {
		return false
	}
	s = initialismsRep.Replace(s)
	s = numRep.Replace(s)
	if len(s) > 1 && strings.ToUpper(s) == s {
		return false
	}
	return true
}

func NoUnderscore(s string) bool {
	if s == "_" {
		return true
	}
	return !strings.Contains(s, "_")
}

func HasGetPrefix(s string) bool {
	return strings.HasPrefix(strings.ToLower(s), "get")
}
