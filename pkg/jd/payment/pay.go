package payment

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"payment_demo/internal/common/code"
	"payment_demo/pkg/jd/util"
	"time"
)

const (
	Practicality               = "0" //实物
	BusinessServiceConsumeCode = "100006"
	TimeLayout                 = "20060102150405"
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

func (payment *Payment) CreatePaymentForm(arg PayArg) (form string, errCode int, err error) {
	goodsInfoBytes, err := json.Marshal(arg.GoodsInfo)
	if err != nil {
		return form, code.PaymentGoodsInfoFormatErrCode, errors.New(code.PaymentGoodsInfoFormatErrMessage)
	}
	kjInfoBytes, err := json.Marshal(arg.KjInfo)
	if err != nil {
		return form, code.PaymentKjInfoFormatErrCode, errors.New(code.PaymentKjInfoFormatErrMessage)
	}

	paramMap := map[string]string{
		"version":        Version,
		"merchant":       arg.Merchant,
		"tradeNum":       arg.TradeNum,
		"tradeName":      arg.TradeName,
		"tradeTime":      time.Now().Format(TimeLayout),
		"amount":         fmt.Sprintf("%v", arg.Amount),
		"orderType":      Practicality,
		"currency":       arg.Currency,
		"callbackUrl":    arg.CallbackUrl,
		"notifyUrl":      arg.NotifyUrl,
		"userId":         arg.UserId,
		"expireTime":     fmt.Sprintf("%v", arg.ExpireTime),
		"goodsInfo":      string(goodsInfoBytes),
		"kjInfo":         string(kjInfoBytes),
		"bizTp":          BusinessServiceConsumeCode,
		"settleCurrency": arg.TransCurrency,
	}

	//签名
	paramMap["sign"], err = payment.getSign(paramMap, arg.PrivateKey)
	if err != nil {
		return form, code.PaymentSignErrCode, errors.New(code.PaymentSignErrMessage)
	}

	//加密
	paramMap, err = payment.encrypt3DES(paramMap, arg.DesKey)
	if err != nil {
		return form, 10153, errors.New("发起支付，数据加密错误")
	}

	//生成form表单
	form = payment.getForm(paramMap, arg.PayWay)

	return form, 0, nil
}

func (payment *Payment) getForm(paramMap map[string]string, payWay string) string {
	payUrl := "<form action='" + payWay + "' method='post' id='pay_form'>"
	for k, v := range paramMap {
		payUrl += "<input value='" + v + "' name='" + k + "' type='hidden'/>"
	}
	payUrl += "</form>"
	payUrl += "<script>var form = document.getElementById('pay_form');form.submit()</script>"
	return payUrl
}

func (payment *Payment) encrypt3DES(paramMap map[string]string, desKey string) (map[string]string, error) {
	desKey = util.BASE64DecodeStr(desKey)
	for k, v := range paramMap {
		if k == "sign" || k == "merchant" || k == "version" {
			continue
		}
		encrypt, err := util.TripleEcbDesEncrypt([]byte(v), []byte(desKey))
		if err != nil {
			return paramMap, err
		}
		paramMap[k] = util.DecimalByteSlice2HexString(encrypt)
	}
	return paramMap, nil
}

func (payment *Payment) getSign(paramMap map[string]string, privateKey string) (string, error) {
	payString := util.GetPayString(paramMap)
	sha256 := util.HaSha256(payString)
	sign, err := util.SignPKCS1v15([]byte(sha256), []byte(privateKey))
	if err != nil {
		return "", err
	}
	base64String := base64.StdEncoding.EncodeToString(sign)

	return base64String, nil
}
