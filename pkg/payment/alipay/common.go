package alipay

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"net/url"
	"sort"
)

type Alipay struct{}

//支付字符串拼接
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
