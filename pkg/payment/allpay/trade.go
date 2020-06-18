package allpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"payment_demo/api/response"
	"payment_demo/internal/common/request"
	"payment_demo/pkg/curl"
	"payment_demo/pkg/payment/consts"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
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

func (allpay *Allpay) SearchTrade(paramMap map[string]string, req request.SearchTradeReq) (tradeRsp response.SearchTradeRsp, errCode int, err error) {
	md5Key := paramMap["md5_key"]
	delete(paramMap, "md5_key")
	gateWay := getSearchTradeGateWay(paramMap["gate_way"])
	delete(paramMap, "gate_way")
	sapiWay := paramMap["sapi_way"]
	delete(paramMap, "sapi_way")

	transTime := time.Now().Format(TimeLayout)
	paramMap["version"] = Version
	paramMap["charSet"] = CharsetUTF8
	paramMap["transType"] = SearchTradeTransType
	paramMap["orderNum"] = req.OrderId
	paramMap["paymentSchema"] = getSearchTradePaymentSchema(req.MethodCode)
	paramMap["transTime"] = transTime
	paramMap["signType"] = MD5SignType
	paramMap["signature"] = getSearchTradeSign(paramMap, md5Key)
	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}
	fmt.Println("values.Encode()", gateWay+"?"+values.Encode())
	returnBytes, err := curl.GetJSONReturnByte(gateWay + "?" + values.Encode())
	if err != nil {
		logrus.Errorf("org:allpay,"+SearchTradeNetErrMessage+",order id %v,errCode:%v,err:%v", req.OrderId, SearchTradeNetErrCode, err.Error())
		return tradeRsp, SearchTradeNetErrCode, errors.New(SearchTradeNetErrMessage)
	}
	rspMap := make(map[string]string)
	err = json.Unmarshal(returnBytes, &rspMap)
	if err != nil {
		logrus.Errorf("org:allpay,"+SearchTradeResponseDataFormatErrMessage+",order id %v,errCode:%v,err:%v", req.OrderId, SearchTradeResponseDataFormatErrCode, err.Error())
		return tradeRsp, SearchTradeResponseDataFormatErrCode, errors.New(SearchTradeResponseDataFormatErrMessage)
	}
	fmt.Printf("%+v\n", rspMap)
	tradeRsp.OrderId = req.OrderId
	//校验签名
	sign := rspMap["signature"]
	delete(rspMap, "signature")
	if !checkSearchTradeSign(rspMap, md5Key, sign) {
		logrus.Errorf("org:allpay,"+SearchTradeResponseDataFormatErrMessage+",order id %v,errCode:%v", req.OrderId, SearchTradeResponseDataFormatErrCode)
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
		MerId:                  paramMap["merID"],
		OriginalCurrencyCode:   req.Currency,
		ConversionCurrencyCode: CNY,
		Md5Key:                 md5Key,
		PaymentSchema:          rspMap["paymentSchema"],
		GateWay:                sapiWay,
	}
	tradeRsp.Rate, errCode, err = allpay.SearchRate(rateArg)
	if err != nil {
		logrus.Errorf("org:allpay,"+SearchTradeRateFormatErrMessage+",errCode:%v,err:%v", SearchTradeRateFormatErrCode, err.Error())
		return tradeRsp, SearchTradeRateFormatErrCode, errors.New(SearchTradeRateFormatErrMessage)
	}

	//人民币金额
	tradeRsp.RmbFee, err = strconv.ParseFloat(fmt.Sprintf("%.2f", req.TotalFee*tradeRsp.Rate), 64)
	if err != nil {
		logrus.Errorf("org:allpay,"+SearchTradeRateFormatErrMessage+",errCode:%v,err:%v", SearchTradeRateFormatErrCode, err.Error())
		return tradeRsp, SearchTradeRateFormatErrCode, errors.New(SearchTradeRateFormatErrMessage)
	}

	return tradeRsp, 0, nil
}

func getSearchTradeSign(paramMap map[string]string, signKey string) string {
	sortString := GetSortString(paramMap)
	return Md5(sortString + signKey)
}

func checkSearchTradeSign(rspMap map[string]string, md5Key, sign string) bool {
	sortString := GetSortString(rspMap)
	calculateSign := Md5(sortString + md5Key)
	return calculateSign == sign
}

func getSearchTradeGateWay(gateWay string) string {
	return gateWay + TradeRoute
}

func getSearchTradePaymentSchema(methodCode string) string {
	switch methodCode {
	case consts.AlipayMethod:
		return ApSchema
	case consts.UnionpayMethod:
		return UpSchema
	default:
		return ""
	}
}

func (allpay *Allpay) GetSearchTradeConfigCode() []string {
	return []string{
		"merID", "acqID", "md5_key", "gate_way", "sapi_way",
	}
}
