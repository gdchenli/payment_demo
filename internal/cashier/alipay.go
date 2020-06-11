package cashier

import (
	"errors"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/config"
	"payment_demo/internal/common/defs"

	"github.com/gdchenli/pay/dialects/alipay/payment"
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

func (alipay *Alipay) Closed(arg defs.Closed) (closedRsp defs.ClosedRsp, errCode int, err error) {
	logrus.Errorf("org:alipay,"+code.NotSupportPaymentMethodErrMessage+",errCode:%v,err:%v", code.NotSupportPaymentMethodErrCode)
	closedRsp.Status = true
	return closedRsp, 0, nil
}

func (alipay *Alipay) OrderQrCode(arg defs.Order) (form string, errCode int, err error) {
	logrus.Errorf("org:jd,"+code.NotSupportPaymentMethodErrMessage+",errCode:%v,err:%v", code.NotSupportPaymentMethodErrCode)
	return "", 0, nil
}

func (alipay *Alipay) getPayArg(arg defs.Order) (payArg payment.PayArg, errCode int, err error) {
	merchant := config.GetInstance().GetString(AlipayMerchant)
	if merchant == "" {
		logrus.Errorf("org:alipay,"+code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return payArg, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}
	md5key := config.GetInstance().GetString(AlipayMd5Key)
	if md5key == "" {
		logrus.Errorf("org:alipay,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return payArg, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(AlipayGateWay)
	if gateWay == "" {
		logrus.Errorf("org:alipay,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return payArg, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	notifyUrl := config.GetInstance().GetString(AlipayNotifyUrl)
	if notifyUrl == "" {
		logrus.Errorf("org:alipay,"+code.NotifyUrlNotExistsErrMessage+",errCode:%v,err:%v", code.NotifyUrlNotExistsErrCode)
		return payArg, code.NotifyUrlNotExistsErrCode, errors.New(code.NotifyUrlNotExistsErrMessage)
	}

	callbackUrl := config.GetInstance().GetString(AlipayReturnUrl)
	if callbackUrl == "" {
		logrus.Errorf("org:alipay,"+code.CallbackUrlNotExistsErrMessage+",errCode:%v,err:%v", code.CallbackUrlNotExistsErrCode)
		return payArg, code.CallbackUrlNotExistsErrCode, errors.New(code.CallbackUrlNotExistsErrMessage)
	}

	supplier := config.GetInstance().GetString(AlipaySupplier)

	expireTime := config.GetInstance().GetString(AlipayTimeout)
	if expireTime == "" {
		logrus.Errorf("org:alipay,"+code.ExpireTimeNotExistsErrMessage+",errCode:%v,err:%v", code.ExpireTimeNotExistsErrCode)
		return payArg, code.ExpireTimeNotExistsErrCode, errors.New(code.ExpireTimeNotExistsErrMessage)
	}

	referUrl := config.GetInstance().GetString(AlipayReferUrl)
	userAgentType := alipay.getUserAgentType(arg.UserAgentType)
	transCurrency := config.GetInstance().GetString(AlipayTransCurrency)
	payWay := config.GetInstance().GetString(AlipayPayWay)

	payArg = payment.PayArg{
		Merchant:      merchant,
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
		Md5Key:        md5key,
		Items:         []payment.Item{{Name: "test", Qty: 1}},
		TransCurrency: transCurrency,
		UserAgentType: userAgentType,
		PayWay:        payWay,
	}

	return payArg, 0, nil
}

func (alipay *Alipay) Pay(arg defs.Order) (form string, errCode int, err error) {
	payArg, errCode, err := alipay.getPayArg(arg)
	if err != nil {
		return form, errCode, err
	}

	if arg.UserAgentType == code.AlipayMiniProgramUserAgentType {
		return new(payment.Payment).CreateAmpPayStr(payArg)
	} else {
		return new(payment.Payment).CreateForm(payArg)
	}
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
	var alipayNotifyRsp payment.NotifyRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,notify data:%+v",
			alipayNotifyRsp.OrderId, code.AlipayOrg, methodCode, alipayNotifyRsp.Rsp)
	}()

	merchant := config.GetInstance().GetString(AlipayMerchant)
	if merchant == "" {
		logrus.Errorf("org:alipay,"+code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return notifyRsp, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(AlipayGateWay)
	if gateWay == "" {
		logrus.Errorf("org:alipay,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return notifyRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	md5key := config.GetInstance().GetString(AlipayMd5Key)
	if md5key == "" {
		logrus.Errorf("org:alipay,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return notifyRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	notifyArg := payment.NotifyArg{Merchant: merchant, Md5Key: md5key, GateWay: gateWay}
	alipayNotifyRsp, errCode, err = new(payment.Notify).Validate(query, notifyArg)
	if err != nil {
		return notifyRsp, errCode, err
	}

	notifyRsp.TradeNo = alipayNotifyRsp.TradeNo
	notifyRsp.Status = alipayNotifyRsp.Status
	notifyRsp.OrderId = alipayNotifyRsp.OrderId
	notifyRsp.Message = "success"
	notifyRsp.RmbFee = alipayNotifyRsp.RmbFee
	notifyRsp.Rate = alipayNotifyRsp.Rate

	return notifyRsp, 0, nil
}

func (alipay *Alipay) Callback(query, methodCode string) (callbackRsp defs.CallbackRsp, errCode int, err error) {
	var alipayCallbackRsp payment.CallbackRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,callback data:%+v",
			alipayCallbackRsp.OrderId, code.AlipayOrg, methodCode, alipayCallbackRsp.Rsp)
	}()

	md5key := config.GetInstance().GetString(AlipayMd5Key)
	if md5key == "" {
		logrus.Errorf("org:alipay,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return callbackRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	alipayCallbackRsp, errCode, err = new(payment.Callback).Validate(query, md5key)
	if err != nil {
		return callbackRsp, errCode, err
	}

	callbackRsp.Status = alipayCallbackRsp.Status
	callbackRsp.OrderId = alipayCallbackRsp.OrderId

	return callbackRsp, 0, nil
}

func (alipay *Alipay) Trade(orderId, methodCode, currency string, totalFee float64) (tradeRsp defs.TradeRsp, errCode int, err error) {
	var alipayTradeRsp payment.TradeRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,callback data:%+v",
			alipayTradeRsp.OrderId, code.AlipayOrg, methodCode, alipayTradeRsp.Rsp)
	}()

	merchant := config.GetInstance().GetString(AlipayMerchant)
	if merchant == "" {
		logrus.Errorf("org:alipay,"+code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return tradeRsp, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}
	md5key := config.GetInstance().GetString(AlipayMd5Key)
	if md5key == "" {
		logrus.Errorf("org:alipay,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return tradeRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(AlipayGateWay)
	if gateWay == "" {
		logrus.Errorf("org:alipay,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return tradeRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	tradeArg := payment.TradeArg{
		Merchant:   merchant,
		OutTradeNo: orderId,
		Md5Key:     md5key,
		TradeWay:   gateWay,
		TotalFee:   totalFee,
	}
	alipayTradeRsp, errCode, err = new(payment.Trade).Search(tradeArg)
	if err != nil {
		return tradeRsp, 0, nil
	}
	tradeRsp.Status = alipayTradeRsp.Status
	tradeRsp.OrderId = alipayTradeRsp.OrderId
	tradeRsp.TradeNo = alipayTradeRsp.TradeNo
	tradeRsp.Rate = alipayTradeRsp.Rate
	tradeRsp.RmbFee = alipayTradeRsp.RmbFee

	return tradeRsp, 0, nil
}
