package allpay

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"sort"
	"strings"
)

type Allpay struct{}

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

func ParseQueryString(str string) (map[string]string, error) {
	queryMap := make(map[string]string)
	values, err := url.ParseQuery(str)
	if err != nil {
		return queryMap, err
	}
	for k, v := range values {
		queryMap[k] = v[0]
	}

	return queryMap, nil
}

//json转map
func JsonToMap(paramString string) (paramMap map[string]string, err error) {
	if err := json.Unmarshal([]byte(paramString), &paramMap); err != nil {
		return paramMap, err
	}
	return paramMap, nil
}

// base解码
func BASE64DecodeStr(src string) string {
	a, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return ""
	}
	return string(a)
}

//MD5加密
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)

	return hex.EncodeToString(cipherStr)
}

// base编码
func BASE64EncodeStr(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

//sha256加密
func Hsha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	return hex.EncodeToString(cipherStr)
}
