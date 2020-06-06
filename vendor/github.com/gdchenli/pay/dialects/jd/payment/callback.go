package payment

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/gdchenli/pay/dialects/jd/util"
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

type Callback struct{}

type CallbackArg struct {
	PublicKey string `json:"public_key"` //公钥
	DesKey    string `json:"des_key"`
}

type CallbackQuery struct {
	TradeNum  string `json:"tradeNum"`  //订单号
	Amount    string `json:"amount"`    //交易金额
	Currency  string `json:"currency"`  //货币类型
	TradeTime string `json:"tradeTime"` //交易时间
	Status    string `json:"status"`    //交易状态
	Sign      string `json:"sign"`      //签名
}

type CallbackRsp struct {
	OrderId    string `json:"order_id"`    //订单号
	Status     bool   `json:"status"`      //交易状态，true交易成功 false交易失败
	EncryptRsp string `json:"encrypt_rsp"` //返回的加密数据
	DecryptRsp string `json:"decrypt_rsp"` //返回的解密数据
}

func (callback *Callback) Validate(query string, arg CallbackArg) (callbackRsp CallbackRsp, errCode int, err error) {
	callbackRsp.EncryptRsp = query

	//解析参数
	urlValuesMap, err := url.ParseQuery(query)
	if err != nil {
		logrus.Errorf(CallbackEncryptFormatErrMessage+",query:%v,errCode:%v,err:%v", query, CallbackEncryptFormatErrCode, err.Error())
		return callbackRsp, CallbackEncryptFormatErrCode, errors.New(CallbackEncryptFormatErrMessage)
	}
	queryMap := make(map[string]string)
	for k, v := range urlValuesMap {
		queryMap[k] = v[0]
	}

	//解密
	decryptMap, err := callback.decryptArg(queryMap, arg.DesKey)
	if err != nil {
		logrus.Errorf(CallbackDecryptFailedErrMessage+",query:%v,errCode:%v,err:%v", query, CallbackDecryptFailedErrCode, err.Error())
		return callbackRsp, CallbackDecryptFailedErrCode, errors.New(CallbackDecryptFailedErrMessage)
	}
	decryptBytes, err := json.Marshal(decryptMap)
	if err != nil {
		logrus.Errorf(CallbackDecryptFormatErrMessage+",query:%v,errCode:%v,err:%v", query, CallbackDecryptFormatErrCode, err.Error())
		return callbackRsp, CallbackDecryptFormatErrCode, errors.New(CallbackDecryptFormatErrMessage)
	}
	callbackRsp.DecryptRsp = string(decryptBytes)
	fmt.Println("decryptBytes", string(decryptBytes))

	//解析为结构体
	var callbackQuery CallbackQuery
	err = json.Unmarshal(decryptBytes, &callbackQuery)
	if err != nil {
		logrus.Errorf(CallbackDecryptFormatErrMessage+",query:%v,errCode:%v,err:%v", query, CallbackDecryptFormatErrCode, err.Error())
		return callbackRsp, CallbackDecryptFormatErrCode, errors.New(CallbackDecryptFormatErrMessage)
	}
	callbackRsp.OrderId = callbackQuery.TradeNum

	//校验签名
	if !callback.checkSign(decryptMap, arg.PublicKey) {
		logrus.Errorf(CallbackSignErrMessage+",query:%v,errCode:%v,err:%v", query, CallbackSignErrCode)
		return callbackRsp, CallbackSignErrCode, errors.New(CallbackSignErrMessage)
	}

	//交易状态
	if callbackQuery.Status != CallbackSuccessCode {
		logrus.Errorf(CallbackStatusErrMessage+",query:%v,errCode:%v,err:%v", query, CallbackStatusErrCode)
		return callbackRsp, CallbackStatusErrCode, errors.New(CallbackStatusErrMessage)
	}
	callbackRsp.Status = true

	return callbackRsp, 0, nil
}

//解密
func (callback *Callback) decryptArg(encryptMap map[string]string, desKey string) (decryptMap map[string]string, err error) {
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
		encryptBytes, err := util.HexString2Bytes(v)
		decryptBytes, err := util.TripleEcbDesDecrypt(encryptBytes, desKeyBytes)
		if err != nil {
			return nil, err
		}
		decryptMap[k] = string(decryptBytes)
	}

	return decryptMap, nil
}

//校验签名
func (callback *Callback) checkSign(urlValuesMap map[string]string, publicKey string) bool {
	sign, ok := urlValuesMap["sign"]
	if !ok {
		return false
	}
	if sign == "" {
		return false
	}
	delete(urlValuesMap, "sign")
	encodePayString := util.GetNotEmptySortString(urlValuesMap)

	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false
	}

	sha256 := util.HaSha256(encodePayString)

	return util.VerifyPKCS1v15([]byte(sha256), signBytes, []byte(publicKey), crypto.Hash(0))
}
