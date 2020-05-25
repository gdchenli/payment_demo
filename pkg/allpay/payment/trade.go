package payment

import (
	"encoding/json"
	"errors"
	"net/url"
	"payment_demo/pkg/allpay/util"
	"payment_demo/tools/curl"
	"time"
)

const (
	TradeProcess = "2" //交易成功
	TradeClosed  = "3" //交易关闭
)
const (
	SearchTradeNetErrCode                   = 10504
	SearchTradeNetErrMessage                = "查询交易流水,网络错误"
	SearchTradeResponseDataFormatErrCode    = 10505
	SearchTradeResponseDataFormatErrMessage = "查询交易流水,返回数据格式错误"
	SearchTradeResponseDataSignErrCode      = 10508
	SearchTradeResponseDataSignErrMessage   = "查询交易流水,返回数据签名校验错误"
)

type Trade struct{}

type TradeArg struct {
	OrderNum      string `json:"orderNum"`
	PaymentSchema string `json:"paymentSchema"` //渠道代码
	MerId         string `json:"merID"`         //商户号ID ，由 GoAllPay 分配
	AcqId         string `json:"acqID"`         //收单行ID "99020344"
	Md5Key        string `json:"md5Key"`        //安全code
	PayWay        string `json:"payWay"`        //网关地址
}

type TradeResult struct {
	RespCode string `json:"RespCode"`
	RespMsg  string `json:"RespMsg"`
	TransId  string `json:"transID"`
}

type TradeRsp struct {
	Status  string `json:"status"`   //交易状态
	OrderId string `json:"order_id"` //订单号
	TradeNo string `json:"trade_no"` //支付机构交易流水号
	Res     string `json:"res"`
	Rsp     string `json:"rsp"`
}

func (trade *Trade) Search(arg TradeArg) (tradeRsp TradeRsp, errCode int, err error) {
	transTime := time.Now().Format(TimeLayout)
	paramMap := map[string]string{
		"version":       Version,
		"charSet":       CharSet,
		"orderNum":      arg.OrderNum,
		"merID":         arg.MerId,
		"acqID":         arg.AcqId,
		"paymentSchema": arg.PaymentSchema,
		"transTime":     transTime,
		"signType":      PayMD5SignType,
	}
	paramMap["signature"] = trade.getSign(paramMap, arg.Md5Key)
	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}

	returnBytes, err := curl.GetJSONReturnByte(arg.PayWay + "?" + values.Encode())
	if err != nil {
		return tradeRsp, SearchTradeNetErrCode, errors.New(SearchTradeNetErrMessage)
	}
	rspMap := make(map[string]string)
	err = json.Unmarshal(returnBytes, &rspMap)
	if err != nil {
		return tradeRsp, SearchTradeResponseDataFormatErrCode, errors.New(SearchTradeResponseDataFormatErrMessage)
	}

	//校验签名
	sign := rspMap["signature"]
	delete(rspMap, "signature")
	if !trade.checkSign(rspMap, arg.Md5Key, sign) {
		return tradeRsp, SearchTradeResponseDataSignErrCode, errors.New(SearchTradeResponseDataSignErrMessage)
	}

	//交易状态
	if rspMap["RespCode"] == "00" {
		tradeRsp.Status = TradeProcess
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
