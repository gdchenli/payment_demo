package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"payment_demo/pkg/epayments/util"
	"payment_demo/tools/curl"

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

type Trade struct{}

type TradeArg struct {
	Merchant    string `json:"merchant"`
	IncrementId string `json:"increment_id"`
	Md5Key      string `json:"md5_key"`
	TradeWay    string `json:"trade_way"`
}

type TradeRsp struct {
	Status  string `json:"status"`   //交易状态
	OrderId string `json:"order_id"` //订单号
	TradeNo string `json:"trade_no"` //支付机构交易流水号
	Res     string `json:"res"`
	Rsp     string `json:"rsp"`
}

func (trade *Trade) Search(arg TradeArg) (tradeRsp TradeRsp, errCode int, err error) {
	paramMap := map[string]string{
		"merchant_id":  arg.Merchant,
		"increment_id": arg.IncrementId,
		"nonce_str":    util.GetRandomString(20),
		"service":      SearchServiceType,
	}
	sortString := util.GetSortString(paramMap)
	paramMap["signature"] = util.Md5(sortString + arg.Md5Key)
	paramMap["sign_type"] = SignTypeMD5
	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}

	fmt.Println(arg.TradeWay + "?" + values.Encode())
	returnBytes, err := curl.GetJSONReturnByte(arg.TradeWay + "?" + values.Encode())
	if err != nil {
		return tradeRsp, SearchTradeNetErrCode, errors.New(SearchTradeNetErrMessage)
	}

	rspMap := make(map[string]interface{})
	err = json.Unmarshal(returnBytes, &rspMap)
	if err != nil {
		logrus.Errorf(SearchTradeResponseDataFormatErrMessage+",query:%v,errCode:%v,err:%v", sortString, SearchTradeResponseDataFormatErrCode, err.Error())
		return tradeRsp, SearchTradeResponseDataFormatErrCode, errors.New(SearchTradeResponseDataFormatErrMessage)
	}
	tradeRsp.OrderId = arg.IncrementId

	//校验签名
	sign := rspMap["signature"].(string)
	tradeRspMap := make(map[string]string)
	for k, v := range rspMap {
		if k == "signature" || k == "sign_type" {
			continue
		}
		tradeRspMap[k] = fmt.Sprintf("%v", v)
	}
	if !trade.checkSign(tradeRspMap, arg.Md5Key, sign) {
		logrus.Errorf(SearchTradeResponseDataSignErrMessage+",query:%v,errCode:%v", sortString, SearchTradeResponseDataSignErrCode)
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

	return tradeRsp, 0, nil
}

func (trade *Trade) checkSign(rspMap map[string]string, md5Key, sign string) bool {
	sortString := util.GetSortString(rspMap)
	calculateSign := util.Md5(sortString + md5Key)
	return calculateSign == sign
}
