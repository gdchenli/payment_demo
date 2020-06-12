package allpay

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gdchenli/pay/dialects/allpay/util"
)

const (
	NotifyQueryFormatErrCode    = 10201
	NotifyQueryFormatErrMessage = "异步通知，支付数据格式错误"
	NotifySignErrCode           = 10205
	NotifySignErrMessage        = "异步通知，签名校验失败"
	NotifyQueryRateErrCode      = 10220
	NotifyQueryRateErrMessage   = "异步通知，汇率查询失败"
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

type NotifyArg struct {
	MerId       string `json:"mer_id"`
	Md5Key      string `json:"md5_key"`
	SapiGateWay string `json:"sapi_gate_way"`
}

func (notify *Notify) Validate(query string, arg NotifyArg) (notifyRsp NotifyRsp, errCode int, err error) {
	notifyRsp.Rsp = query

	//解析参数
	queryMap, err := util.JsonToMap(query)
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
	if !notify.checkSign(queryMap, arg.Md5Key, sign) {
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

	//汇率
	rateArg := RateArg{
		MerId:                  arg.MerId,
		OriginalCurrencyCode:   queryMap["orderCurrency"],
		ConversionCurrencyCode: CNY,
		Md5Key:                 arg.Md5Key,
		PaymentSchema:          queryMap["paymentSchema"],
		GateWay:                arg.SapiGateWay,
	}
	notifyRsp.Rate, errCode, err = new(Rate).Search(rateArg)
	if err != nil {
		logrus.Errorf("org:allpay,"+NotifyQueryRateErrMessage+",errCode:%v,err:%v", NotifyQueryRateErrCode, err.Error())
		return notifyRsp, NotifyQueryRateErrCode, errors.New(NotifyQueryRateErrMessage)
	}

	//订单外币金额
	orderAmount, err := strconv.ParseFloat(queryMap["orderAmount"], 64)
	if err != nil {
		logrus.Errorf("org:allpay,"+NotifyQueryFormatErrMessage+",errCode:%v,err:%v", NotifyQueryFormatErrCode, err.Error())
		return notifyRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
	}

	//人民币金额
	notifyRsp.RmbFee, err = strconv.ParseFloat(fmt.Sprintf("%.2f", orderAmount*notifyRsp.Rate), 64)
	if err != nil {
		logrus.Errorf("org:allpay,"+NotifyQueryFormatErrMessage+",errCode:%v,err:%v", NotifyQueryFormatErrCode, err.Error())
		return notifyRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
	}

	return notifyRsp, 0, nil
}

func (notify *Notify) checkSign(queryMap map[string]string, md5Key, sign string) bool {
	sortString := util.GetSortString(queryMap)
	calculateSign := util.Md5(sortString + md5Key)
	return calculateSign == sign
}
