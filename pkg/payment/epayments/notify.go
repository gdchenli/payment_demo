package epayments

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gdchenli/pay/dialects/epayments/util"
	"github.com/sirupsen/logrus"
)

const (
	NotifyQueryFormatErrCode    = 10201
	NotifyQueryFormatErrMessage = "异步通知，支付数据格式错误"
	NotifySignErrCode           = 10205
	NotifySignErrMessage        = "异步通知，签名校验失败"
)

type Notify struct{}

type NotifyRsp struct {
	OrderId string  `json:"order_id"` //订单号
	Status  bool    `json:"status"`   //交易状态，true交易成功 false交易失败
	TradeNo string  `json:"trade_no"` //支付机构交易流水号
	PaidAt  string  `json:"paid_at"`  //支付gmt时间
	RmbFee  float64 `json:"rmb_fee"`  //人民币金额
	Rate    float64 `json:"rate"`     //汇率
	Rsp     string  `json:"rsp"`      //返回的数据
}

func (notify *Notify) Validate(query, md5Key string) (notifyRsp NotifyRsp, errCode int, err error) {
	notifyRsp.Rsp = query

	//解析参数
	queryMap, err := util.ParseQueryString(query)
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

	if !notify.checkSign(queryMap, md5Key, sign) {
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

	return notifyRsp, 0, nil
}

func (notify *Notify) checkSign(queryMap map[string]string, md5Key, sign string) bool {
	sortString := util.GetSortString(queryMap)
	calculateSign := util.Md5(sortString + md5Key)
	return calculateSign == sign
}
