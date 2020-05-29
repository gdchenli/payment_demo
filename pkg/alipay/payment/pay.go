package payment

import (
	"encoding/json"
	"payment_demo/pkg/alipay/util"
	"strconv"
)

const (
	pcServiceType          = "create_forex_trade"      //PC端支付类型
	wapServiceType         = "create_forex_trade_wap"  //移动端支付类型
	newPay                 = "2"                       //新接口
	newOverseasSeller      = "NEW_OVERSEAS_SELLER"     //海外销售产品代码
	newWapOverseasSeller   = "NEW_WAP_OVERSEAS_SELLER" //海外移动端销售产品代码
	businessTypeSalesGoods = 4
	amp                    = "amp"
)

type Payment struct{}

type PayArg struct {
	Service          string  `json:"service"`           //service,PC端:create_forex_trade,手机端:create_forex_trade_wap
	Partner          string  `json:"partner"`           //PartnerId
	NotifyUrl        string  `json:"notify_url"`        //支付结果异步通知到该地址
	ReturnUrl        string  `json:"return_url"`        //支付结果异步通知到该地址
	Subject          string  `json:"subject"`           //主题
	Body             string  `json:"body"`              //主体
	OutTradeNo       string  `json:"out_trade_no"`      //订单号
	TotalFee         float64 `json:"total_fee"`         //订单金额(CNY)
	Currency         string  `json:"currency"`          //订单币种
	Supplier         string  `json:"supplier"`          //供应商
	TimeoutRule      string  `json:"timeout_rule"`      //超时时间,如:12h,10m
	ReferUrl         string  `json:"refer_url"`         //商家url(站点url)
	ProductCode      string  `json:"product_code"`      //NEW_OVERSEAS_SELLER，海外卖家
	GateWay          string  `json:"gate_way"`          //网关地址
	SecretKey        string  `json:"secret_key"`        //密钥
	TradeInformation string  `json:"trade_information"` //交易明细json串
	Items            []Item  `json:"items"`             //订单明细
	TransCurrency    string  `json:"trade_information"` //结算币种
	OrderSource      string  `json:"order_source"`      //订单客户端 PC:web  手机端：web_mobile
	PayWay           string  `json:"pay_way"`           //版本 1旧版本 2新版本
}

type Item struct {
	Id    string  `form:"id" json:"id"`
	Name  string  `form:"name" json:"name" `
	Price float64 `form:"price" json:"price"`
	Qty   int     `form:"qty_ordered" json:"qty_ordered"`
}
type TradeInformation struct {
	BusinessType  int    `json:"business_type"`
	GoodsInfo     string `json:"goods_info"`
	TotalQuantity int    `json:"total_quantity"`
}

func (payment *Payment) CreateForm(arg PayArg) (form string, errCode int, err error) {
	paramMap := map[string]string{
		"service":           payment.getServiceType(arg.OrderSource),
		"partner":           arg.Partner,
		"return_url":        arg.ReturnUrl,
		"notify_url":        arg.NotifyUrl,
		"_input_charset":    CharsetUTF8,
		"subject":           arg.Subject,
		"body":              arg.Body,
		"out_trade_no":      arg.OutTradeNo,
		"total_fee":         payment.getTotalFee(arg),
		"currency":          arg.TransCurrency,
		"supplier":          arg.Supplier,
		"refer_url":         arg.ReferUrl,
		"trade_information": payment.getTradeInformationJson(arg.Items),
	}
	//超时时间
	if arg.TimeoutRule != "" {
		paramMap["timeout_rule"] = arg.TimeoutRule
	}
	//判断新旧版本
	if arg.PayWay == newPay {
		paramMap["product_code"] = payment.getProductCode(arg.OrderSource)
	}
	//签名
	payString := util.GetPayString(paramMap)
	paramMap["sign"] = util.Md5(payString + arg.SecretKey)

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

//获取服务类型
func (payment *Payment) getServiceType(orderSource string) (serviceType string) {
	if orderSource == Web {
		serviceType = pcServiceType
	} else if orderSource == MobileWeb {
		serviceType = wapServiceType
	} else if orderSource == amp {
		serviceType = wapServiceType
	}
	return serviceType
}

//获取交易明细
func (payment *Payment) getTradeInformationJson(items []Item) string {
	var goodsInfoStr string
	var totalQuantity int
	for _, item := range items {
		qty := strconv.Itoa(item.Qty)
		goodsInfoStr += item.Name + "^" + qty + "|"
		totalQuantity += item.Qty
	}
	goodsInfoStr = goodsInfoStr[0 : len(goodsInfoStr)-1]

	tradeInfo := TradeInformation{
		BusinessType:  businessTypeSalesGoods,
		GoodsInfo:     goodsInfoStr,
		TotalQuantity: totalQuantity,
	}
	bty, _ := json.Marshal(tradeInfo)

	return string(bty)
}

//获取支付金额
func (payment *Payment) getTotalFee(arg PayArg) (totalFee string) {
	//支付金额处理
	if arg.Currency == KRW || arg.Currency == JPY {
		totalFee = strconv.FormatFloat(arg.TotalFee, 'f', 0, 64)
	} else {
		totalFee = strconv.FormatFloat(arg.TotalFee, 'f', 2, 64)
	}
	return totalFee
}

//获取产品代码
func (payment *Payment) getProductCode(orderSource string) (productCode string) {
	if orderSource == Web {
		productCode = newOverseasSeller
	} else if orderSource == MobileWeb {
		productCode = newWapOverseasSeller
	}
	return productCode
}
