package alipay

import (
	"payment_demo/api/response"

	"github.com/gdchenli/pay/dialects/alipay/util"
)

type Verify struct{}

type VerifyRsp struct {
	OrderId string `json:"order_id"` //订单号
	Status  bool   `json:"status"`   //交易状态，true交易成功 false交易失败
	Rsp     string `json:"rsp"`      //返回的数据
}

func (vreify *Verify) Validate(query, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) {
	/*callbackRsp.Rsp = query

	//解析参数
	queryMap, err := util.ParseQueryString(query)
	if err != nil {
		logrus.Errorf("org:alipay,"+NotifyQueryFormatErrMessage+",errCode:%v,err:%v", NotifyQueryFormatErrCode, err.Error())
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
		logrus.Errorf("org:alipay,"+NotifySignErrMessage+",orderId:%v,errCode:%v", callbackRsp.OrderId, NotifySignErrCode)
		return callbackRsp, NotifySignErrCode, errors.New(NotifySignErrMessage)
	}

	//交易状态
	if queryMap["trade_status"] == TradeFinished || queryMap["trade_status"] == TradeSuccess {
		callbackRsp.Status = true
	}*/

	return verifyRsp, 0, nil
}

func (vreify *Verify) checkSign(queryMap map[string]string, signKey, sign string) bool {
	sortString := util.GetSortString(queryMap)
	calculateSign := util.Md5(sortString + signKey)
	return calculateSign == sign
}
