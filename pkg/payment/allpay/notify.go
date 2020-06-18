package allpay

import (
	"errors"
	"fmt"
	"payment_demo/api/response"
	"payment_demo/internal/common/request"
	"payment_demo/pkg/payment/consts"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	NotifyQueryFormatErrCode    = 10201
	NotifyQueryFormatErrMessage = "异步通知，支付数据格式错误"
	NotifySignErrCode           = 10205
	NotifySignErrMessage        = "异步通知，签名校验失败"
	NotifyQueryRateErrCode      = 10220
	NotifyQueryRateErrMessage   = "异步通知，汇率查询失败"
)

func (allpay *Allpay) Notify(configParamMap map[string]string, query, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error) {
	//解析参数
	queryMap, err := JsonToMap(query)
	if err != nil {
		logrus.Errorf("org:allpay,"+NotifyQueryFormatErrMessage+",errCode:%v,err:%v", NotifyQueryFormatErrCode, err.Error())
		return notifyRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
	}

	//订单编号
	notifyRsp.OrderId = queryMap["orderNum"]

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
	if !checkNotifySign(queryMap, configParamMap["md5_key"], sign) {
		logrus.Errorf("org:allpay,"+NotifySignErrMessage+",order id %v,errCode:%v", notifyRsp.OrderId, NotifySignErrCode)
		return notifyRsp, NotifySignErrCode, errors.New(NotifySignErrMessage)
	}

	//交易状态
	if queryMap["RespCode"] == "00" {
		notifyRsp.Status = true
	}

	//allpay交易流水号
	notifyRsp.TradeNo = queryMap["transID"]

	//allpay不返回支付成功时间，取当前异步通知时的服务器时间
	notifyRsp.PaidAt = time.Now().UTC().Format(DateTimeFormatLayout)

	//订单外币金额
	orderAmount, err := strconv.ParseFloat(queryMap["orderAmount"], 64)
	if err != nil {
		logrus.Errorf("org:allpay,"+NotifyQueryFormatErrMessage+",errCode:%v,err:%v", NotifyQueryFormatErrCode, err.Error())
		return notifyRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
	}

	//汇率
	rateArg := request.SearchTradeReq{
		OrderId:    queryMap["orderNum"],
		MethodCode: methodCode,
		OrgCode:    consts.AllpayOrg,
		Currency:   queryMap["orderCurrency"],
		TotalFee:   orderAmount,
	}
	trade, errCode, err := allpay.SearchTrade(configParamMap, rateArg)
	if err != nil {
		logrus.Errorf("org:allpay,"+NotifyQueryRateErrMessage+",errCode:%v,err:%v", NotifyQueryRateErrCode, err.Error())
		return notifyRsp, NotifyQueryRateErrCode, errors.New(NotifyQueryRateErrMessage)
	}
	notifyRsp.Rate = trade.Rate

	//人民币金额
	notifyRsp.RmbFee, err = strconv.ParseFloat(fmt.Sprintf("%.2f", orderAmount*notifyRsp.Rate), 64)
	if err != nil {
		logrus.Errorf("org:allpay,"+NotifyQueryFormatErrMessage+",errCode:%v,err:%v", NotifyQueryFormatErrCode, err.Error())
		return notifyRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
	}
	notifyRsp.Message = "Ok"
	return notifyRsp, 0, nil
}

func checkNotifySign(queryMap map[string]string, md5Key, sign string) bool {
	sortString := GetSortString(queryMap)
	calculateSign := Md5(sortString + md5Key)
	return calculateSign == sign
}

func (allpay *Allpay) GetNotifyConfigCode() []string {
	return []string{
		"merID", "acqID", "md5_key", "sapi_way",
	}
}
