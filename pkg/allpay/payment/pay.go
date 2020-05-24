package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"payment_demo/pkg/allpay/util"
	"time"
)

const (
	PayTransType = "PURC"
)

const (
	PayGoodsInfoFormatErrCode    = 10101
	PayGoodsInfoFormatErrMessage = "发起支付，商品数据转换失败"
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

func (payment *Payment) CreateForm(arg PayArg) (form string, errCode int, err error) {
	orderAmount := fmt.Sprintf("%.2f", arg.OrderAmount)
	transTime := time.Now().Format(TimeLayout)
	detailInfoBytes, err := json.Marshal(arg.DetailInfo)
	if err != nil {
		return form, PayGoodsInfoFormatErrCode, errors.New(PayGoodsInfoFormatErrMessage)
	}
	detailInfo := util.BASE64EncodeStr(detailInfoBytes)

	paramMap := map[string]string{
		"version":       Version,
		"charSet":       CharSet,
		"transType":     PayTransType,
		"orderNum":      arg.OrderNum,
		"orderAmount":   orderAmount,
		"orderCurrency": arg.OrderCurrency,
		"frontURL":      arg.FrontUrl,
		"backURL":       arg.BackUrl,
		"merID":         arg.MerId,
		"acqID":         arg.AcqId,
		"paymentSchema": arg.PaymentSchema,
		"goodsInfo":     arg.OrderNum,
		"detailInfo":    detailInfo,
		"transTime":     transTime,
		"signType":      PayMD5SignType,
		"tradeFrom":     arg.TradeFrom,
		"timeout":       arg.Timeout,
		"merReserve":    "",
	}
	paramMap["signature"] = payment.getSign(paramMap, arg.Md5Key)

	form = payment.buildForm(paramMap, arg.PayWay)

	return form, 0, nil
}

func (payment *Payment) getSign(paramMap map[string]string, signKey string) string {
	sortString := util.GetSortString(paramMap)
	fmt.Println("sortString", sortString)
	return util.Md5(sortString + signKey)
}

func (payment *Payment) buildForm(paramMap map[string]string, payWay string) (form string) {
	payUrl := "<form action='" + payWay + "' method='post' id='pay_form'>"
	for k, v := range paramMap {
		payUrl += "<input value='" + v + "' name='" + k + "' type='hidden'/>"
	}
	payUrl += "</form>"
	payUrl += "<script>var form = document.getElementById('pay_form');form.submit()</script>"
	return payUrl
}
