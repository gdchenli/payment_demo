package jd

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"payment_demo/pkg/payment/common"

	"github.com/sirupsen/logrus"
)

const (
	CallbackSuccessCode = "0" //支付成功
)

const (
	CallbackEncryptFormatErrCode    = 10301
	CallbackEncryptFormatErrMessage = "同步通知，加密数据格式错误"
	CallbackDecryptFailedErrCode    = 10302
	CallbackDecryptFailedErrMessage = "同步通知，解密失败"
	CallbackDecryptFormatErrCode    = 10303
	CallbackDecryptFormatErrMessage = "同步通知，解密据格式错误"
	CallbackSignErrCode             = 10304
	CallbackSignErrMessage          = "同步通知，签名校验失败"
	CallbackStatusErrCode           = 10305
	CallbackStatusErrMessage        = "同步通知，交易状态不正确"
)

type VerifyQuery struct {
	TradeNum  string `json:"tradeNum"`  //订单号
	Amount    string `json:"amount"`    //交易金额
	Currency  string `json:"currency"`  //货币类型
	TradeTime string `json:"tradeTime"` //交易时间
	Status    string `json:"status"`    //交易状态
	Sign      string `json:"sign"`      //签名
}

func (jd *Jd) Verify(configParamMap map[string]string, query, methodCode string) (verifykRsp common.VerifyRsp, errCode int, err error) {
	//verifykRsp.EncryptRsp = query

	//解析参数
	urlValuesMap, err := url.ParseQuery(query)
	if err != nil {
		logrus.Errorf("org:jd,"+CallbackEncryptFormatErrMessage+",query:%v,errCode:%v,err:%v", query, CallbackEncryptFormatErrCode, err.Error())
		return verifykRsp, CallbackEncryptFormatErrCode, errors.New(CallbackEncryptFormatErrMessage)
	}
	queryMap := make(map[string]string)
	for k, v := range urlValuesMap {
		queryMap[k] = v[0]
	}

	//解密
	decryptMap, err := decryptVerifyArg(queryMap, configParamMap["des_key"])
	if err != nil {
		logrus.Errorf("org:jd,"+CallbackDecryptFailedErrMessage+",query:%v,errCode:%v,err:%v", query, CallbackDecryptFailedErrCode, err.Error())
		return verifykRsp, CallbackDecryptFailedErrCode, errors.New(CallbackDecryptFailedErrMessage)
	}
	decryptBytes, err := json.Marshal(decryptMap)
	if err != nil {
		logrus.Errorf("org:jd,"+CallbackDecryptFormatErrMessage+",query:%v,errCode:%v,err:%v", query, CallbackDecryptFormatErrCode, err.Error())
		return verifykRsp, CallbackDecryptFormatErrCode, errors.New(CallbackDecryptFormatErrMessage)
	}
	//verifykRsp.DecryptRsp = string(decryptBytes)
	fmt.Println("decryptBytes", string(decryptBytes))

	//解析为结构体
	var verifyQuery VerifyQuery
	err = json.Unmarshal(decryptBytes, &verifyQuery)
	if err != nil {
		logrus.Errorf("org:jd,"+CallbackDecryptFormatErrMessage+",query:%v,errCode:%v,err:%v", query, CallbackDecryptFormatErrCode, err.Error())
		return verifykRsp, CallbackDecryptFormatErrCode, errors.New(CallbackDecryptFormatErrMessage)
	}
	verifykRsp.OrderId = verifyQuery.TradeNum

	//校验签名
	if !checkVerifySign(decryptMap, configParamMap["public_key"]) {
		logrus.Errorf("org:jd,"+CallbackSignErrMessage+",query:%v,errCode:%v,err:%v", query, CallbackSignErrCode)
		return verifykRsp, CallbackSignErrCode, errors.New(CallbackSignErrMessage)
	}

	//交易状态
	if verifyQuery.Status != CallbackSuccessCode {
		logrus.Errorf("org:jd,"+CallbackStatusErrMessage+",query:%v,errCode:%v,err:%v", query, CallbackStatusErrCode)
		return verifykRsp, CallbackStatusErrCode, errors.New(CallbackStatusErrMessage)
	}
	verifykRsp.Status = true

	return verifykRsp, 0, nil
}

//解密
func decryptVerifyArg(encryptMap map[string]string, desKey string) (decryptMap map[string]string, err error) {
	//解密
	desKeyBytes, err := base64.StdEncoding.DecodeString(desKey)
	if err != nil {
		return decryptMap, err
	}

	//遍历map解密
	decryptMap = make(map[string]string)
	for k, v := range encryptMap {
		if k == "sign" || v == "" {
			decryptMap[k] = v
			continue
		}
		encryptBytes, err := HexString2Bytes(v)
		decryptBytes, err := TripleEcbDesDecrypt(encryptBytes, desKeyBytes)
		if err != nil {
			return nil, err
		}
		decryptMap[k] = string(decryptBytes)
	}

	return decryptMap, nil
}

//校验签名
func checkVerifySign(urlValuesMap map[string]string, publicKey string) bool {
	sign, ok := urlValuesMap["sign"]
	if !ok {
		return false
	}
	if sign == "" {
		return false
	}
	delete(urlValuesMap, "sign")
	encodePayString := GetNotEmptySortString(urlValuesMap)

	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false
	}

	sha256 := HaSha256(encodePayString)

	return VerifyPKCS1v15([]byte(sha256), signBytes, []byte(publicKey), crypto.Hash(0))
}

func (jd *Jd) GetVerifyConfigCode() []string {
	return []string{
		"merchant",
		"des_key", "public_key",
	}
}
