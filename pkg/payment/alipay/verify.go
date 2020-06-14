package alipay

import (
	"errors"
	"payment_demo/api/response"

	"github.com/sirupsen/logrus"

	"github.com/gdchenli/pay/dialects/alipay/util"
)

type Verify struct{}

type VerifyRsp struct {
	OrderId string `json:"order_id"` //订单号
	Status  bool   `json:"status"`   //交易状态，true交易成功 false交易失败
	Rsp     string `json:"rsp"`      //返回的数据
}

func (vreify *Verify) Validate(configParamMap map[string]string, query, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) {
	//callbackRsp.Rsp = query

	//解析参数
	queryMap, err := util.ParseQueryString(query)
	if err != nil {
		logrus.Errorf("org:alipay,"+NotifyQueryFormatErrMessage+",errCode:%v,err:%v", NotifyQueryFormatErrCode, err.Error())
		return verifyRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
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

	if !vreify.checkSign(queryMap, configParamMap["md5_key"], sign) {
		logrus.Errorf("org:alipay,"+NotifySignErrMessage+",orderId:%v,errCode:%v", queryMap["out_trade_no"], NotifySignErrCode)
		return verifyRsp, NotifySignErrCode, errors.New(NotifySignErrMessage)
	}

	//交易状态
	if queryMap["trade_status"] == TradeFinished || queryMap["trade_status"] == TradeSuccess {
		verifyRsp.Status = true
	}

	return verifyRsp, 0, nil
}

func (vreify *Verify) checkSign(queryMap map[string]string, signKey, sign string) bool {
	sortString := util.GetSortString(queryMap)
	calculateSign := util.Md5(sortString + signKey)
	return calculateSign == sign
}

func (vreify *Verify) GetConfigCode() []string {
	return []string{
		"pay_way", "md5_key", "gate_way", "partner",
	}
}
