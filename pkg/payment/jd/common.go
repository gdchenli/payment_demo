package jd

import (
	"bytes"
	"crypto"
	"crypto/des"
	"crypto/md5"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type Jd struct{}

//生成随机字符串
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	byts := []byte(str)
	bytesLen := len(byts)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, byts[r.Intn(bytesLen)])
	}
	return string(result)
}

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

//支付字符串拼接
func GetNotEmptySortString(m map[string]string) string {
	var buf bytes.Buffer
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		vs := m[k]
		if vs == "" {
			continue
		}
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(k)
		buf.WriteByte('=')
		buf.WriteString(vs)
	}
	return buf.String()
}

//JDAPP填充规则
func JDPadding(origData []byte) []byte {
	merchantData := len(origData)
	x := (merchantData + 4) % 8

	y := 0
	if x == 0 {
		y = 0
	} else {
		y = 8 - x
	}

	sizeByte := IntegerToBytes(merchantData)
	var resultByte []byte

	//填充byte数据长度
	for i := 0; i < 4; i++ {
		resultByte = append(resultByte, sizeByte[i])
	}

	//填充原数据长度
	for j := 0; j < merchantData; j++ {
		resultByte = append(resultByte, origData[j])
	}

	//填充0
	for k := 0; k < y; k++ {
		resultByte = append(resultByte, 0x00)
	}

	return resultByte
}

func JDUnPadding(unPaddingResult []byte) []byte {

	var Result []byte
	var dataSizeByte []byte
	for i := 0; i < 4; i++ {
		dataSizeByte = append(dataSizeByte, unPaddingResult[i])
	}

	decimalDataSize := ByteArrayToInt(dataSizeByte)

	for j := 0; j < decimalDataSize; j++ {
		Result = append(Result, unPaddingResult[4+j])
	}

	return Result
}

//字节数组表示的实际长度
func ByteArrayToInt(dataSizeByte []byte) int {

	value := 0
	for i := 0; i < 4; i++ {
		shift := byte((4 - 1 - i) * 8)
		value = value + int(dataSizeByte[i]&0x000000FF)<<shift
	}
	return value
}

func IntegerToBytes(val int) [4]byte {
	byt := [4]byte{}
	byt[0] = byte(val >> 24 & 0xff)
	byt[1] = byte(val >> 16 & 0xff)
	byt[2] = byte(val >> 8 & 0xff)
	byt[3] = byte(val & 0xff)
	return byt
}

//byte转16进制字符串
func DecimalByteSlice2HexString(DecimalSlice []byte) string {
	var sa = make([]string, 0)
	for _, v := range DecimalSlice {
		sa = append(sa, fmt.Sprintf("%02X", v))
	}
	ss := strings.Join(sa, "")
	return ss
}

//十六进制字符串转byte
func HexString2Bytes(str string) ([]byte, error) {
	Bys, err := hex.DecodeString(str)
	if err != nil {
		return nil, err
	}
	return Bys, nil
}

// base解码
func BASE64DecodeStr(src string) string {
	a, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return ""
	}
	return string(a)
}

//Des解密
func decrypt(crypted, key []byte) ([]byte, error) {
	if len(crypted) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	out := make([]byte, len(crypted))
	dst := out
	bs := block.BlockSize()
	if len(crypted)%bs != 0 {
		return nil, errors.New("wrong crypted size")
	}

	for len(crypted) > 0 {
		block.Decrypt(dst, crypted[:bs])
		crypted = crypted[bs:]
		dst = dst[bs:]
	}

	return out, nil
}

//[golang ECB 3DES Decrypt]
func TripleEcbDesDecrypt(crypted, key []byte) ([]byte, error) {
	tkey := make([]byte, 24, 24)
	copy(tkey, key)
	k1 := tkey[:8]
	k2 := tkey[8:16]
	k3 := tkey[16:]
	buf1, err := decrypt(crypted, k3)
	if err != nil {
		return nil, err
	}
	buf2, err := encrypt(buf1, k2)
	if err != nil {
		return nil, err
	}
	out, err := decrypt(buf2, k1)
	if err != nil {
		return nil, err
	}
	out = JDUnPadding(out)
	return out, nil
}

func VerifyPKCS1v15(msg, sign, publicKey []byte, hashType crypto.Hash) bool {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return false
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	err = rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), hashType, msg, sign)
	return err == nil
}

//MD5加密
func Md5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)

	return hex.EncodeToString(cipherStr)
}

//sha256加密
func HaSha256(str string) string {
	h := sha256.New()
	h.Write([]byte(str))
	cipherStr := h.Sum(nil)
	//return cipherStr
	return hex.EncodeToString(cipherStr)
}

// base编码
func BASE64EncodeStr(src string) string {
	return base64.StdEncoding.EncodeToString([]byte(src))
}

//对消息的散列值进行数字签名
func SignPKCS1v15(msg, privateKey []byte, hashType crypto.Hash) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key format error")
	}
	pri, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("parse private key error")
	}

	sign, err := rsa.SignPKCS1v15(cryptoRand.Reader, pri, hashType, msg)
	if err != nil {
		return nil, errors.New("sign error")
	}
	return sign, nil
}

//Des加密
func encrypt(origData, key []byte) ([]byte, error) {
	if len(origData) < 1 || len(key) < 1 {
		return nil, errors.New("wrong data or key")
	}
	block, err := des.NewCipher(key)
	if err != nil {
		return nil, err
	}
	bs := block.BlockSize()
	if len(origData)%bs != 0 {
		return nil, errors.New("wrong padding")
	}
	out := make([]byte, len(origData))
	dst := out
	for len(origData) > 0 {
		block.Encrypt(dst, origData[:bs])
		origData = origData[bs:]
		dst = dst[bs:]
	}
	return out, nil
}

//[golang ECB 3DES Encrypt]
func TripleEcbDesEncrypt(origData, key []byte) ([]byte, error) {

	tkey := make([]byte, 24, 24)
	copy(tkey, key)
	k1 := tkey[:8]
	k2 := tkey[8:16]
	k3 := tkey[16:]

	origData = JDPadding(origData) // PKCS5Padding(origData, bs)
	buf1, err := encrypt(origData, k1)
	if err != nil {
		return nil, err
	}
	buf2, err := decrypt(buf1, k2)
	if err != nil {
		return nil, err
	}
	out, err := encrypt(buf2, k3)
	if err != nil {
		return nil, err
	}
	return out, nil
}
