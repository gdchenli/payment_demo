package util

import (
	"bytes"
	"sort"
	"strings"
)

func GetSortString(m map[string]string) string {
	var buf bytes.Buffer
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		vs := m[k]
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(vs)
	}
	return buf.String()
}

//特殊字符过滤
func SpecialReplace(str string) string {
	specialChars := []string{"+", "-", "×", "<", ">", "#", "[", "]", "(", ")", "（", "）", "/", "?", "&", ".", "{", "}", "「", "」"}
	for _, spc := range specialChars {
		str = strings.Replace(str, spc, "", -1)
	}
	return str
}
