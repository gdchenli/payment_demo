package allpay

import (
	"errors"
	"payment_demo/api/payment/response"

	"github.com/sirupsen/logrus"
)

const (
	VerifyQueryFormatErrCode    = 10301
	VerifyQueryFormatErrMessage = "同步通知，支付数据格式错误"
	VerifySignErrCode           = 10305
	VerifySignErrMessage        = "同步通知，签名校验失败"
)

func (allpay *Allpay) Verify(configParamMap map[string]string, query, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) {
	//callbackRsp.Rsp = query

	//解析参数
	queryMap, err := ParseQueryString(query)
	if err != nil {
		logrus.Errorf("org:allpay,"+VerifyQueryFormatErrMessage+",errCode:%v,err:%v", VerifyQueryFormatErrCode, err.Error())
		return verifyRsp, VerifyQueryFormatErrCode, errors.New(VerifyQueryFormatErrMessage)
	}

	//订单编号
	verifyRsp.OrderId = queryMap["orderNum"]

	//校验签名
	var sign string
	if value, ok := queryMap["sign"]; ok {
		sign = value
		delete(queryMap, "sign")
	}
	if value, ok := queryMap["signature"]; ok {
		sign = value
		delete(queryMap, "signature")
	}
	if !checkVerifySign(queryMap, configParamMap["md5_key"], sign) {
		logrus.Errorf("org:allpay,"+VerifySignErrMessage+",order id %v,errCode:%v", queryMap["orderNum"], VerifySignErrCode)
		return verifyRsp, VerifySignErrCode, errors.New(VerifySignErrMessage)
	}

	//交易状态
	if queryMap["RespCode"] == "00" {
		verifyRsp.Status = true
	}

	return verifyRsp, 0, nil
}

func checkVerifySign(queryMap map[string]string, signKey, sign string) bool {
	sortString := GetSortString(queryMap)
	calculateSign := Md5(sortString + signKey)
	return calculateSign == sign
}

func (allpay *Allpay) GetVerifyConfigCode() []string {
	return []string{
		"md5_key",
	}
}
