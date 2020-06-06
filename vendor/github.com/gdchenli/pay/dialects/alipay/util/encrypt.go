package util

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

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
