package detector

import "strings"

var initialisms = map[string]struct{}{
	"ACL":   {},
	"API":   {},
	"ASCII": {},
	"CPU":   {},
	"CSS":   {},
	"DDoS":  {},
	"DNS":   {},
	"EOF":   {},
	"GUID":  {},
	"HTML":  {},
	"ID":    {},
	"UUID":  {},
	"iOS":   {},
	"IP":    {},
	"JSON":  {},
	"QPS":   {},
	"RAM":   {},
	"RPC":   {},
	"URL":   {},
	"XML":   {},
}

var initialismsRep *strings.Replacer = func() *strings.Replacer {
	var r []string
	for k := range initialisms {
		r = append(r, k, "")
	}
	return strings.NewReplacer(r...)
}()

func IsMixedCaps(s string) bool {
	s = strings.TrimPrefix(s, "_")
	if strings.Contains(s, "_") {
		return false
	}
	s = initialismsRep.Replace(s)
	if len(s) > 1 && strings.ToUpper(s) == s {
		return false
	}
	return true
}
