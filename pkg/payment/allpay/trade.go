package allpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gdchenli/pay/dialects/allpay/util"
	"github.com/gdchenli/pay/pkg/curl"
)

const (
	SearchTradeTransType = "INQY"
	TradeRoute           = "/api/orderquery"
)
const (
	SearchTradeWait    = "1" //等待支付
	SearchTradeProcess = "2" //交易成功
	SearchTradeClosed  = "3" //交易关闭
	SearchTradeError   = "4" //交易失败
	SearchTradeRevoked = "5" //交撤销
	SearchTradeNotPay  = "6" //未支付
	SearchTradeRefund  = "7" //转入退款
)

const (
	SearchTradeNetErrCode                   = 10504
	SearchTradeNetErrMessage                = "查询交易流水,网络错误"
	SearchTradeResponseDataFormatErrCode    = 10505
	SearchTradeResponseDataFormatErrMessage = "查询交易流水,返回数据格式错误"
	SearchTradeResponseDataSignErrCode      = 10508
	SearchTradeResponseDataSignErrMessage   = "查询交易流水,返回数据签名校验错误"
	SearchTradeRateFormatErrCode            = 10520
	SearchTradeRateFormatErrMessage         = "查询交易流水,汇率查询结果格式错误"
)

type Trade struct{}

type TradeArg struct {
	OrderNum      string  `json:"orderNum"`
	PaymentSchema string  `json:"paymentSchema"` //渠道代码
	MerId         string  `json:"merID"`         //商户号ID ，由 GoAllPay 分配
	AcqId         string  `json:"acqID"`         //收单行ID "99020344"
	Md5Key        string  `json:"md5Key"`        //安全code
	TradeGateWay  string  `json:"tradeGgateWay"` //查询交易网关地址
	SapiGateWay   string  `json:"sapiGateWay"`   //查询汇率网关地址
	Currency      string  `json:"currency"`      //币种
	TotalFee      float64 `json:"total_fee"`     //金额
}

type TradeResult struct {
	RespCode string `json:"RespCode"`
	RespMsg  string `json:"RespMsg"`
	TransId  string `json:"transID"`
}

type TradeRsp struct {
	Status  string  `json:"status"`   //交易状态
	OrderId string  `json:"order_id"` //订单号
	TradeNo string  `json:"trade_no"` //支付机构交易流水号
	PaidAt  string  `json:"paid_at"`  //支付gmt时间
	RmbFee  float64 `json:"rmb_fee"`  //人民币金额
	Rate    float64 `json:"rate"`     //汇率
	Res     string  `json:"res"`
	Rsp     string  `json:"rsp"`
}

func (trade *Trade) Search(arg TradeArg) (tradeRsp TradeRsp, errCode int, err error) {
	transTime := time.Now().Format(TimeLayout)
	paramMap := map[string]string{
		"version":       Version,
		"charSet":       CharsetUTF8,
		"transType":     SearchTradeTransType,
		"orderNum":      arg.OrderNum,
		"merID":         arg.MerId,
		"acqID":         arg.AcqId,
		"paymentSchema": arg.PaymentSchema,
		"transTime":     transTime,
		"signType":      MD5SignType,
	}
	paramMap["signature"] = trade.getSign(paramMap, arg.Md5Key)
	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}
	//fmt.Println("values.Encode()", arg.TradeGateWay+"?"+values.Encode())
	returnBytes, err := curl.GetJSONReturnByte(trade.getGateWay(arg.TradeGateWay) + "?" + values.Encode())
	if err != nil {
		logrus.Errorf("org:allpay,"+SearchTradeNetErrMessage+",order id %v,errCode:%v,err:%v", arg.OrderNum, SearchTradeNetErrCode, err.Error())
		return tradeRsp, SearchTradeNetErrCode, errors.New(SearchTradeNetErrMessage)
	}
	rspMap := make(map[string]string)
	err = json.Unmarshal(returnBytes, &rspMap)
	if err != nil {
		logrus.Errorf("org:allpay,"+SearchTradeResponseDataFormatErrMessage+",order id %v,errCode:%v,err:%v", arg.OrderNum, SearchTradeResponseDataFormatErrCode, err.Error())
		return tradeRsp, SearchTradeResponseDataFormatErrCode, errors.New(SearchTradeResponseDataFormatErrMessage)
	}
	fmt.Printf("%+v\n", rspMap)
	tradeRsp.OrderId = arg.OrderNum
	//校验签名
	sign := rspMap["signature"]
	delete(rspMap, "signature")
	if !trade.checkSign(rspMap, arg.Md5Key, sign) {
		logrus.Errorf("org:allpay,"+SearchTradeResponseDataFormatErrMessage+",order id %v,errCode:%v", arg.OrderNum, SearchTradeResponseDataFormatErrCode)
		return tradeRsp, SearchTradeResponseDataSignErrCode, errors.New(SearchTradeResponseDataSignErrMessage)
	}

	tradeRsp.Status = SearchTradeWait
	//交易状态
	if rspMap["RespCode"] == "00" {
		tradeRsp.Status = SearchTradeProcess
		tradeRsp.TradeNo = rspMap["transID"]
	}

	//allpay不返回支付成功时间，取当前异步通知时的服务器时间
	tradeRsp.PaidAt = time.Now().UTC().Format(DateTimeFormatLayout)

	//汇率
	rateArg := RateArg{
		MerId:                  arg.MerId,
		OriginalCurrencyCode:   arg.Currency,
		ConversionCurrencyCode: CNY,
		Md5Key:                 arg.Md5Key,
		PaymentSchema:          rspMap["paymentSchema"],
		GateWay:                arg.SapiGateWay,
	}
	tradeRsp.Rate, errCode, err = new(Rate).Search(rateArg)
	if err != nil {
		logrus.Errorf("org:allpay,"+SearchTradeRateFormatErrMessage+",errCode:%v,err:%v", SearchTradeRateFormatErrCode, err.Error())
		return tradeRsp, SearchTradeRateFormatErrCode, errors.New(SearchTradeRateFormatErrMessage)
	}

	//人民币金额
	tradeRsp.RmbFee, err = strconv.ParseFloat(fmt.Sprintf("%.2f", arg.TotalFee*tradeRsp.Rate), 64)
	if err != nil {
		logrus.Errorf("org:allpay,"+SearchTradeRateFormatErrMessage+",errCode:%v,err:%v", SearchTradeRateFormatErrCode, err.Error())
		return tradeRsp, SearchTradeRateFormatErrCode, errors.New(SearchTradeRateFormatErrMessage)
	}

	return tradeRsp, 0, nil
}

func (trade *Trade) getSign(paramMap map[string]string, signKey string) string {
	sortString := util.GetSortString(paramMap)
	return util.Md5(sortString + signKey)
}

func (trade *Trade) checkSign(rspMap map[string]string, md5Key, sign string) bool {
	sortString := util.GetSortString(rspMap)
	calculateSign := util.Md5(sortString + md5Key)
	return calculateSign == sign
}

func (trade *Trade) getGateWay(gateWay string) string {
	return gateWay + TradeRoute
}
