package payment

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/gdchenli/pay/dialects/alipay/util"
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
	Merchant string `json:"merchant"`
	Md5Key   string `json:"md5_key"`
	GateWay  string `json:"gate_way"`
}

func (notify *Notify) Validate(query string, arg NotifyArg) (notifyRsp NotifyRsp, errCode int, err error) {
	notifyRsp.Rsp = query

	//解析参数
	queryMap, err := util.ParseQueryString(query)
	if err != nil {
		return notifyRsp, NotifyQueryFormatErrCode, errors.New(NotifyQueryFormatErrMessage)
	}
	fmt.Printf("%+v\n", queryMap)

	//订单编号
	notifyRsp.OrderId = queryMap["out_trade_no"]

	//校验签名
	var sign string
	if value, ok := queryMap["sign"]; ok {
		sign = value
		delete(queryMap, "sign")
		delete(queryMap, "sign_type")
	}

	if !notify.checkSign(queryMap, arg.Md5Key, sign) {
		logrus.Errorf(NotifySignErrMessage+",order id %v,errCode:%v", notifyRsp.OrderId, NotifySignErrCode)
		return notifyRsp, NotifySignErrCode, errors.New(NotifySignErrMessage)
	}

	//交易状态
	if queryMap["trade_status"] == TradeFinished || queryMap["trade_status"] == TradeSuccess {
		notifyRsp.Status = true
	}

	//alipay交易流水号，
	notifyRsp.TradeNo = queryMap["trade_no"]

	tradeArg := TradeArg{
		Merchant:   arg.Merchant,
		OutTradeNo: queryMap["out_trade_no"],
		Md5Key:     arg.Md5Key,
		TradeWay:   arg.GateWay,
	}
	alipayTradeRsp, errCode, err := new(Trade).Search(tradeArg)
	if err != nil {
		return notifyRsp, errCode, err
	}

	//人民币金额
	notifyRsp.RmbFee = alipayTradeRsp.RmbFee

	//支付时间
	notifyRsp.PaidAt = alipayTradeRsp.PaidAt

	//汇率
	totalFee, err := strconv.ParseFloat(queryMap["total_fee"], 64)
	if err != nil {
		logrus.Errorf("org:alipay,"+NotifyDecryptFormatErrMessage+",order id %v,errCode:%v,err:%v", notifyRsp.OrderId, NotifyDecryptFormatErrCode, err.Error())
		return notifyRsp, NotifyDecryptFormatErrCode, errors.New(NotifyDecryptFormatErrMessage)
	}
	rate := alipayTradeRsp.RmbFee / totalFee
	rate, err = strconv.ParseFloat(fmt.Sprintf("%.8f", rate), 64)
	if err != nil {
		logrus.Errorf("org:alipay,"+NotifyDecryptFormatErrMessage+",order id %v,errCode:%v,err:%v", notifyRsp.OrderId, NotifyDecryptFormatErrCode, err.Error())
		return notifyRsp, NotifySignErrCode, errors.New(NotifySignErrMessage)
	}
	notifyRsp.Rate = rate

	return notifyRsp, 0, nil
}

func (notify *Notify) checkSign(queryMap map[string]string, md5Key, sign string) bool {
	sortString := util.GetSortString(queryMap)
	fmt.Println("sortString", sortString)
	calculateSign := util.Md5(sortString + md5Key)
	fmt.Println("calculateSign", calculateSign)
	return calculateSign == sign
}
