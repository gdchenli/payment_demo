package payment

import (
	"errors"
	"payment_demo/pkg/alipay/util"
)

type Callback struct{}

type CallbackRsp struct {
	OrderId string `json:"order_id"` //订单号
	Status  bool   `json:"status"`   //交易状态，true交易成功 false交易失败
	Rsp     string `json:"rsp"`      //返回的数据
}

func (callback *Callback) Validate(query, md5Key string) (callbackRsp CallbackRsp, errCode int, err error) {
	callbackRsp.Rsp = query

	//解析参数
	queryMap, err := util.ParseQueryString(query)
	if err != nil {
		return callbackRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
	}

	//订单编号
	callbackRsp.OrderId = queryMap["out_trade_no"]

	//校验签名
	var sign string
	if value, ok := queryMap["sign"]; ok {
		sign = value
		delete(queryMap, "sign")
		delete(queryMap, "sign_type")
	}

	if !callback.checkSign(queryMap, md5Key, sign) {
		return callbackRsp, NotifySignErrCode, errors.New(NotifySignErrMessage)
	}

	//交易状态
	if queryMap["trade_status"] == TradeFinished || queryMap["trade_status"] == TradeSuccess {
		callbackRsp.Status = true
	}

	return callbackRsp, 0, nil
}

func (callback *Callback) checkSign(queryMap map[string]string, signKey, sign string) bool {
	sortString := util.GetSortString(queryMap)
	calculateSign := util.Md5(sortString + signKey)
	return calculateSign == sign
}
