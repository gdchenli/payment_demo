package epayments

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"payment_demo/pkg/curl"
	"payment_demo/pkg/payment/common"
	"payment_demo/pkg/payment/consts"
	"strconv"

	"github.com/skip2/go-qrcode"
)

const (
	createSmartPay   = "create_smart_pay"
	AggregateCodePay = "aggregate_code_pay"
)

const (
	SuccessCode = "0"
)

const (
	PayNetErrCode    = 10110
	PayNetErrMessage = "发起支付,网络错误"
)

func (epayments *Epayments) CreatePayUrl(paramMap map[string]string, order common.OrderArg) (payUrl string, errCode int, err error) {
	if order.UserAgentType == consts.WebUserAgentType {
		marshal, _ := json.Marshal(order)
		reqParamMap := make(map[string]interface{})
		json.Unmarshal(marshal, &reqParamMap)
		values := url.Values{}
		for k, v := range reqParamMap {
			values.Add(k, fmt.Sprintf("%v", v))
		}
		return "/payment/qrcode?" + values.Encode(), 0, nil
	}

	gateWay := paramMap["gate_way"]
	delete(paramMap, "gate_way")

	md5Key := paramMap["md5_key"]
	delete(paramMap, "md5_key")

	//支付金额处理
	var grandTotal string
	if order.Currency == KRW || order.Currency == JPY {
		grandTotal = strconv.FormatFloat(order.TotalFee, 'f', 0, 64)
	} else {
		grandTotal = strconv.FormatFloat(order.TotalFee, 'f', 2, 64)
	}

	paramMap["service"] = createSmartPay
	paramMap["subject"] = order.OrderId
	paramMap["grandtotal"] = grandTotal
	paramMap["increment_id"] = order.OrderId
	paramMap["payment_channels"] = getPayPaymentChannels(order.MethodCode)
	paramMap["describe"] = order.OrderId
	paramMap["nonce_str"] = GetRandomString(20)
	sortString := GetSortString(paramMap)
	paramMap["signature"] = Md5(sortString + md5Key)
	paramMap["sign_type"] = SignTypeMD5

	payUrl = buildPayUrl(paramMap, gateWay)

	return payUrl, 0, nil
}

func buildPayUrl(paramMap map[string]string, gateWay string) (payUrl string) {
	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}
	payUrl = gateWay + "?" + values.Encode()
	return payUrl
}

type QrCodeResult struct {
	Code        string `json:"code"`         //结果代码 0 成功
	Message     string `json:"message"`      //结果描述，如果返回错误，则为错误描述原因
	MerchantId  string `json:"merchant_id"`  //商户ID
	IncrementId string `json:"increment_id"` //商户订单号
	CodeUrl     string `json:"code_url"`     //二维码原链接，用于生成二维码
	NonceStr    string `json:"nonce_str"`    //随机字符串
	Signature   string `json:"signature"`    //签名
	SignType    string `json:"sign_type"`    //签名类型
}

func (epayments *Epayments) CreateQrCode(paramMap map[string]string, order common.OrderArg) (qrCodeUrl string, errCode int, err error) {
	gateWay := paramMap["gate_way"]
	delete(paramMap, "gate_way")

	md5Key := paramMap["md5_key"]
	delete(paramMap, "md5_key")

	delete(paramMap, "return_url")

	//支付金额处理
	var grandTotal string
	if order.Currency == KRW || order.Currency == JPY {
		grandTotal = strconv.FormatFloat(order.TotalFee, 'f', 0, 64)
	} else {
		grandTotal = strconv.FormatFloat(order.TotalFee, 'f', 2, 64)
	}

	paramMap["service"] = AggregateCodePay
	paramMap["subject"] = order.OrderId
	paramMap["grandtotal"] = grandTotal
	paramMap["increment_id"] = order.OrderId
	paramMap["payment_channels"] = getPayPaymentChannels(order.MethodCode)
	paramMap["describe"] = order.OrderId
	paramMap["nonce_str"] = GetRandomString(20)
	sortString := GetSortString(paramMap)
	paramMap["signature"] = Md5(sortString + md5Key)
	paramMap["sign_type"] = SignTypeMD5

	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}
	returnBytes, err := curl.GetJSONReturnByte(gateWay + "?" + values.Encode())
	if err != nil {
		return qrCodeUrl, PayNetErrCode, errors.New(PayNetErrMessage)
	}

	var qrCodeResult QrCodeResult
	if err = json.Unmarshal(returnBytes, &qrCodeResult); err != nil {
		return qrCodeUrl, PayNetErrCode, errors.New(PayNetErrMessage)
	}

	if qrCodeResult.Code != SuccessCode {
		return qrCodeUrl, PayNetErrCode, errors.New(PayNetErrMessage)
	}
	if qrCodeResult.CodeUrl == "" {
		return qrCodeUrl, PayNetErrCode, errors.New(PayNetErrMessage)
	}
	qrCodeBytes, err := qrcode.Encode(qrCodeResult.CodeUrl, qrcode.Medium, 256)
	if qrCodeResult.CodeUrl == "" {
		return qrCodeUrl, PayNetErrCode, errors.New(PayNetErrMessage)
	}
	imgStr := "data:image/png;base64," + BASE64EncodeStr(qrCodeBytes)
	return imgStr, 0, nil
}

//获取支付通道
func getPayPaymentChannels(methodCode string) (paymentChannels string) {
	if methodCode == consts.WechatMethod {
		paymentChannels = ChannelWechat
	} else if methodCode == consts.AlipayMethod {
		paymentChannels = ChannelAlipay
	}
	return paymentChannels
}

func (epayments *Epayments) GetPayConfigCode() []string {
	return []string{
		"merchant_id", "notify_url", "return_url", "currency", "valid_mins",
		"md5_key", "gate_way",
	}
}
