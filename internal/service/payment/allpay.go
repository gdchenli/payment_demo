package payment

import (
	"errors"
	"fmt"
	"payment_demo/api/validate"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/config"
	"payment_demo/internal/common/response"
	"time"

	"github.com/gdchenli/pay/dialects/allpay/payment"
	"github.com/gdchenli/pay/dialects/allpay/util"
	"github.com/sirupsen/logrus"
)

const (
	AllpayFrontUrl = "allpay.front_url"
	AllpayBackUrl  = "allpay.back_url"
	AllpayMerchant = "allpay.merchant"
	AllpayAcqId    = "allpay.acq_id"
	AllpayMd5Key   = "allpay.md5_key"
	AllpayGateWay  = "allpay.gate_way"
	AllpayTimeout  = "allpay.timeout"
	AllpaySapiWay  = "allpay.sapi_way"
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

func (allpay *Allpay) getPayArg(arg validate.Order) (payArg payment.PayArg, errCode int, err error) {
	merchant := config.GetInstance().GetString(AllpayMerchant)
	if merchant == "" {
		logrus.Errorf("org:allpay,"+code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return payArg, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	notifyUrl := config.GetInstance().GetString(AllpayBackUrl)
	if notifyUrl == "" {
		logrus.Errorf("org:allpay,"+code.NotifyUrlNotExistsErrMessage+",errCode:%v,err:%v", code.NotifyUrlNotExistsErrCode)
		return payArg, code.NotifyUrlNotExistsErrCode, errors.New(code.NotifyUrlNotExistsErrMessage)
	}

	callbackUrl := config.GetInstance().GetString(AllpayFrontUrl)
	if callbackUrl == "" {
		logrus.Errorf("org:allpay,"+code.CallbackUrlNotExistsErrMessage+",errCode:%v,err:%v", code.CallbackUrlNotExistsErrCode)
		return payArg, code.CallbackUrlNotExistsErrCode, errors.New(code.CallbackUrlNotExistsErrMessage)
	}

	expireTime := config.GetInstance().GetString(AllpayTimeout)
	if expireTime == "" {
		logrus.Errorf("org:allpay,"+code.ExpireTimeNotExistsErrMessage+",errCode:%v,err:%v", code.ExpireTimeNotExistsErrCode)
		return payArg, code.ExpireTimeNotExistsErrCode, errors.New(code.ExpireTimeNotExistsErrMessage)
	}

	md5key := config.GetInstance().GetString(AllpayMd5Key)
	if md5key == "" {
		logrus.Errorf("org:allpay,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return payArg, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	acqId := config.GetInstance().GetString(AllpayAcqId)
	if acqId == "" {
		logrus.Errorf("org:allpay,"+code.AcqIdNotExistsErrMessage+",errCode:%v,err:%v", code.AcqIdNotExistsErrCode)
		return payArg, code.AcqIdNotExistsErrCode, errors.New(code.AcqIdNotExistsErrMessage)
	}

	detailInfo := []payment.DetailInfo{
		{
			GoodsName: util.SpecialReplace("test goods name" + time.Now().Format(payment.TimeLayout)),
			Quantity:  1,
		},
	}
	paymentSchema, errCode, err := allpay.getPaymentSchema(arg.MethodCode)
	if err != nil {
		return payArg, errCode, err
	}

	gateWay := allpay.getPayWay(arg.UserAgentType)
	if gateWay == "" {
		logrus.Errorf("org:allpay,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return payArg, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	tradeFrom := allpay.getTradeFrom(arg.MethodCode, arg.UserAgentType)

	payArg = payment.PayArg{
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

	return payArg, 0, nil
}

func (allpay *Allpay) Pay(arg validate.Order) (form string, errCode int, err error) {
	payArg, errCode, err := allpay.getPayArg(arg)
	if err != nil {
		return form, errCode, err
	}

	if arg.UserAgentType == code.AlipayMiniProgramUserAgentType {
		return new(payment.Payment).CreateAmpPayStr(payArg)
	} else {
		return new(payment.Payment).CreateForm(payArg)
	}
}

func (allpay *Allpay) getPaymentSchema(methodCode string) (string, int, error) {
	switch methodCode {
	case code.AlipayMethod:
		return AlipayPaymentSchema, 0, nil
	case code.UnionpayMethod:
		return UpPaymentSchema, 0, nil
	}
	return "", code.NotSupportPaymentMethodErrCode, errors.New(code.NotSupportPaymentMethodErrMessage)
}

func (allpay *Allpay) getTradeFrom(methodCode string, userAgentType int) string {
	if methodCode == code.AlipayMethod {
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
		return config.GetInstance().GetString(AllpayGateWay)
	case 2:
		return config.GetInstance().GetString(AllpayGateWay)
	}
	return config.GetInstance().GetString(AllpayGateWay)
}

func (allpay *Allpay) Notify(query, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error) {
	var allpayNotifyRsp payment.NotifyRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,notify data:%+v",
			allpayNotifyRsp.OrderId, code.AllpayOrg, methodCode, allpayNotifyRsp.Rsp)
	}()

	md5key := config.GetInstance().GetString(AllpayMd5Key)
	if md5key == "" {
		logrus.Errorf("org:allpay,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return notifyRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	merchant := config.GetInstance().GetString(AllpayMerchant)
	if merchant == "" {
		logrus.Errorf("org:allpay,"+code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return notifyRsp, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	sapiGateWay := config.GetInstance().GetString(AllpaySapiWay)
	if sapiGateWay == "" {
		logrus.Errorf("org:allpay,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return notifyRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	notifyArg := payment.NotifyArg{
		MerId:       merchant,
		Md5Key:      md5key,
		SapiGateWay: sapiGateWay,
	}
	fmt.Println("notifyArg.SapiGateWay", notifyArg.SapiGateWay)
	allpayNotifyRsp, errCode, err = new(payment.Notify).Validate(query, notifyArg)
	if err != nil {
		return notifyRsp, errCode, err
	}
	notifyRsp.TradeNo = allpayNotifyRsp.TradeNo
	notifyRsp.Status = allpayNotifyRsp.Status
	notifyRsp.OrderId = allpayNotifyRsp.OrderId
	notifyRsp.Message = "OK"
	notifyRsp.RmbFee = allpayNotifyRsp.RmbFee
	notifyRsp.Rate = allpayNotifyRsp.Rate

	return notifyRsp, 0, nil
}

func (allpay *Allpay) Verify(query, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) {
	var allpayCallbackRsp payment.CallbackRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,callback data:%+v",
			allpayCallbackRsp.OrderId, code.AllpayOrg, methodCode, allpayCallbackRsp.Rsp)
	}()

	md5key := config.GetInstance().GetString(AllpayMd5Key)
	if md5key == "" {
		logrus.Errorf("org:allpay,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return verifyRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	allpayCallbackRsp, errCode, err = new(payment.Callback).Validate(query, md5key)
	if err != nil {
		return verifyRsp, errCode, err
	}

	verifyRsp.Status = allpayCallbackRsp.Status
	verifyRsp.OrderId = allpayCallbackRsp.OrderId

	return verifyRsp, 0, nil
}

func (allpay *Allpay) SearchTrade(orderId, methodCode, currency string, totalFee float64) (searchtradeRsp response.SearchTradeRsp, errCode int, err error) {
	var allpayTradeRsp payment.TradeRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,trade data:%+v",
			allpayTradeRsp.OrderId, code.AllpayOrg, methodCode, allpayTradeRsp.Rsp)
	}()

	merchant := config.GetInstance().GetString(AllpayMerchant)
	if merchant == "" {
		logrus.Errorf("org:allpay,"+code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return searchtradeRsp, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	acqId := config.GetInstance().GetString(AllpayAcqId)
	if acqId == "" {
		logrus.Errorf("org:allpay,"+code.AcqIdNotExistsErrMessage+",errCode:%v,err:%v", code.AcqIdNotExistsErrCode)
		return searchtradeRsp, code.AcqIdNotExistsErrCode, errors.New(code.AcqIdNotExistsErrMessage)
	}

	md5key := config.GetInstance().GetString(AllpayMd5Key)
	if md5key == "" {
		logrus.Errorf("org:allpay,"+code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return searchtradeRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	tradeGateWay := config.GetInstance().GetString(AllpayGateWay)
	if tradeGateWay == "" {
		logrus.Errorf("org:allpay,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return searchtradeRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	sapiGateWay := config.GetInstance().GetString(AllpaySapiWay)
	if sapiGateWay == "" {
		logrus.Errorf("org:allpay,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return searchtradeRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	paymentSchema, errCode, err := allpay.getPaymentSchema(methodCode)
	if err != nil {
		return searchtradeRsp, errCode, err
	}

	tradeArg := payment.TradeArg{
		OrderNum:      orderId,
		MerId:         merchant,
		AcqId:         acqId,
		Md5Key:        md5key,
		TradeGateWay:  tradeGateWay,
		PaymentSchema: paymentSchema,
		SapiGateWay:   sapiGateWay,
		TotalFee:      totalFee,
		Currency:      currency,
	}
	allpayTradeRsp, errCode, err = new(payment.Trade).Search(tradeArg)
	if err != nil {
		return searchtradeRsp, errCode, err
	}
	searchtradeRsp.OrderId = allpayTradeRsp.OrderId
	searchtradeRsp.TradeNo = allpayTradeRsp.TradeNo
	searchtradeRsp.Status = allpayTradeRsp.Status
	searchtradeRsp.Rate = allpayTradeRsp.Rate
	searchtradeRsp.RmbFee = allpayTradeRsp.RmbFee

	return searchtradeRsp, 0, nil
}
