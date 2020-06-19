package epayments

import (
	"errors"
	"fmt"
	"payment_demo/api/payment/response"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	NotifyQueryFormatErrCode    = 10201
	NotifyQueryFormatErrMessage = "异步通知，支付数据格式错误"
	NotifySignErrCode           = 10205
	NotifySignErrMessage        = "异步通知，签名校验失败"
)

type NotifyRsp struct {
	OrderId string  `json:"order_id"` //订单号
	Status  bool    `json:"status"`   //交易状态，true交易成功 false交易失败
	TradeNo string  `json:"trade_no"` //支付机构交易流水号
	PaidAt  string  `json:"paid_at"`  //支付gmt时间
	RmbFee  float64 `json:"rmb_fee"`  //人民币金额
	Rate    float64 `json:"rate"`     //汇率
	Rsp     string  `json:"rsp"`      //返回的数据
}

func (epayments *Epayments) Notify(configParamMap map[string]string, query, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error) {
	//解析参数
	queryMap, err := ParseQueryString(query)
	if err != nil {
		logrus.Errorf(NotifyQueryFormatErrMessage+",query:%v,errCode:%v,err:%v", query, NotifyQueryFormatErrCode, err.Error())
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

	if !checkNotifySign(queryMap, configParamMap["md5_key"], sign) {
		logrus.Errorf("org:epayments,"+NotifySignErrMessage+",orderId%v,query:%v,errCode:%v", notifyRsp.OrderId, query, NotifySignErrCode)
		return notifyRsp, NotifySignErrCode, errors.New(NotifySignErrMessage)
	}

	//交易状态
	if queryMap["trade_status"] == TradeFinished || queryMap["trade_status"] == TradeSuccess {
		notifyRsp.Status = true
	}

	//交易流水号
	notifyRsp.TradeNo = queryMap["trade_no"]

	//支付时间
	parseTime, err := time.Parse(DateTimeFormatLayout, queryMap["gmt_payment"])
	if err != nil {
		logrus.Errorf("org:epayments,"+NotifyQueryFormatErrMessage+",orderId:%v,query:%v,errCode:%v,err:%v", notifyRsp.OrderId, query, NotifyQueryFormatErrCode, err.Error())
		return notifyRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
	}
	notifyRsp.PaidAt = parseTime.UTC().Format(DateTimeFormatLayout)

	//汇率
	notifyRsp.Rate, err = strconv.ParseFloat(queryMap["rate"], 64)
	if err != nil {
		logrus.Errorf("org:epayments,"+NotifyQueryFormatErrMessage+",orderId:%v,query:%v,errCode:%v,err:%v", notifyRsp.OrderId, query, NotifyQueryFormatErrCode, err.Error())
		return notifyRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
	}

	//人民币金额
	grandTotal, err := strconv.ParseFloat(queryMap["grandtotal"], 64)
	if err != nil {
		logrus.Errorf("org:epayments,"+NotifyQueryFormatErrMessage+",orderId:%v,query:%v,errCode:%v,err:%v", notifyRsp.OrderId, query, NotifyQueryFormatErrCode, err.Error())
		return notifyRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
	}
	rmbFee := grandTotal * notifyRsp.Rate
	rmbFee, err = strconv.ParseFloat(fmt.Sprintf("%.2f", rmbFee), 64)
	if err != nil {
		logrus.Errorf("org:epayments,"+NotifyQueryFormatErrMessage+",orderId:%v,query:%v,errCode:%v,err:%v", notifyRsp.OrderId, query, NotifyQueryFormatErrCode, err.Error())
		return notifyRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
	}
	notifyRsp.RmbFee = rmbFee
	notifyRsp.Message = "success"

	return notifyRsp, 0, nil
}

func checkNotifySign(queryMap map[string]string, md5Key, sign string) bool {
	sortString := GetSortString(queryMap)
	calculateSign := Md5(sortString + md5Key)
	return calculateSign == sign
}

func (epayments *Epayments) GetNotifyConfigCode() []string {
	return []string{
		"md5_key",
	}
}
