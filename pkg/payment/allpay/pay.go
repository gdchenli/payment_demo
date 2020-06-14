package allpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"payment_demo/api/validate"
	"payment_demo/pkg/curl"
	"payment_demo/pkg/payment/consts"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	PayTransType = "PURC"
	PayRoute     = "/api/unifiedorder"
)

const (
	PayGoodsInfoFormatErrCode    = 10101
	PayGoodsInfoFormatErrMessage = "发起支付，商品数据转换失败"
	PayNetErrCode                = 10110
	PayNetErrMessage             = "发起支付,网络错误"
)

type Payment struct{}

type PayArg struct {
	OrderNum      string       `json:"orderNum"`      //订单号：商户自行定义，需保证同一商户号下订单号不能重复
	OrderAmount   float64      `json:"orderAmount"`   //订单金额：如 100 元，表示为 100 或 100.00
	OrderCurrency string       `json:"orderCurrency"` //订单币种：ISO标准。如：人民币填写“CNY”，美元填写"USD"
	FrontUrl      string       `json:"frontURL"`      //支付完成后跳转到该地址
	BackUrl       string       `json:"backURL"`       //支付结果异步通知到该地址
	MerId         string       `json:"merID"`         //商户号ID ，由 GoAllPay 分配
	AcqId         string       `json:"acqID"`         //收单行ID "99020344"
	PaymentSchema string       `json:"paymentSchema"` //渠道代码
	GoodsInfo     string       `json:"goodsInfo"`     //商品信息
	DetailInfo    []DetailInfo `json:"detailInfo"`    //商品明细，格式：[{"goods_name":"iPhone X","quantity":"2"},{"goods_name":"iPhone 8","quantity":"4"}]，需对该字段进行base-64编码后签名上送
	Md5Key        string       `json:"md5Key"`        //安全code
	PayWay        string       `json:"payWay"`        //网关地址
	TradeFrom     string       `json:"tradeFrom"`     //渠道
	Timeout       string       `json:"timeout"`
}

type DetailInfo struct {
	GoodsName string `json:"goods_name"`
	Quantity  int    `json:"quantity"`
}

func (payment *Payment) CreatePayUrl(configParamMap map[string]string, order validate.Order) (payUrl string, errCode int, err error) {
	gateWay := payment.getGateWay(configParamMap["gate_way"])
	delete(configParamMap, "gata_way")

	paramMap, errCode, err := payment.getPayParamMap(configParamMap, order)
	if err != nil {
		return payUrl, errCode, err
	}

	payUrl = payment.buildPayUrl(paramMap, gateWay)

	return payUrl, 0, nil
}

func (payment *Payment) CreateAmpPayStr(configParamMap map[string]string, order validate.Order) (payStr string, errCode int, err error) {
	gateWay := payment.getGateWay(configParamMap["gate_way"])
	delete(configParamMap, "gata_way")

	paramMap, errCode, err := payment.getPayParamMap(configParamMap, order)
	if err != nil {
		return payStr, errCode, err
	}
	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}
	var ampProgramRsp AmpProgramRsp

	err = curl.PostJSON(gateWay, &ampProgramRsp, strings.NewReader(values.Encode()))
	if err != nil {
		logrus.Errorf("org:allpay,"+PayNetErrMessage+",order id %v,errCode:%v,err:%v", order.OrderId, PayNetErrCode, err.Error())
		return payStr, PayNetErrCode, errors.New(PayNetErrMessage)
	}

	if ampProgramRsp.RespCode != "00" {
		return payStr, PayNetErrCode, errors.New(PayNetErrMessage)
	}

	return ampProgramRsp.SdkParams, 0, nil
}

func (payment *Payment) getPayParamMap(paramMap map[string]string, order validate.Order) (map[string]string, int, error) {
	orderAmount := fmt.Sprintf("%.2f", order.TotalFee)
	transTime := time.Now().Format(TimeLayout)
	detailInfoBytes, err := json.Marshal([]DetailInfo{{
		GoodsName: SpecialReplace("test goods name" + time.Now().Format(TimeLayout)),
		Quantity:  1,
	}})
	if err != nil {
		logrus.Errorf("org:allpay,"+PayGoodsInfoFormatErrMessage+",order id %v,errCode:%v,err:%v", order.OrderId, PayGoodsInfoFormatErrCode, err.Error())
		return paramMap, PayGoodsInfoFormatErrCode, errors.New(PayGoodsInfoFormatErrMessage)
	}
	detailInfo := BASE64EncodeStr(detailInfoBytes)
	paramMap["version"] = Version
	paramMap["charSet"] = CharsetUTF8
	paramMap["transType"] = PayTransType
	paramMap["orderNum"] = order.OrderId
	paramMap["orderAmount"] = orderAmount
	paramMap["orderCurrency"] = order.Currency
	paramMap["paymentSchema"] = payment.getPaymentSchema(order.MethodCode)
	paramMap["goodsInfo"] = order.OrderId
	paramMap["detailInfo"] = detailInfo
	paramMap["transTime"] = transTime
	paramMap["signType"] = MD5SignType
	paramMap["tradeFrom"] = payment.getTradeFrom(order.MethodCode, order.UserAgentType)
	paramMap["merReserve"] = ""

	md5key := paramMap["md5_key"]
	delete(paramMap, "md5_key")
	paramMap["signature"] = payment.getSign(paramMap, md5key)

	return paramMap, 0, nil
}
func (payment *Payment) getPaymentSchema(methodCode string) string {
	switch methodCode {
	case consts.AlipayMethod:
		return ApSchema
	case consts.UnionpayMethod:
		return UpSchema
	default:
		return ""
	}
}

func (payment *Payment) getTradeFrom(methodCode string, userAgentType int) string {
	if methodCode == consts.AlipayMethod {
		switch userAgentType {
		case consts.WebUserAgentType:
			return AlipayWebTradeFrom
		case consts.MobileUserAgentType:
			return AlipayMobileTradeFrom
		case consts.AlipayMiniProgramUserAgentType:
			return AlipayMiniProgramTradeFrom
		case consts.AndroidAppUserAgentType, consts.IOSAppUserAgentType:
			return AppTradeFrom
		}
	}

	if methodCode == consts.UnionpayMethod {
		return UpTradeFrom
	}

	return ""
}

func (payment *Payment) getSign(paramMap map[string]string, signKey string) string {
	sortString := GetSortString(paramMap)
	return Md5(sortString + signKey)
}

func (payment *Payment) buildPayUrl(paramMap map[string]string, gateWay string) (payUrl string) {
	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}
	payUrl = gateWay + "?" + values.Encode()
	return payUrl
}

type AmpProgramRsp struct {
	RespCode  string `json:"RespCode"`
	ResMsg    string `json:"RespMsg"`
	SdkParams string `json:"sdk_params"`
}

func (payment *Payment) getGateWay(gateWay string) string {
	return gateWay + PayRoute
}

func (payment *Payment) GetConfigCode() []string {
	return []string{
		"merID", "frontURL", "backURL", "acqID", "timeout",
		"md5_key", "gate_way",
	}
}