package alipay

import (
	"encoding/json"
	"net/url"
	"payment_demo/pkg/payment/common"
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

type Item struct {
	Name string `form:"name" json:"name" `
	Qty  int    `form:"qty_ordered" json:"qty_ordered"`
}
type TradeInformation struct {
	BusinessType  int    `json:"business_type"`
	GoodsInfo     string `json:"goods_info"`
	TotalQuantity int    `json:"total_quantity"`
}

func getParamMap(paramMap map[string]string, order common.OrderArg) map[string]string {
	paramMap["service"] = getServiceType(order.UserAgentType)
	paramMap["_input_charset"] = CharsetUTF8
	paramMap["subject"] = order.OrderId
	paramMap["body"] = order.OrderId
	paramMap["out_trade_no"] = order.OrderId
	paramMap["total_fee"] = getTotalFee(order.Currency, order.TotalFee)
	paramMap["trade_information"] = getTradeInformationJson([]Item{{Name: "test", Qty: 1}})

	//超时时间
	if paramMap["timeout_rule"] == "" {
		delete(paramMap, "timeout_rule")
	}
	//判断新旧版本
	if paramMap["pay_way"] == newPay {
		paramMap["product_code"] = getProductCode(order.UserAgentType)
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

func (alipay *Alipay) CreateAmpPayStr(configParamMap map[string]string, order common.OrderArg) (payString string, errCode int, err error) {
	delete(configParamMap, "gate_way")

	paramMap := getParamMap(configParamMap, order)

	return GetSortString(paramMap), 0, nil
}

func (alipay *Alipay) CreatePayUrl(configParamMap map[string]string, order common.OrderArg) (url string, errCode int, err error) {
	geteWay := configParamMap["gate_way"]
	delete(configParamMap, "gate_way")

	paramMap := getParamMap(configParamMap, order)

	return buildPayUrl(paramMap, geteWay), 0, nil
}

func (alipay *Alipay) CreateAppPayStr(paramMap map[string]string, order common.OrderArg) (appRsp common.AppPayRsp, errCode int, err error) {
	return appRsp, 0, nil
}

func buildPayUrl(paramMap map[string]string, gateWay string) (payUrl string) {
	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}
	payUrl = gateWay + "?" + values.Encode()
	return payUrl
}

//获取服务类型
func getServiceType(orderSource int) (serviceType string) {
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
func getTradeInformationJson(items []Item) string {
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
func getTotalFee(currency string, totalFeeF float64) (totalFee string) {
	//支付金额处理
	if currency == KRW || currency == JPY {
		totalFee = strconv.FormatFloat(totalFeeF, 'f', 0, 64)
	} else {
		totalFee = strconv.FormatFloat(totalFeeF, 'f', 2, 64)
	}
	return totalFee
}

//获取产品代码
func getProductCode(orderSource int) (productCode string) {
	if orderSource == consts.WebUserAgentType {
		productCode = newOverseasSeller
	} else if orderSource == consts.MobileUserAgentType || orderSource == consts.AlipayMiniProgramUserAgentType {
		productCode = newWapOverseasSeller
	}
	return productCode
}

func (alipay *Alipay) GetPayConfigCode() []string {
	return []string{
		"partner", "notify_url", "return_url", "supplier", "refer_url", "currency", "supplier", "timeout_rule",
		"pay_way", "md5_key", "gate_way",
	}
}
