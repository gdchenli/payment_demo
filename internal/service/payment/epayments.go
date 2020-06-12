package payment

import (
	"errors"
	payment2 "payment_demo/api/validate/payment"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/config"
	payment3 "payment_demo/internal/common/response/payment"

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

func (e *Epayments) getPayArg(order payment2.Order) (payArg payment.PayArg, errCode int, err error) {
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

func (e *Epayments) Pay(order payment2.Order) (form string, errCode int, err error) {
	payArg, errCode, err := e.getPayArg(order)
	if err != nil {
		return form, errCode, err
	}

	if order.UserAgentType == code.MobileUserAgentType {
		return new(payment.Payment).CreateQrCode(payArg)
	} else {
		return new(payment.Payment).CreateForm(payArg)
	}
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

func (e *Epayments) Notify(query, methodCode string) (notifyRsp payment3.NotifyRsp, errCode int, err error) {
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
	notifyRsp.RmbFee = epaymentsNotifyRsp.RmbFee
	notifyRsp.Rate = epaymentsNotifyRsp.Rate

	return notifyRsp, 0, nil
}

func (e *Epayments) Verify(query, methodCode string) (verifyRsp payment3.VerifyRsp, errCode int, err error) {
	var epaymentsCallbackRsp payment.CallbackRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,callback data:%+v",
			epaymentsCallbackRsp.OrderId, code.EpaymentsOrg, methodCode, epaymentsCallbackRsp.Rsp)
	}()

	md5key := config.GetInstance().GetString(EpaymentsMd5Key)
	if md5key == "" {
		logrus.Errorf("org:epayments,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return verifyRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	epaymentsCallbackRsp, errCode, err = new(payment.Callback).Validate(query, md5key)
	if err != nil {
		return verifyRsp, errCode, err
	}

	verifyRsp.Status = epaymentsCallbackRsp.Status
	verifyRsp.OrderId = epaymentsCallbackRsp.OrderId

	return verifyRsp, 0, nil
}

func (e *Epayments) SearchTrade(orderId, methodCode, currency string, totalFee float64) (searchtradeRsp payment3.SearchTradeRsp, errCode int, err error) {
	var epaymentsTradeRsp payment.TradeRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,callback data:%+v",
			epaymentsTradeRsp.OrderId, code.EpaymentsOrg, methodCode, epaymentsTradeRsp.Rsp)
	}()

	merchant := config.GetInstance().GetString(EpaymentsMerchant)
	if merchant == "" {
		logrus.Errorf("org:epayments,"+code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return searchtradeRsp, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}
	md5key := config.GetInstance().GetString(EpaymentsMd5Key)
	if md5key == "" {
		logrus.Errorf("org:epayments,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return searchtradeRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(EpaymentsGateWay)
	if gateWay == "" {
		logrus.Errorf("org:epayments,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return searchtradeRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	tradeArg := payment.TradeArg{
		Merchant:    merchant,
		IncrementId: orderId,
		Md5Key:      md5key,
		TradeWay:    gateWay,
	}
	epaymentsTradeRsp, errCode, err = new(payment.Trade).Search(tradeArg)
	if err != nil {
		return searchtradeRsp, 0, nil
	}
	searchtradeRsp.Status = epaymentsTradeRsp.Status
	searchtradeRsp.OrderId = epaymentsTradeRsp.OrderId
	searchtradeRsp.TradeNo = epaymentsTradeRsp.TradeNo
	searchtradeRsp.RmbFee = epaymentsTradeRsp.RmbFee
	searchtradeRsp.Rate = epaymentsTradeRsp.Rate

	return searchtradeRsp, 0, nil
}

func (e *Epayments) CloseTrade(arg payment2.CloseTradeReq) (closeTradeRsp payment3.CloseTradeRsp, errCode int, err error) {
	logrus.Errorf("org:allpay,"+code.NotSupportPaymentMethodErrMessage+",errCode:%v,err:%v", code.NotSupportPaymentMethodErrCode)
	closeTradeRsp.Status = true
	return closeTradeRsp, 0, nil
}
