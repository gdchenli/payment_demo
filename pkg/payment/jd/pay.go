package jd

import (
	"crypto"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"payment_demo/api/validate"
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

type Payment struct{}

type PayArg struct {
	Merchant      string      `json:"merchant"`      //商户号
	TradeNum      string      `json:"tradeNum"`      //订单编号
	TradeName     string      `json:"tradeName"`     //交易名称
	Amount        int64       `json:"amount"`        //交易金额，单位分，大于0
	Currency      string      `json:"currency"`      //货币种类
	CallbackUrl   string      `json:"callbackUrl"`   //支付成功后跳转路径
	NotifyUrl     string      `json:"notifyUrl"`     //异步通知页面地址
	UserId        string      `json:"userId"`        //用户账号
	ExpireTime    string      `json:"expireTime"`    //订单失效时长，单位：秒，失效后则不能再支付，默认失效时间为604800秒(7天)，最大失效时间为7776000秒（90天），超过则按90天计算
	GoodsInfo     []GoodsInfo `json:"goodsInfo"`     //商品信息
	KjInfo        KjInfo      `json:"kInfo"`         //业务信息
	BizTp         string      `json:"bizTp"`         //通道业务类型
	PrivateKey    string      `json:"privateKey"`    //商户私钥
	DesKey        string      `json:"desKey"`        //des key
	TransCurrency string      `json:"transCurrency"` //结算币种
	PayWay        string      `json:"pay_way"`
}

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

func (payment *Payment) CreatePayUrl(paramMap map[string]string, order validate.Order) (form string, errCode int, err error) {
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
	paramMap["sign"], err = payment.getSign(paramMap, privateKey)
	if err != nil {
		logrus.Errorf("org:jd,"+PaySignErrMessage+",order id %v,errCode:%v,err:%v", order.OrderId, PaySignErrCode, err.Error())
		return form, PaySignErrCode, errors.New(PaySignErrMessage)
	}

	//加密
	paramMap, err = payment.encrypt3DES(paramMap, desKey)
	if err != nil {
		logrus.Errorf("org:jd,"+PayEncryptErrMessage+",order id %v,errCode:%v,err:%v", order.OrderId, PayEncryptErrCode, err.Error())
		return form, PayEncryptErrCode, errors.New(PayEncryptErrMessage)
	}

	//生成form表单
	form = payment.buildPayUrl(paramMap, payWay)

	return form, 0, nil
}

func (payment *Payment) buildPayUrl(paramMap map[string]string, gateWay string) (payUrl string) {
	payUrl = "<form action='" + gateWay + "' method='post' id='pay_form'>"
	for k, v := range paramMap {
		payUrl += "<input value='" + v + "' name='" + k + "' type='hidden'/>"
	}
	payUrl += "</form>"
	payUrl += "<script>var form = document.getElementById('pay_form');form.submit()</script>"
	return payUrl
}

func (payment *Payment) encrypt3DES(paramMap map[string]string, desKey string) (map[string]string, error) {
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

func (payment *Payment) getSign(paramMap map[string]string, privateKey string) (string, error) {
	payString := GetSortString(paramMap)
	sha256 := HaSha256(payString)
	sign, err := SignPKCS1v15([]byte(sha256), []byte(privateKey), crypto.Hash(0))
	if err != nil {
		return "", err
	}
	base64String := base64.StdEncoding.EncodeToString(sign)

	return base64String, nil
}

func (payment *Payment) GetConfigCode() []string {
	return []string{
		"merchant", "callbackUrl", "notifyUrl", "settleCurrency", "expireTime",
		"des_key", "pc_pay_way", "h5_pay_way", "private_key",
	}
}
