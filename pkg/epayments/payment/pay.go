package payment

import (
	"encoding/json"
	"errors"
	"net/url"
	"payment_demo/pkg/epayments/util"
	"payment_demo/tools/curl"
	"strconv"
)

const (
	createSmartPay   = "create_smart_pay"
	AggregateCodePay = "aggregate_code_pay"
	AlipayOrgCode    = "alipay"         //支付宝
	AlipayMethodCode = "alipay_payment" //支付宝支付
	WechatMethodCode = "wechat"         //微信支付
)

const (
	SuccessCode = "0"
)

const (
	PayNetErrCode    = 10110
	PayNetErrMessage = "发起支付,网络错误"
)

type Payment struct{}

type PayArg struct {
	MerchantId      string  `json:"merchant_id"`    //商户ID
	NotifyUrl       string  `json:"notify_url"`     //支付结果异步通知到该地址
	ReturnUrl       string  `json:"return_url"`     //支付结果异步通知到该地址
	ValidMins       string  `json:"valid_mins"`     //创建交易有效期，单位为分钟，超过时间，订单失效，不传入，默认1小时。
	IncrementId     string  `json:"increment_id"`   //订单号
	GrandTotal      float64 `json:"grandtotal"`     //订单金额
	Currency        string  `json:"currency"`       //订单币种
	GateWay         string  `json:"gate_way"`       //网关地址
	SecretKey       string  `json:"secret_key"`     //密钥
	TransCurrency   string  `json:"trans_currency"` //结算币种
	PaymentChannels string  `json:"payment_channels"`
}

func (payment *Payment) CreateForm(arg PayArg) (form string, errCode int, err error) {
	//支付金额处理
	var grandTotal string
	if arg.Currency == KRW || arg.Currency == JPY {
		grandTotal = strconv.FormatFloat(arg.GrandTotal, 'f', 0, 64)
	} else {
		grandTotal = strconv.FormatFloat(arg.GrandTotal, 'f', 2, 64)
	}

	paramMap := map[string]string{
		"service":          createSmartPay,
		"merchant_id":      arg.MerchantId,
		"notify_url":       arg.NotifyUrl,
		"return_url":       arg.ReturnUrl,
		"subject":          arg.IncrementId,
		"grandtotal":       grandTotal,
		"increment_id":     arg.IncrementId,
		"currency":         arg.TransCurrency,
		"payment_channels": arg.PaymentChannels,
		"describe":         arg.IncrementId,
		"nonce_str":        util.GetRandomString(20),
	}

	//超时时间
	if arg.ValidMins != "" {
		paramMap["valid_mins"] = arg.ValidMins
	}

	sortString := util.GetSortString(paramMap)
	paramMap["signature"] = util.Md5(sortString)
	paramMap["sign_type"] = SignTypeMD5

	//生成form表单
	form = payment.buildForm(paramMap, arg.GateWay)

	return form, 0, nil
}

func (payment *Payment) buildForm(paramMap map[string]string, gateWay string) (form string) {
	payUrl := "<form action='" + gateWay + "' method='post' id='pay_form'>"
	for k, v := range paramMap {
		payUrl += "<input value='" + v + "' name='" + k + "' type='hidden'/>"
	}
	payUrl += "</form>"
	payUrl += "<script>var form = document.getElementById('pay_form');form.submit()</script>"
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

func (payment *Payment) CreateQrCode(arg PayArg) (qrCodeUrl string, errCode int, err error) {
	//支付金额处理
	var grandTotal string
	if arg.Currency == KRW || arg.Currency == JPY {
		grandTotal = strconv.FormatFloat(arg.GrandTotal, 'f', 0, 64)
	} else {
		grandTotal = strconv.FormatFloat(arg.GrandTotal, 'f', 2, 64)
	}

	paramMap := map[string]string{
		"service":          AggregateCodePay,
		"merchant_id":      arg.MerchantId,
		"notify_url":       arg.NotifyUrl,
		"subject":          arg.IncrementId,
		"grandtotal":       grandTotal,
		"increment_id":     arg.IncrementId,
		"currency":         arg.TransCurrency,
		"payment_channels": arg.PaymentChannels,
		"describe":         arg.IncrementId,
		"nonce_str":        util.GetRandomString(20),
	}

	//超时时间
	if arg.ValidMins != "" {
		paramMap["valid_mins"] = arg.ValidMins
	}

	sortString := util.GetSortString(paramMap)
	paramMap["signature"] = util.Md5(sortString)
	paramMap["sign_type"] = SignTypeMD5

	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}
	returnBytes, err := curl.GetJSONReturnByte(arg.GateWay + "?" + values.Encode())
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

	return qrCodeResult.CodeUrl, 0, nil
}
