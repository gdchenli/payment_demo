package cashier

import (
	"errors"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/config"
	"payment_demo/internal/common/defs"

	"github.com/gdchenli/pay/dialects/epayments/payment"
	"github.com/sirupsen/logrus"
)

const (
	EpaymentsTimeout       = "epayments.timeout"
	EpaymentsMerchant      = "epayments.merchant"
	EpaymentsMd5Key        = "epayments.md5_key"
	EpaymentsGateWay       = "epayments.gate_way"
	EpaymentsNotifyUrl     = "epayments.notify_url"
	EpaymentsReturnUrl     = "epayments.return_url"
	EpaymentsTransCurrency = "epayments.trans_currency"
)

type Epayments struct{}

func (e *Epayments) AmpSubmit(arg defs.Order) (form string, errCode int, err error) {
	return form, code.NotSupportPaymentMethodErrCode, errors.New(code.NotSupportPaymentMethodErrMessage)
}

func (e *Epayments) getPayArg(order defs.Order) (payArg payment.PayArg, errCode int, err error) {
	merchant := config.GetInstance().GetString(EpaymentsMerchant)
	if merchant == "" {
		logrus.Errorf("org:epayments,"+code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return payArg, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	md5key := config.GetInstance().GetString(EpaymentsMd5Key)
	if md5key == "" {
		logrus.Errorf("org:epayments,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return payArg, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(EpaymentsGateWay)
	if gateWay == "" {
		logrus.Errorf("org:epayments,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return payArg, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	notifyUrl := config.GetInstance().GetString(EpaymentsNotifyUrl)
	if notifyUrl == "" {
		logrus.Errorf("org:epayments,"+code.NotifyUrlNotExistsErrMessage+",errCode:%v,err:%v", code.NotifyUrlNotExistsErrCode)
		return payArg, code.NotifyUrlNotExistsErrCode, errors.New(code.NotifyUrlNotExistsErrMessage)
	}

	callbackUrl := config.GetInstance().GetString(EpaymentsReturnUrl)
	if callbackUrl == "" {
		logrus.Errorf("org:epayments,"+code.CallbackUrlNotExistsErrMessage+",errCode:%v,err:%v", code.CallbackUrlNotExistsErrCode)
		return payArg, code.CallbackUrlNotExistsErrCode, errors.New(code.CallbackUrlNotExistsErrMessage)
	}

	expireTime := config.GetInstance().GetString(EpaymentsTimeout)
	if expireTime == "" {
		logrus.Errorf("org:epayments,"+code.ExpireTimeNotExistsErrMessage+",errCode:%v,err:%v", code.ExpireTimeNotExistsErrCode)
		return payArg, code.ExpireTimeNotExistsErrCode, errors.New(code.ExpireTimeNotExistsErrMessage)
	}

	transCurrency := config.GetInstance().GetString(EpaymentsTransCurrency)
	paymentChannels := e.getPaymentChannels(order.MethodCode)

	payArg = payment.PayArg{
		MerchantId:      merchant,
		NotifyUrl:       notifyUrl,
		ReturnUrl:       callbackUrl,
		ValidMins:       expireTime,
		IncrementId:     order.OrderId,
		GrandTotal:      order.TotalFee,
		Currency:        order.Currency,
		GateWay:         gateWay,
		Md5Key:          md5key,
		TransCurrency:   transCurrency,
		PaymentChannels: paymentChannels,
	}

	return payArg, 0, nil
}

func (e *Epayments) OrderSubmit(order defs.Order) (form string, errCode int, err error) {
	payArg, errCode, err := e.getPayArg(order)
	if err != nil {
		return form, errCode, err
	}
	return new(payment.Payment).CreateForm(payArg)
}

func (e *Epayments) OrderQrCode(order defs.Order) (form string, errCode int, err error) {
	payArg, errCode, err := e.getPayArg(order)
	if err != nil {
		return form, errCode, err
	}

	return new(payment.Payment).CreateQrCode(payArg)
}

//获取支付通道
func (e *Epayments) getPaymentChannels(methodCode string) (paymentChannels string) {
	if methodCode == code.WechatMethod {
		paymentChannels = payment.ChannelWechat
	} else if methodCode == code.AlipayMethod {
		paymentChannels = payment.ChannelAlipay
	}
	return paymentChannels
}

func (e *Epayments) Notify(query, methodCode string) (notifyRsp defs.NotifyRsp, errCode int, err error) {
	var epaymentsNotifyRsp payment.NotifyRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,notify data:%+v",
			epaymentsNotifyRsp.OrderId, code.EpaymentsOrg, methodCode, epaymentsNotifyRsp.Rsp)
	}()

	md5key := config.GetInstance().GetString(EpaymentsMd5Key)
	if md5key == "" {
		logrus.Errorf("org:epayments,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return notifyRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	epaymentsNotifyRsp, errCode, err = new(payment.Notify).Validate(query, md5key)
	if err != nil {
		return notifyRsp, errCode, err
	}
	notifyRsp.TradeNo = epaymentsNotifyRsp.TradeNo
	notifyRsp.Status = epaymentsNotifyRsp.Status
	notifyRsp.OrderId = epaymentsNotifyRsp.OrderId
	notifyRsp.Message = "success"

	return notifyRsp, 0, nil
}

func (e *Epayments) Callback(query, methodCode string) (callbackRsp defs.CallbackRsp, errCode int, err error) {
	var epaymentsCallbackRsp payment.CallbackRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,callback data:%+v",
			epaymentsCallbackRsp.OrderId, code.EpaymentsOrg, methodCode, epaymentsCallbackRsp.Rsp)
	}()

	md5key := config.GetInstance().GetString(EpaymentsMd5Key)
	if md5key == "" {
		logrus.Errorf("org:epayments,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return callbackRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	epaymentsCallbackRsp, errCode, err = new(payment.Callback).Validate(query, md5key)
	if err != nil {
		return callbackRsp, errCode, err
	}

	callbackRsp.Status = epaymentsCallbackRsp.Status
	callbackRsp.OrderId = epaymentsCallbackRsp.OrderId

	return callbackRsp, 0, nil
}

func (e *Epayments) Trade(orderId, methodCode string) (tradeRsp defs.TradeRsp, errCode int, err error) {
	var epaymentsTradeRsp payment.TradeRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,callback data:%+v",
			epaymentsTradeRsp.OrderId, code.EpaymentsOrg, methodCode, epaymentsTradeRsp.Rsp)
	}()

	merchant := config.GetInstance().GetString(EpaymentsMerchant)
	if merchant == "" {
		logrus.Errorf("org:epayments,"+code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return tradeRsp, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}
	md5key := config.GetInstance().GetString(EpaymentsMd5Key)
	if md5key == "" {
		logrus.Errorf("org:epayments,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return tradeRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(EpaymentsGateWay)
	if gateWay == "" {
		logrus.Errorf("org:epayments,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return tradeRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	tradeArg := payment.TradeArg{
		Merchant:    merchant,
		IncrementId: orderId,
		Md5Key:      md5key,
		TradeWay:    gateWay,
	}
	epaymentsTradeRsp, errCode, err = new(payment.Trade).Search(tradeArg)
	if err != nil {
		return tradeRsp, 0, nil
	}
	tradeRsp.Status = epaymentsTradeRsp.Status
	tradeRsp.OrderId = epaymentsTradeRsp.OrderId
	tradeRsp.TradeNo = epaymentsTradeRsp.TradeNo

	return tradeRsp, 0, nil
}

func (e *Epayments) Closed(arg defs.Closed) (closedRsp defs.ClosedRsp, errCode int, err error) {
	return closedRsp, code.NotSupportPaymentMethodErrCode, errors.New(code.NotSupportPaymentMethodErrMessage)
}
