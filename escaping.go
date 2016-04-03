package vcardenc

import "strings"

func escapedJoin(v []string, delimiter string) string {
	var vo []string
	for _, e := range v {
		vo = append(vo, escape(e))
	}
	return strings.Join(vo, delimiter)
}

func escapeQuoted(s string) string {
	if len(s) < 3 {
		return s
	}
	return `"` + strings.Replace(s[1:len(s)-2], `"`, `\"`, -1) + `"`
}

func escapeParam(s string) (o string) {
	if s == "" {
		return s
	}
	rs := []rune(s)
	if len(rs) > 2 && (rs[0] == '"' && rs[len(rs)-1] == '"') {
		return escapeQuoted(s)
	}
	return escape(s)
}

func escape(s string) (o string) {
	if s == "" {
		return s
	}
	o = strings.Replace(s, "\\", "\\\\", -1)
	o = strings.Replace(o, "\r\n", "\n", -1)
	o = strings.Replace(o, "\n", "\\n", -1)
	o = strings.Replace(o, ";", "\\;", -1)
	//o = strings.Replace(o, ":", "\\:", -1)
	//o = strings.Replace(o, ",", "\\,", -1)
	return o
}
