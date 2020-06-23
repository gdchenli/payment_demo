package alipay

import (
	"errors"
	"payment_demo/api/notice/response"

	"github.com/sirupsen/logrus"
)

const (
	VerifyQueryFormatErrCode    = 10301
	VerifyQueryFormatErrMessage = "同步通知，支付数据格式错误"
	VerifySignErrCode           = 10305
	VerifySignErrMessage        = "同步通知，签名校验失败"
)

func (alipay *Alipay) Verify(configParamMap map[string]string, query, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) {
	//callbackRsp.Rsp = query

	//解析参数
	queryMap, err := ParseQueryString(query)
	if err != nil {
		logrus.Errorf("org:alipay,"+VerifyQueryFormatErrMessage+",errCode:%v,err:%v", VerifyQueryFormatErrCode, err.Error())
		return verifyRsp, VerifyQueryFormatErrCode, errors.New(VerifyQueryFormatErrMessage)
	}

	//订单编号
	verifyRsp.OrderId = queryMap["out_trade_no"]

	//校验签名
	var sign string
	if value, ok := queryMap["sign"]; ok {
		sign = value
		delete(queryMap, "sign")
		delete(queryMap, "sign_type")
	}

	if !checkVerifySign(queryMap, configParamMap["md5_key"], sign) {
		logrus.Errorf("org:alipay,"+VerifySignErrMessage+",orderId:%v,errCode:%v", queryMap["out_trade_no"], VerifySignErrCode)
		return verifyRsp, VerifySignErrCode, errors.New(VerifySignErrMessage)
	}

	//交易状态
	if queryMap["trade_status"] == TradeFinished || queryMap["trade_status"] == TradeSuccess {
		verifyRsp.Status = true
	}

	return verifyRsp, 0, nil
}

func checkVerifySign(queryMap map[string]string, signKey, sign string) bool {
	sortString := GetSortString(queryMap)
	calculateSign := Md5(sortString + signKey)
	return calculateSign == sign
}

func (alipay *Alipay) GetVerifyConfigCode() []string {
	return []string{
		"md5_key",
	}
}
