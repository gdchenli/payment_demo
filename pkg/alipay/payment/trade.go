package payment

import (
	"encoding/xml"
	"errors"
	"net/url"
	"payment_demo/pkg/alipay/util"
	"payment_demo/tools/curl"
	"strconv"
)

const (
	SearchServiceType = "single_trade_query"
)

const (
	SearchTradeWait    = "1" //等待交易
	SearchTradeProcess = "2" //交易成功
	SearchTradeClosed  = "3" //交易关闭
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
	Merchant   string `json:"merchant"`
	OutTradeNo string `json:"out_trade_no"`
	Md5Key     string `json:"md5_key"`
	TradeWay   string `json:"trade_way"`
}

type SearchResult struct {
	XMLName   xml.Name      `xml:"alipay" json:"alipay"`         //指定最外层的标签为alipay
	IsSuccess string        `xml:"is_success" json:"is_success"` //读取is_success
	Response  TradeResponse `xml:"response" json:"response"`     //读取response
	Sign      string        `xml:"sign" json:"sign"`             //读取sign
	SignType  string        `xml:"sign_type" json:"sign_type"`   //读取sign_type
}

type TradeResponse struct {
	TradeXml TradeXml `xml:"trade" json:"trade"` //读取trade
}

type TradeXml struct {
	Body                string `xml:"body" json:"body"`                                     //读取body
	BuyerEmail          string `xml:"buyer_email" json:"buyer_email"`                       //读取buyer_email
	BuyerId             string `xml:"buyer_id" json:"buyer_id"`                             //读取buyer_id
	Discount            string `xml:"discount" json:"discount" xml:"discount"`              //读取discount
	FlagTradeLocked     int    `xml:"flag_trade_locked" json:"flag_trade_locked"`           //读取flag_trade_locked
	GmtCreate           string `xml:"gmt_create" json:"gmt_create"`                         //读取gmt_create
	GmtLastModifiedTime string `xml:"gmt_last_modified_time" json:"gmt_last_modified_time"` //读取gmt_last_modified_time
	GmtPayment          string `xml:"gmt_payment" json:"gmt_payment"`                       //读取gmt_payment
	IsTotalFeeAdjust    string `xml:"is_total_fee_adjust" json:"is_total_fee_adjust"`       //读取is_total_fee_adjust
	OperatorRole        string `xml:"operator_role" json:"operator_role"`                   //读取operator_role
	OutTradeNo          string `xml:"out_trade_no" json:"out_trade_no"`                     //读取out_trade_no
	PaymentType         string `xml:"payment_type" json:"payment_type"`                     //读取payment_type
	Price               string `xml:"price" json:"price"`                                   //读取price
	Quantity            int    `xml:"quantity" json:"quantity"`                             //读取quantity
	SellerEmail         string `xml:"seller_email" json:"seller_email"`                     //读取seller_email
	SellerId            string `xml:"seller_id" json:"seller_id"`                           //读取seller_id
	Subject             string `xml:"subject" json:"subject"`                               //读取subject
	ToBuyerFee          string `xml:"to_buyer_fee" json:"to_buyer_fee"`                     //读取to_buyer_fee
	ToSellerFee         string `xml:"to_seller_fee" json:"to_seller_fee"`                   //读取to_seller_fee
	TotalFee            string `xml:"total_fee" json:"total_fee"`                           //读取total_fee
	TradeNo             string `xml:"trade_no" json:"trade_no"`                             //读取trade_no
	TradeStatus         string `xml:"trade_status" json:"trade_status"`                     //读取trade_status
	UseCoupon           string `xml:"use_coupon" json:"use_coupon"`                         //读取use_coupon
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
		"service":        SearchServiceType, //交易查询服务
		"partner":        arg.Merchant,      //商户ID
		"_input_charset": CharsetUTF8,       //编码
		"out_trade_no":   arg.OutTradeNo,    //订单编号
	}
	payString := util.GetSortString(paramMap)
	paramMap["sign"] = util.Md5(payString + arg.Md5Key)
	paramMap["sign_type"] = SignTypeMD5
	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}
	tradeRsp.Res = values.Encode()
	returnBytes, err := curl.GetJSONReturnByte(arg.TradeWay + "?" + values.Encode())
	if err != nil {
		return tradeRsp, SearchTradeNetErrCode, errors.New(SearchTradeNetErrMessage)
	}
	tradeRsp.Rsp = string(returnBytes)

	var searchResult SearchResult
	if err = xml.Unmarshal(returnBytes, &searchResult); err != nil {
		return tradeRsp, SearchTradeResponseDataFormatErrCode, errors.New(SearchTradeResponseDataFormatErrMessage)
	}
	if !trade.checkSign(searchResult, arg.Md5Key) {
		return tradeRsp, SearchTradeResponseDataSignErrCode, errors.New(SearchTradeResponseDataSignErrMessage)
	}

	//交易状态
	switch searchResult.Response.TradeXml.TradeStatus {
	case TradeSuccess:
		tradeRsp.Status = SearchTradeProcess
	case TradeFinished:
		tradeRsp.Status = SearchTradeProcess
	case TradeWaitBuyPay:
		tradeRsp.Status = SearchTradeWait
	case TradeClosed:
		tradeRsp.Status = SearchTradeClosed
	}
	tradeRsp.TradeNo = searchResult.Response.TradeXml.TradeNo

	return tradeRsp, 0, nil
}

//验证查询交易结果
func (trade *Trade) checkSign(searchResult SearchResult, md5Key string) bool {
	payMap := map[string]string{
		"body":                   searchResult.Response.TradeXml.Body,
		"buyer_email":            searchResult.Response.TradeXml.BuyerEmail,
		"buyer_id":               searchResult.Response.TradeXml.BuyerId,
		"discount":               searchResult.Response.TradeXml.Discount,
		"flag_trade_locked":      strconv.Itoa(searchResult.Response.TradeXml.FlagTradeLocked),
		"gmt_create":             searchResult.Response.TradeXml.GmtCreate,
		"gmt_last_modified_time": searchResult.Response.TradeXml.GmtLastModifiedTime,
		"gmt_payment":            searchResult.Response.TradeXml.GmtPayment,
		"is_total_fee_adjust":    searchResult.Response.TradeXml.IsTotalFeeAdjust,
		"operator_role":          searchResult.Response.TradeXml.OperatorRole,
		"out_trade_no":           searchResult.Response.TradeXml.OutTradeNo,
		"payment_type":           searchResult.Response.TradeXml.PaymentType,
		"price":                  searchResult.Response.TradeXml.Price,
		"quantity":               strconv.Itoa(searchResult.Response.TradeXml.Quantity),
		"seller_email":           searchResult.Response.TradeXml.SellerEmail,
		"seller_id":              searchResult.Response.TradeXml.SellerId,
		"subject":                searchResult.Response.TradeXml.Subject,
		"to_buyer_fee":           searchResult.Response.TradeXml.ToBuyerFee,
		"to_seller_fee":          searchResult.Response.TradeXml.ToSellerFee,
		"total_fee":              searchResult.Response.TradeXml.TotalFee,
		"trade_no":               searchResult.Response.TradeXml.TradeNo,
		"trade_status":           searchResult.Response.TradeXml.TradeStatus,
		"use_coupon":             searchResult.Response.TradeXml.UseCoupon,
	}

	payString := util.GetSortString(payMap)
	compareSignature := util.Md5(payString + md5Key)

	return compareSignature == searchResult.Sign
}
