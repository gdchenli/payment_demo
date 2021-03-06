package allpay

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"payment_demo/pkg/curl"
)

const (
	RateSearchNetErrCode       = 10601
	RateSearchNetErrMessage    = "汇率查询网络错误"
	RateSearchFormatErrCode    = 10602
	RateSearchFormatErrMessage = "汇率查询结果格式错误"
)

const (
	RateRoute = "/sapi/v1/get_exchange_rate"
)

type Rate struct{}

type RateArg struct {
	MerId                  string `json:"merchant"`
	PaymentSchema          string `json:"paymentSchema"`
	OriginalCurrencyCode   string `json:"original_currency_code"`
	ConversionCurrencyCode string `json:"conversion_currency_code"`
	Md5Key                 string `json:"md5_key"`
	GateWay                string `json:"gate_way"`
}

func (allpay *Allpay) SearchRate(arg RateArg) (float64, int, error) {
	paramMap := map[string]string{
		"pid":                 arg.MerId,
		"issuer":              getRateIssuer(arg.PaymentSchema),
		"original_currency":   arg.OriginalCurrencyCode,
		"conversion_currency": arg.ConversionCurrencyCode,
		"sign_type":           SignTypeSHA256,
	}
	sortString := GetSortString(paramMap)
	paramMap["sign"] = Hsha256(sortString + arg.Md5Key)

	values := url.Values{}
	for k, v := range paramMap {
		values.Add(k, v)
	}

	fmt.Println(getRateGateWay(arg.GateWay) + "?" + values.Encode())
	returnBytes, err := curl.GetJSONReturnByte(getRateGateWay(arg.GateWay) + "?" + values.Encode())
	if err != nil {
		return 0, RateSearchNetErrCode, errors.New(RateSearchNetErrMessage)
	}

	rateRspMap := make(map[string]interface{})
	if err = json.Unmarshal(returnBytes, &rateRspMap); err != nil {
		return 0, RateSearchFormatErrCode, errors.New(RateSearchFormatErrMessage)
	}

	rateMap := rateRspMap["data"].(map[string]interface{})
	exchangeRate := rateMap["exchange_rate"].(string)

	val, err := strconv.ParseFloat(exchangeRate, 64)
	if err != nil {
		return 0, RateSearchFormatErrCode, errors.New(RateSearchFormatErrMessage)
	}

	return val, 0, nil
}

func getRateIssuer(paymentSchema string) string {
	switch paymentSchema {
	case UpSchema:
		return UpIssuer
	case ApSchema:
		return ApIssuer
	case WxSchema:
		return WxIssuer
	default:
		return ""
	}
}

func getRateGateWay(gateWay string) string {
	return gateWay + RateRoute
}
