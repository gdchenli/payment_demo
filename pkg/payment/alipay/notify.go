package alipay

import (
	"errors"
	"fmt"
	"payment_demo/api/response"
	"payment_demo/internal/common/request"
	"strconv"

	"github.com/sirupsen/logrus"
)

const (
	NotifyQueryFormatErrCode      = 10201
	NotifyQueryFormatErrMessage   = "异步通知，支付数据格式错误"
	NotifyDecryptFormatErrCode    = 10203
	NotifyDecryptFormatErrMessage = "异步通知，解密后数据格式错误"
	NotifySignErrCode             = 10205
	NotifySignErrMessage          = "异步通知，签名校验失败"
)

type Notify struct{}

func (notify *Notify) Validate(configParamMap map[string]string, query, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error) {
	//解析参数
	queryMap, err := ParseQueryString(query)
	if err != nil {
		return notifyRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
	}

	//订单编号
	notifyRsp.OrderId = queryMap["out_trade_no"]

	//校验签名
	var sign string
	if value, ok := queryMap["sign"]; ok {
		sign = value
		delete(queryMap, "sign")
		delete(queryMap, "sign_type")
	}

	if !notify.checkSign(queryMap, configParamMap["md5_key"], sign) {
		logrus.Errorf(NotifySignErrMessage+",order id %v,errCode:%v", notifyRsp.OrderId, NotifySignErrCode)
		return notifyRsp, NotifySignErrCode, errors.New(NotifySignErrMessage)
	}

	//交易状态
	if queryMap["trade_status"] == TradeFinished || queryMap["trade_status"] == TradeSuccess {
		notifyRsp.Status = true
	}

	//alipay交易流水号，
	notifyRsp.TradeNo = queryMap["trade_no"]

	totalFee, err := strconv.ParseFloat(queryMap["total_fee"], 64)
	if err != nil {
		logrus.Errorf("org:alipay,"+NotifyDecryptFormatErrMessage+",order id %v,errCode:%v,err:%v", notifyRsp.OrderId, NotifyDecryptFormatErrCode, err.Error())
		return notifyRsp, NotifyDecryptFormatErrCode, errors.New(NotifyDecryptFormatErrMessage)
	}

	tradeArg := request.SearchTradeReq{
		OrderId:    queryMap["out_trade_no"],
		MethodCode: methodCode,
		OrgCode:    "alipay", TotalFee: totalFee,
		Currency: queryMap["currency"]}
	alipayTradeRsp, errCode, err := new(Trade).Search(configParamMap, tradeArg)
	if err != nil {
		return notifyRsp, errCode, err
	}

	//人民币金额
	notifyRsp.RmbFee = alipayTradeRsp.RmbFee

	//支付时间
	notifyRsp.PaidAt = alipayTradeRsp.PaidAt

	//汇率
	rate := alipayTradeRsp.RmbFee / totalFee
	rate, err = strconv.ParseFloat(fmt.Sprintf("%.8f", rate), 64)
	if err != nil {
		logrus.Errorf("org:alipay,"+NotifyDecryptFormatErrMessage+",order id %v,errCode:%v,err:%v", notifyRsp.OrderId, NotifyDecryptFormatErrCode, err.Error())
		return notifyRsp, NotifySignErrCode, errors.New(NotifySignErrMessage)
	}
	notifyRsp.Rate = rate
	notifyRsp.Message = "success"

	return notifyRsp, 0, nil
}

func (notify *Notify) checkSign(queryMap map[string]string, md5Key, sign string) bool {
	sortString := GetSortString(queryMap)
	calculateSign := Md5(sortString + md5Key)

	return calculateSign == sign
}

func (notify *Notify) GetConfigCode() []string {
	return []string{
		"md5_key", "gate_way", "partner",
	}
}
