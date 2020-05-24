package method

import (
	"errors"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/config"
	"payment_demo/pkg/allpay/payment"
	"payment_demo/pkg/allpay/util"

	"github.com/sirupsen/logrus"
)

const (
	AllpayFrontUrl = "allpay.front_url"
	AllpayBackUrl  = "allpay.back_url"
	AllpayMerchant = "allpay.merchant"
	AllpayAcqId    = "allpay.acq_id"
	AllpayMd5Key   = "allpay.md5_key"
	AllpayPcPayWay = "allpay.pc_pay_way"
	AllpayH5PayWay = "allpay.h5_pay_way"
	AllpayTimeout  = "allpay.timeout"
)

type Allpay struct{}

type AllpayArg struct {
	OrderId       string  `json:"order_id"`
	TotalFee      float64 `json:"total_fee"`
	Currency      string  `json:"currency"`
	UserId        string  `json:"user_id"`
	MethodCode    string  `json:"method_code"`
	UserAgentType int     `json:"user_agent_type"`
}

func (allpay *Allpay) Submit(arg AllpayArg) (form string, errCode int, err error) {
	detailInfo := []payment.DetailInfo{
		{
			GoodsName: util.SpecialReplace("test goods name"),
			Quantity:  1,
		},
	}
	paymentSchema, errCode, err := allpay.getPaymentSchema(arg.MethodCode)
	if err != nil {
		return form, errCode, err
	}

	gateWay := allpay.getPayWay(arg.UserAgentType)
	if gateWay == "" {
		logrus.Errorf(code.GateWayNotExitsErrMessage+",errCode:%v,err:%v", code.GateWayNotExitsErrCode)
		return form, code.GateWayNotExitsErrCode, errors.New(code.GateWayNotExitsErrMessage)
	}
	tradeFrom := allpay.getTradeFrom(arg.MethodCode, arg.UserAgentType)

	payArg := payment.PayArg{
		OrderNum:      arg.OrderId,
		OrderAmount:   arg.TotalFee,
		FrontUrl:      config.GetInstance().GetString(AllpayFrontUrl),
		BackUrl:       config.GetInstance().GetString(AllpayBackUrl),
		MerId:         config.GetInstance().GetString(AllpayMerchant),
		AcqId:         config.GetInstance().GetString(AllpayAcqId),
		PaymentSchema: paymentSchema,
		GoodsInfo:     arg.OrderId,
		DetailInfo:    detailInfo,
		PayWay:        gateWay,
		Md5Key:        config.GetInstance().GetString(AllpayMd5Key),
		TradeFrom:     tradeFrom,
		OrderCurrency: arg.Currency,
		Timeout:       config.GetInstance().GetString(AllpayTimeout),
	}
	form, errCode, err = new(payment.Payment).CreateForm(payArg)
	if err != nil {
		return form, errCode, err
	}
	return form, 0, nil
}

func (allpay *Allpay) getPaymentSchema(methodCode string) (string, int, error) {
	switch methodCode {
	case "alipay_payment":
		//支付宝
		return "AP", 0, nil
	case "vtpayment_payment":
		//银联
		return "UP", 0, nil
	}
	return "", code.NotSupportPaymentMethodErrCode, errors.New(code.NotSupportPaymentMethodErrMessage)
}

func (allpay *Allpay) getTradeFrom(methodCode string, userAgentType int) string {
	if methodCode == "alipay_payment" {
		switch userAgentType {
		case 1:
			return "WEB"
		case 2:
			return "JSAPI"
		}
	}

	if methodCode == "vtpayment_payment" {
		return "H5"
	}

	return ""
}

func (allpay *Allpay) getPayWay(userAgentType int) string {
	switch userAgentType {
	case 1:
		return config.GetInstance().GetString(AllpayPcPayWay)
	case 2:
		return config.GetInstance().GetString(AllpayH5PayWay)
	}
	return config.GetInstance().GetString(AllpayPcPayWay)
}
