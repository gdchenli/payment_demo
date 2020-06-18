package epayments

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"payment_demo/api/request"
	"payment_demo/api/response"
	"payment_demo/pkg/curl"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	SearchTradeWait    = "0" //等待支付
	SearchTradeProcess = "2" //交易成功
	SearchTradeClosed  = "3" //交易关闭
	SearchTradeError   = "4" //交易失败
	SearchTradeRevoked = "5" //交撤销
	SearchTradeNotPay  = "6" //未支付
	SearchTradeRefund  = "7" //转入退款
)
const (
	SearchServiceType = "create_trade_query"
)

const (
	SearchTradeNetErrCode                   = 10504
	SearchTradeNetErrMessage                = "查询交易流水,网络错误"
	SearchTradeResponseDataFormatErrCode    = 10505
	SearchTradeResponseDataFormatErrMessage = "查询交易流水,返回数据格式错误"
	SearchTradeResponseDataSignErrCode      = 10508
	SearchTradeResponseDataSignErrMessage   = "查询交易流水,返回数据签名校验错误"
)

func (epayments *Epayments) SearchTrade(paramMap map[string]string, req request.SearchTradeReq) (tradeRsp response.SearchTradeRsp, errCode int, err error) {
	md5Key := paramMap["md5_key"]
	delete(paramMap, "md5_key")
	gateWay := paramMap["gate_way"]
	delete(paramMap, "gate_way")

	paramMap["service"] = SearchServiceType
	paramMap["increment_id"] = req.OrderId
	paramMap["nonce_str"] = GetRandomString(20)
	sortString := GetSortString(paramMap)
	paramMap["signature"] = Md5(sortString + md5Key)
	paramMap["sign_type"] = SignTypeMD5
	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}

	fmt.Println(gateWay + "?" + values.Encode())
	returnBytes, err := curl.GetJSONReturnByte(gateWay + "?" + values.Encode())
	if err != nil {
		return tradeRsp, SearchTradeNetErrCode, errors.New(SearchTradeNetErrMessage)
	}

	rspMap := make(map[string]interface{})
	err = json.Unmarshal(returnBytes, &rspMap)
	if err != nil {
		logrus.Errorf("org:epayments,"+SearchTradeResponseDataFormatErrMessage+",orderId:%v,query:%v,errCode:%v,err:%v", req.OrderId, sortString, SearchTradeResponseDataFormatErrCode, err.Error())
		return tradeRsp, SearchTradeResponseDataFormatErrCode, errors.New(SearchTradeResponseDataFormatErrMessage)
	}
	tradeRsp.OrderId = req.OrderId

	//校验签名
	sign := rspMap["signature"].(string)
	tradeRspMap := make(map[string]string)
	for k, v := range rspMap {
		if k == "signature" || k == "sign_type" {
			continue
		}
		tradeRspMap[k] = fmt.Sprintf("%v", v)
	}
	if !checkSearchTradeSign(tradeRspMap, md5Key, sign) {
		logrus.Errorf("org:epayments,"+SearchTradeResponseDataSignErrMessage+",orderId:%v,query:%v,errCode:%v", req.OrderId, sortString, SearchTradeResponseDataSignErrCode)
		return tradeRsp, SearchTradeResponseDataSignErrCode, errors.New(SearchTradeResponseDataSignErrMessage)
	}

	tradeRsp.TradeNo = tradeRspMap["trade_no"]
	switch rspMap["trade_status"] {
	case TradeFinished:
		tradeRsp.Status = SearchTradeProcess
	case TradeSuccess:
		tradeRsp.Status = SearchTradeProcess
	case TradeWaitBuyPay:
		tradeRsp.Status = SearchTradeWait
	case TradeError:
		tradeRsp.Status = SearchTradeError
	case TradeClosed:
		tradeRsp.Status = SearchTradeClosed
	case TradeNotPay:
		tradeRsp.Status = SearchTradeNotPay
	case TradeRevoked:
		tradeRsp.Status = SearchTradeRevoked
	case TradeRefund:
		tradeRsp.Status = SearchTradeRefund
	}
	if tradeRsp.Status != TradeSuccess {
		return tradeRsp, 0, nil
	}

	//支付时间
	parseTime, err := time.Parse(DateTimeFormatLayout, tradeRspMap["gmt_payment"])
	if err != nil {
		logrus.Errorf("org:epayments,"+SearchTradeResponseDataFormatErrMessage+",orderId:%v,query:%v,errCode:%v,err:%v", req.OrderId, sortString, SearchTradeResponseDataFormatErrCode, err.Error())
		return tradeRsp, SearchTradeResponseDataFormatErrCode, errors.New(SearchTradeResponseDataFormatErrMessage)
	}
	tradeRsp.PaidAt = parseTime.UTC().Format(DateTimeFormatLayout)

	//汇率
	tradeRsp.Rate, err = strconv.ParseFloat(tradeRspMap["rate"], 64)
	if err != nil {
		logrus.Errorf("org:epayments,"+SearchTradeResponseDataFormatErrMessage+",orderId:%v,query:%v,errCode:%v,err:%v", req.OrderId, sortString, SearchTradeResponseDataFormatErrCode, err.Error())
		return tradeRsp, SearchTradeResponseDataFormatErrCode, errors.New(SearchTradeResponseDataFormatErrMessage)
	}

	//人民币金额
	grandTotal, err := strconv.ParseFloat(tradeRspMap["grandtotal"], 64)
	if err != nil {
		logrus.Errorf("org:epayments,"+SearchTradeResponseDataFormatErrMessage+",orderId:%v,query:%v,errCode:%v,err:%v", req.OrderId, sortString, SearchTradeResponseDataFormatErrCode, err.Error())
		return tradeRsp, SearchTradeResponseDataFormatErrCode, errors.New(SearchTradeResponseDataFormatErrMessage)
	}
	rmbFee := grandTotal * tradeRsp.Rate
	rmbFee, err = strconv.ParseFloat(fmt.Sprintf("%.2f", rmbFee), 64)
	if err != nil {
		logrus.Errorf("org:epayments,"+SearchTradeResponseDataFormatErrMessage+",orderId:%v,query:%v,errCode:%v,err:%v", req.OrderId, sortString, SearchTradeResponseDataFormatErrCode, err.Error())
		return tradeRsp, SearchTradeResponseDataFormatErrCode, errors.New(SearchTradeResponseDataFormatErrMessage)
	}
	tradeRsp.RmbFee = rmbFee

	return tradeRsp, 0, nil
}

func checkSearchTradeSign(rspMap map[string]string, md5Key, sign string) bool {
	sortString := GetSortString(rspMap)
	calculateSign := Md5(sortString + md5Key)
	return calculateSign == sign
}

func (epayments *Epayments) GetSearchTradeConfigCode() []string {
	return []string{
		"merchant_id",
		"md5_key", "gate_way",
	}
}
