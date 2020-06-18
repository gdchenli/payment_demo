package jd

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"payment_demo/api/request"
	"payment_demo/pkg/payment/consts"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	Practicality               = "0" //实物
	BusinessServiceConsumeCode = "100006"
	TimeLayout                 = "20060102150405"
)

const (
	PayGoodsInfoFormatErrCode    = 10101
	PayGoodsInfoFormatErrMessage = "发起支付，商品数据转换失败"
	PayKjInfoFormatErrCode       = 10102
	PayKjInfoFormatErrMessage    = "发起支付，跨境数据转换失败"
	PaySignErrCode               = 10103
	PaySignErrMessage            = "发起支付，签名计算错误"
	PayEncryptErrCode            = 10104
	PayEncryptErrMessage         = "发起支付，数据加密错误"
)

type GoodsInfo struct {
	Id    string `json:"id"`    //商品编号
	Name  string `json:"name"`  //商品名称
	Price int64  `json:"price"` //商品单价，单位分
	Num   int    `json:"num"`   //商品数量
}

type KjInfo struct {
	GoodsSubmittedCustoms string `json:"goodsSubmittedCustoms"` //是否报关Y/N
	GoodsUnderBonded      string `json:"goodsUnderBonded"`      //是否保税货物项下付款Y/N
}

func (jd *Jd) CreatePayUrl(paramMap map[string]string, order request.Order) (form string, errCode int, err error) {
	marshal, _ := json.Marshal(order)
	reqParamMap := make(map[string]interface{})
	json.Unmarshal(marshal, &reqParamMap)
	values := url.Values{}
	for k, v := range reqParamMap {
		values.Add(k, fmt.Sprintf("%v", v))
	}
	return "/payment/form?" + values.Encode(), 0, nil
}

func (jd *Jd) CreatePayForm(paramMap map[string]string, order request.Order) (form string, errCode int, err error) {
	privateKey := paramMap["private_key"]
	delete(paramMap, "private_key")
	desKey := paramMap["des_key"]
	delete(paramMap, "des_key")

	var payWay string
	if order.UserAgentType == consts.WebUserAgentType {
		payWay = paramMap["pc_pay_way"]
	} else {
		payWay = paramMap["h5_pay_way"]
	}
	delete(paramMap, "pc_pay_way")
	delete(paramMap, "h5_pay_way")

	totalFee := fmt.Sprintf("%.f", order.TotalFee*100)
	totalFeeInt, _ := strconv.ParseInt(totalFee, 10, 64)
	date := time.Now().Format(TimeLayout)
	goodsInfoBytes, err := json.Marshal([]GoodsInfo{{Id: "test" + date, Name: "test" + date, Price: totalFeeInt, Num: 1}})
	if err != nil {
		logrus.Errorf("org:jd,"+PayGoodsInfoFormatErrMessage+",order id %v,errCode:%v,err:%v", order.OrderId, PayGoodsInfoFormatErrCode, err.Error())
		return form, PayGoodsInfoFormatErrCode, errors.New(PayGoodsInfoFormatErrMessage)
	}
	kjInfoBytes, err := json.Marshal(KjInfo{GoodsSubmittedCustoms: "N", GoodsUnderBonded: "N"})
	if err != nil {
		logrus.Errorf("org:jd,"+PayKjInfoFormatErrMessage+",errCode:%v,err:%v", PayKjInfoFormatErrCode, err.Error())
		return form, PayKjInfoFormatErrCode, errors.New(PayKjInfoFormatErrMessage)
	}
	paramMap["version"] = Version
	paramMap["tradeNum"] = order.OrderId
	paramMap["tradeName"] = order.OrderId
	paramMap["tradeTime"] = time.Now().Format(TimeLayout)
	paramMap["amount"] = totalFee
	paramMap["orderType"] = Practicality
	paramMap["currency"] = order.Currency
	paramMap["userId"] = order.UserId
	paramMap["goodsInfo"] = string(goodsInfoBytes)
	paramMap["kjInfo"] = string(kjInfoBytes)
	paramMap["bizTp"] = BusinessServiceConsumeCode

	//签名
	paramMap["sign"], err = getPaySign(paramMap, privateKey)
	if err != nil {
		logrus.Errorf("org:jd,"+PaySignErrMessage+",order id %v,errCode:%v,err:%v", order.OrderId, PaySignErrCode, err.Error())
		return form, PaySignErrCode, errors.New(PaySignErrMessage)
	}

	//加密
	paramMap, err = encrypt3DES(paramMap, desKey)
	if err != nil {
		logrus.Errorf("org:jd,"+PayEncryptErrMessage+",order id %v,errCode:%v,err:%v", order.OrderId, PayEncryptErrCode, err.Error())
		return form, PayEncryptErrCode, errors.New(PayEncryptErrMessage)
	}

	//生成form表单
	form = buildPayForm(paramMap, payWay)

	return form, 0, nil
}

func buildPayForm(paramMap map[string]string, gateWay string) (payUrl string) {
	payUrl = "<form action='" + gateWay + "' method='post' id='pay_form'>"
	for k, v := range paramMap {
		payUrl += "<input value='" + v + "' name='" + k + "' type='hidden'/>"
	}
	payUrl += "</form>"
	payUrl += "<script>var form = document.getElementById('pay_form');form.submit()</script>"
	return payUrl
}

func encrypt3DES(paramMap map[string]string, desKey string) (map[string]string, error) {
	desKey = BASE64DecodeStr(desKey)
	for k, v := range paramMap {
		if k == "sign" || k == "merchant" || k == "version" {
			continue
		}
		encrypt, err := TripleEcbDesEncrypt([]byte(v), []byte(desKey))
		if err != nil {
			return paramMap, err
		}
		paramMap[k] = DecimalByteSlice2HexString(encrypt)
	}
	return paramMap, nil
}

func getPaySign(paramMap map[string]string, privateKey string) (string, error) {
	payString := GetSortString(paramMap)
	sha256 := HaSha256(payString)
	sign, err := SignPKCS1v15([]byte(sha256), []byte(privateKey), crypto.Hash(0))
	if err != nil {
		return "", err
	}
	base64String := base64.StdEncoding.EncodeToString(sign)

	return base64String, nil
}

func (jd *Jd) GetPayConfigCode() []string {
	return []string{
		"merchant", "callbackUrl", "notifyUrl", "settleCurrency", "expireTime",
		"des_key", "pc_pay_way", "h5_pay_way", "private_key",
	}
}
