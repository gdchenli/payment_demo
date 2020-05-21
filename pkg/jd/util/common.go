package util

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

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
func GetPayString(m map[string]string) string {
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
func GetNotEmptyPayString(m map[string]string) string {
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
