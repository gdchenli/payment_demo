package alipay

import (
	"encoding/json"
	"net/url"
	"payment_demo/api/response"
	"payment_demo/api/validate"
	"payment_demo/pkg/payment/consts"
	"strconv"
)

const (
	pcServiceType          = "create_forex_trade"      //PC端支付类型
	wapServiceType         = "create_forex_trade_wap"  //移动端支付类型
	newPay                 = "2"                       //新接口
	newOverseasSeller      = "NEW_OVERSEAS_SELLER"     //海外销售产品代码
	newWapOverseasSeller   = "NEW_WAP_OVERSEAS_SELLER" //海外移动端销售产品代码
	businessTypeSalesGoods = 4
)

type Payment struct{}

type PayArg struct {
	Merchant      string  `json:"merchant"`          //PartnerId
	NotifyUrl     string  `json:"notify_url"`        //支付结果异步通知到该地址
	ReturnUrl     string  `json:"return_url"`        //支付结果异步通知到该地址
	Body          string  `json:"body"`              //主体
	OutTradeNo    string  `json:"out_trade_no"`      //订单号
	TotalFee      float64 `json:"total_fee"`         //订单金额(CNY)
	Currency      string  `json:"currency"`          //订单币种
	Supplier      string  `json:"supplier"`          //供应商
	TimeoutRule   string  `json:"timeout_rule"`      //超时时间,如:12h,10m
	ReferUrl      string  `json:"refer_url"`         //商家url(站点url)
	GateWay       string  `json:"gate_way"`          //网关地址
	Md5Key        string  `json:"md5_key"`           //密钥
	TransCurrency string  `json:"trade_information"` //结算币种
	UserAgentType string  `json:"user_agent_type"`   //订单客户端 PC:web  手机端：web_mobile
	PayWay        string  `json:"pay_way"`           //版本 1旧版本 2新版本
	Items         []Item  `json:"item"`
}

type Item struct {
	Name string `form:"name" json:"name" `
	Qty  int    `form:"qty_ordered" json:"qty_ordered"`
}
type TradeInformation struct {
	BusinessType  int    `json:"business_type"`
	GoodsInfo     string `json:"goods_info"`
	TotalQuantity int    `json:"total_quantity"`
}

func (payment *Payment) getParamMap(paramMap map[string]string, order validate.Order) map[string]string {
	paramMap["service"] = payment.getServiceType(order.UserAgentType)
	paramMap["_input_charset"] = CharsetUTF8
	paramMap["subject"] = order.OrderId
	paramMap["body"] = order.OrderId
	paramMap["out_trade_no"] = order.OrderId
	paramMap["total_fee"] = payment.getTotalFee(order.Currency, order.TotalFee)
	paramMap["trade_information"] = payment.getTradeInformationJson([]Item{{Name: "test", Qty: 1}})

	//超时时间
	if paramMap["timeout_rule"] == "" {
		delete(paramMap, "timeout_rule")
	}
	//判断新旧版本
	if paramMap["pay_way"] == newPay {
		paramMap["product_code"] = payment.getProductCode(order.UserAgentType)
	}
	delete(paramMap, "pay_way")

	//人民币币种替换字段
	if paramMap["currency"] == CNY {
		paramMap["rmb_fee"] = paramMap["total_fee"]
		delete(paramMap, "total_fee")
	}

	md5Key := paramMap["md5_key"]
	delete(paramMap, "md5_key")

	//签名
	payString := GetSortString(paramMap)
	paramMap["sign"] = Md5(payString + md5Key)
	paramMap["sign_type"] = SignTypeMD5

	return paramMap
}

func (payment *Payment) CreateAmpPayStr(configParamMap map[string]string, order validate.Order) (payString string, errCode int, err error) {
	delete(configParamMap, "gate_way")

	paramMap := payment.getParamMap(configParamMap, order)

	return GetSortString(paramMap), 0, nil
}

func (payment *Payment) CreatePayUrl(configParamMap map[string]string, order validate.Order) (url string, errCode int, err error) {
	geteWay := configParamMap["gate_way"]
	delete(configParamMap, "gate_way")

	paramMap := payment.getParamMap(configParamMap, order)

	return payment.buildPayUrl(paramMap, geteWay), 0, nil
}

func (payment *Payment) CreateAppPayStr(paramMap map[string]string, order validate.Order) (appRsp response.AppRsp, errCode int, err error) {
	return appRsp, 0, nil
}

func (payment *Payment) buildPayUrl(paramMap map[string]string, gateWay string) (payUrl string) {
	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}
	payUrl = gateWay + "?" + values.Encode()
	return payUrl
}

//获取服务类型
func (payment *Payment) getServiceType(orderSource int) (serviceType string) {
	if orderSource == 1 {
		serviceType = pcServiceType
	} else if orderSource == 2 {
		serviceType = wapServiceType
	} else if orderSource == 7 {
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
func (payment *Payment) getTotalFee(currency string, totalFeeF float64) (totalFee string) {
	//支付金额处理
	if currency == KRW || currency == JPY {
		totalFee = strconv.FormatFloat(totalFeeF, 'f', 0, 64)
	} else {
		totalFee = strconv.FormatFloat(totalFeeF, 'f', 2, 64)
	}
	return totalFee
}

//获取产品代码
func (payment *Payment) getProductCode(orderSource int) (productCode string) {
	if orderSource == consts.WebUserAgentType {
		productCode = newOverseasSeller
	} else if orderSource == consts.MobileUserAgentType || orderSource == consts.AlipayMiniProgramUserAgentType {
		productCode = newWapOverseasSeller
	}
	return productCode
}

func (payment *Payment) GetConfigCode() []string {
	return []string{
		"partner", "notify_url", "return_url", "supplier", "refer_url", "currency", "supplier", "timeout_rule",
		"pay_way", "md5_key", "gate_way",
	}
}
