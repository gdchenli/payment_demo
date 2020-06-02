package payment

import (
	"errors"
	"payment_demo/pkg/epayments/util"
)

const (
	NotifyQueryFormatErrCode    = 10201
	NotifyQueryFormatErrMessage = "异步通知，支付数据格式错误"
	NotifySignErrCode           = 10205
	NotifySignErrMessage        = "异步通知，签名校验失败"
)

type Notify struct{}

type NotifyRsp struct {
	OrderId string `json:"order_id"` //订单号
	Status  bool   `json:"status"`   //交易状态，true交易成功 false交易失败
	TradeNo string `json:"trade_no"` //支付机构交易流水号
	Rsp     string `json:"rsp"`      //返回的数据
}

func (notify *Notify) Validate(query, md5Key string) (notifyRsp NotifyRsp, errCode int, err error) {
	notifyRsp.Rsp = query

	//解析参数
	queryMap, err := util.ParseQueryString(query)
	if err != nil {
		return notifyRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
	}

	//订单编号
	notifyRsp.OrderId = queryMap["increment_id"]

	//校验签名
	var sign string
	if value, ok := queryMap["signature"]; ok {
		sign = value
		delete(queryMap, "signature")
		delete(queryMap, "sign_type")
	}

	if !notify.checkSign(queryMap, md5Key, sign) {
		return notifyRsp, NotifySignErrCode, errors.New(NotifySignErrMessage)
	}

	//交易状态
	if queryMap["trade_status"] == TradeFinished || queryMap["trade_status"] == TradeSuccess {
		notifyRsp.Status = true
	}

	//alipay交易流水号，
	notifyRsp.TradeNo = queryMap["trade_no"]

	return notifyRsp, 0, nil
}

func (notify *Notify) checkSign(queryMap map[string]string, md5Key, sign string) bool {
	sortString := util.GetSortString(queryMap)
	calculateSign := util.Md5(sortString + md5Key)
	return calculateSign == sign
}
