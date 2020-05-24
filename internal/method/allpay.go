package method

import (
	"errors"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/config"
	"payment_demo/pkg/allpay/payment"
	"payment_demo/pkg/allpay/util"
	"time"

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

const (
	AlipayWebTradeFrom         = "WEB"
	AlipayMobileTradeFrom      = "JSAPI"
	AlipayMiniProgramTradeFrom = "APPLET"
	UpTradeFrom                = "H5"
	AppTradeFrom               = "APP"
	AlipayPaymentSchema        = "AP"
	UpPaymentSchema            = "UP"
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
	merchant := config.GetInstance().GetString(AllpayMerchant)
	if merchant == "" {
		logrus.Errorf(code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return form, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	notifyUrl := config.GetInstance().GetString(AllpayBackUrl)
	if notifyUrl == "" {
		logrus.Errorf(code.NotifyUrlNotExistsErrMessage+",errCode:%v,err:%v", code.NotifyUrlNotExistsErrCode)
		return form, code.NotifyUrlNotExistsErrCode, errors.New(code.NotifyUrlNotExistsErrMessage)
	}

	callbackUrl := config.GetInstance().GetString(AllpayFrontUrl)
	if callbackUrl == "" {
		logrus.Errorf(code.CallbackUrlNotExistsErrMessage+",errCode:%v,err:%v", code.CallbackUrlNotExistsErrCode)
		return form, code.CallbackUrlNotExistsErrCode, errors.New(code.CallbackUrlNotExistsErrMessage)
	}

	expireTime := config.GetInstance().GetString(AllpayTimeout)
	if expireTime == "" {
		logrus.Errorf(code.ExpireTimeNotExistsErrMessage+",errCode:%v,err:%v", code.ExpireTimeNotExistsErrCode)
		return form, code.ExpireTimeNotExistsErrCode, errors.New(code.ExpireTimeNotExistsErrMessage)
	}

	md5key := config.GetInstance().GetString(AllpayMd5Key)
	if md5key == "" {
		logrus.Errorf(code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return form, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	acqId := config.GetInstance().GetString(AllpayAcqId)
	if acqId == "" {
		logrus.Errorf(code.AcqIdNotExistsErrMessage+",errCode:%v,err:%v", code.AcqIdNotExistsErrCode)
		return form, code.AcqIdNotExistsErrCode, errors.New(code.AcqIdNotExistsErrMessage)
	}

	detailInfo := []payment.DetailInfo{
		{
			GoodsName: util.SpecialReplace("test goods name" + time.Now().Format(payment.TimeLayout)),
			Quantity:  1,
		},
	}
	paymentSchema, errCode, err := allpay.getPaymentSchema(arg.MethodCode)
	if err != nil {
		return form, errCode, err
	}

	gateWay := allpay.getPayWay(arg.UserAgentType)
	if gateWay == "" {
		logrus.Errorf(code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return form, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}
	tradeFrom := allpay.getTradeFrom(arg.MethodCode, arg.UserAgentType)

	payArg := payment.PayArg{
		OrderNum:      arg.OrderId,
		OrderAmount:   arg.TotalFee,
		FrontUrl:      callbackUrl,
		BackUrl:       notifyUrl,
		MerId:         merchant,
		AcqId:         acqId,
		PaymentSchema: paymentSchema,
		GoodsInfo:     arg.OrderId,
		DetailInfo:    detailInfo,
		PayWay:        gateWay,
		Md5Key:        md5key,
		TradeFrom:     tradeFrom,
		OrderCurrency: arg.Currency,
		Timeout:       expireTime,
	}
	form, errCode, err = new(payment.Payment).CreateForm(payArg)
	if err != nil {
		return form, errCode, err
	}
	return form, 0, nil
}

func (allpay *Allpay) getPaymentSchema(methodCode string) (string, int, error) {
	switch methodCode {
	case code.AliapayMethod:
		return AlipayPaymentSchema, 0, nil
	case code.UnionpayMethod:
		return UpPaymentSchema, 0, nil
	}
	return "", code.NotSupportPaymentMethodErrCode, errors.New(code.NotSupportPaymentMethodErrMessage)
}

func (allpay *Allpay) getTradeFrom(methodCode string, userAgentType int) string {
	if methodCode == code.AliapayMethod {
		switch userAgentType {
		case code.WebUserAgentType:
			return AlipayWebTradeFrom
		case code.MobileUserAgentType:
			return AlipayMobileTradeFrom
		case code.AlipayMiniProgramUserAgentType:
			return AlipayMiniProgramTradeFrom
		case code.AndroidAppUserAgentType, code.IOSAppUserAgentType:
			return AppTradeFrom
		}
	}

	if methodCode == code.UnionpayMethod {
		return UpTradeFrom
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
