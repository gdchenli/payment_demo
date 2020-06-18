package epayments

import (
	"errors"
	"fmt"
	"payment_demo/api/response"

	"github.com/sirupsen/logrus"
)

const (
	VerifyQueryFormatErrCode    = 10301
	VerifyQueryFormatErrMessage = "同步通知，支付数据格式错误"
	VerifySignErrCode           = 10305
	VerifySignErrMessage        = "同步通知，签名校验失败"
)

type CallbackRsp struct {
	OrderId string `json:"order_id"` //订单号
	Status  bool   `json:"status"`   //交易状态，true交易成功 false交易失败
	Rsp     string `json:"rsp"`      //返回的数据
}

func (epayments *Epayments) Verify(configParamMap map[string]string, query, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) {
	//解析参数
	queryMap, err := ParseQueryString(query)
	if err != nil {
		logrus.Errorf("org:epayments,"+VerifyQueryFormatErrMessage+",query:%v,errCode:%v,err:%v", query, VerifyQueryFormatErrCode, err.Error())
		return verifyRsp, VerifyQueryFormatErrCode, errors.New(VerifyQueryFormatErrMessage)
	}

	//订单编号
	verifyRsp.OrderId = queryMap["increment_id"]

	//校验签名
	var sign string
	if value, ok := queryMap["signature"]; ok {
		sign = value
		delete(queryMap, "signature")
		delete(queryMap, "sign_type")
	}

	if !checkVerifySign(queryMap, configParamMap["md5_key"], sign) {
		logrus.Errorf("org:epayments,"+VerifySignErrMessage+",query:%v,errCode:%v", query, VerifySignErrCode)
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
	fmt.Println("sortString", sortString)
	calculateSign := Md5(sortString + signKey)
	return calculateSign == sign
}

func (epayments *Epayments) GetVerifyConfigCode() []string {
	return []string{
		"md5_key",
	}
}
