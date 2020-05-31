package method

import (
	"errors"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/config"
	"payment_demo/internal/common/defs"
	"payment_demo/pkg/alipay/payment"

	"github.com/sirupsen/logrus"
)

type Alipay struct{}

const (
	AlipayTimeout       = "alipay.timeout"
	AlipayMerchant      = "alipay.merchant"
	AlipayMd5Key        = "alipay.md5_key"
	AlipayGateWay       = "alipay.gate_way"
	AlipayNotifyUrl     = "alipay.notify_url"
	AlipayReturnUrl     = "alipay.return_url"
	AlipaySupplier      = "alipay.supplier"
	AlipayReferUrl      = "alipay.refer_url"
	AlipayPayWay        = "alipay.pay_way"
	AlipayTransCurrency = "alipay.trans_currency"
)

func (alipay *Alipay) OrderSubmit(arg defs.Order) (form string, errCode int, err error) {
	merchant := config.GetInstance().GetString(AlipayMerchant)
	if merchant == "" {
		logrus.Errorf(code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return form, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}
	md5key := config.GetInstance().GetString(AlipayMd5Key)
	if md5key == "" {
		logrus.Errorf(code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return form, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(AlipayGateWay)
	if gateWay == "" {
		logrus.Errorf(code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return form, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	notifyUrl := config.GetInstance().GetString(AlipayNotifyUrl)
	if notifyUrl == "" {
		logrus.Errorf(code.NotifyUrlNotExistsErrMessage+",errCode:%v,err:%v", code.NotifyUrlNotExistsErrCode)
		return form, code.NotifyUrlNotExistsErrCode, errors.New(code.NotifyUrlNotExistsErrMessage)
	}

	callbackUrl := config.GetInstance().GetString(AlipayReturnUrl)
	if callbackUrl == "" {
		logrus.Errorf(code.CallbackUrlNotExistsErrMessage+",errCode:%v,err:%v", code.CallbackUrlNotExistsErrCode)
		return form, code.CallbackUrlNotExistsErrCode, errors.New(code.CallbackUrlNotExistsErrMessage)
	}

	supplier := config.GetInstance().GetString(AlipaySupplier)

	expireTime := config.GetInstance().GetString(AlipayTimeout)
	if expireTime == "" {
		logrus.Errorf(code.ExpireTimeNotExistsErrMessage+",errCode:%v,err:%v", code.ExpireTimeNotExistsErrCode)
		return form, code.ExpireTimeNotExistsErrCode, errors.New(code.ExpireTimeNotExistsErrMessage)
	}

	referUrl := config.GetInstance().GetString(AlipayReferUrl)
	userAgentType := alipay.getUserAgentType(arg.UserAgentType)
	transCurrency := config.GetInstance().GetString(AlipayTransCurrency)
	payWay := config.GetInstance().GetString(AlipayPayWay)

	payArg := payment.PayArg{
		Partner:       merchant,
		NotifyUrl:     notifyUrl,
		ReturnUrl:     callbackUrl,
		Body:          arg.OrderId,
		OutTradeNo:    arg.OrderId,
		TotalFee:      arg.TotalFee,
		Currency:      arg.Currency,
		Supplier:      supplier,
		TimeoutRule:   expireTime,
		ReferUrl:      referUrl,
		GateWay:       gateWay,
		Md5key:        md5key,
		Items:         []payment.Item{{Name: "test", Qty: 1}},
		TransCurrency: transCurrency,
		UserAgentType: userAgentType,
		PayWay:        payWay,
	}

	form, errCode, err = new(payment.Payment).CreateForm(payArg)
	if err != nil {
		return form, errCode, err
	}

	return form, 0, nil
}

func (alipay *Alipay) getUserAgentType(userAgentType int) string {
	switch userAgentType {
	case code.WebUserAgentType:
		return payment.Web
	case code.MobileUserAgentType:
		return payment.MobileWeb
	case code.AlipayMiniProgramUserAgentType:
		return payment.Amp
	}

	return ""
}

func (alipay *Alipay) Notify(query, methodCode string) (notifyRsp defs.NotifyRsp, errCode int, err error) {
	return notifyRsp, 0, nil
}

func (alipay *Alipay) Callback(query, methodCode string) (callbackRsp defs.CallbackRsp, errCode int, err error) {
	return callbackRsp, 0, nil
}

func (alipay *Alipay) Trade(orderId, methodCode string) (tradeRsp defs.TradeRsp, errCode int, err error) {
	return tradeRsp, 0, nil
}
