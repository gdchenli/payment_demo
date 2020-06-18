package allpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"payment_demo/api/request"
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

type DetailInfo struct {
	GoodsName string `json:"goods_name"`
	Quantity  int    `json:"quantity"`
}

func (allpay *Allpay) CreatePayUrl(configParamMap map[string]string, order request.Order) (payUrl string, errCode int, err error) {
	gateWay := getPayGateWay(configParamMap["gate_way"])
	delete(configParamMap, "gata_way")

	paramMap, errCode, err := getPayParamMap(configParamMap, order)
	if err != nil {
		return payUrl, errCode, err
	}

	payUrl = buildPayUrl(paramMap, gateWay)

	return payUrl, 0, nil
}

func (allpay *Allpay) CreateAmpPayStr(configParamMap map[string]string, order request.Order) (payStr string, errCode int, err error) {
	gateWay := getPayGateWay(configParamMap["gate_way"])
	delete(configParamMap, "gata_way")

	paramMap, errCode, err := getPayParamMap(configParamMap, order)
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

func getPayParamMap(paramMap map[string]string, order request.Order) (map[string]string, int, error) {
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
	paramMap["paymentSchema"] = getPaymentSchema(order.MethodCode)
	paramMap["goodsInfo"] = order.OrderId
	paramMap["detailInfo"] = detailInfo
	paramMap["transTime"] = transTime
	paramMap["signType"] = MD5SignType
	paramMap["tradeFrom"] = getPayTradeFrom(order.MethodCode, order.UserAgentType)
	paramMap["merReserve"] = ""

	md5key := paramMap["md5_key"]
	delete(paramMap, "md5_key")
	paramMap["signature"] = getPaySign(paramMap, md5key)

	return paramMap, 0, nil
}
func getPaymentSchema(methodCode string) string {
	switch methodCode {
	case consts.AlipayMethod:
		return ApSchema
	case consts.UnionpayMethod:
		return UpSchema
	default:
		return ""
	}
}

func getPayTradeFrom(methodCode string, userAgentType int) string {
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

func getPaySign(paramMap map[string]string, signKey string) string {
	sortString := GetSortString(paramMap)
	return Md5(sortString + signKey)
}

func buildPayUrl(paramMap map[string]string, gateWay string) (payUrl string) {
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

func getPayGateWay(gateWay string) string {
	return gateWay + PayRoute
}

func (allpay *Allpay) GetPayConfigCode() []string {
	return []string{
		"merID", "frontURL", "backURL", "acqID", "timeout",
		"md5_key", "gate_way",
	}
}
