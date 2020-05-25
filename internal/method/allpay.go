package method

import (
	"errors"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/config"
	"payment_demo/internal/common/defs"
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
	AllpayTradeWay = "allpay.trade_way"
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

func (allpay *Allpay) Notify(query, methodCode string) (notifyRsp defs.NotifyRsp, errCode int, err error) {
	var allpayNotifyRsp payment.NotifyRsp
	defer func() {
		//记录日志
		logrus.Info("order id:%v,org:%v,method:%v,notify data:%+v",
			allpayNotifyRsp.OrderId, code.AllpayOrg, methodCode, allpayNotifyRsp.Rsp)
	}()

	md5key := config.GetInstance().GetString(AllpayMd5Key)
	if md5key == "" {
		logrus.Errorf(code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return notifyRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	allpayNotifyRsp, errCode, err = new(payment.Notify).Validate(query, md5key)
	if err != nil {
		return notifyRsp, errCode, err
	}
	notifyRsp.TradeNo = allpayNotifyRsp.TradeNo
	notifyRsp.Status = allpayNotifyRsp.Status
	notifyRsp.OrderId = allpayNotifyRsp.OrderId
	notifyRsp.Message = "OK"

	return notifyRsp, 0, nil
}

func (allpay *Allpay) Callback(query, methodCode string) (callbackRsp defs.CallbackRsp, errCode int, err error) {
	var allpayCallbackRsp payment.CallbackRsp
	defer func() {
		//记录日志
		logrus.Info("order id:%v,org:%v,method:%v,callback data:%+v",
			allpayCallbackRsp.OrderId, code.AllpayOrg, methodCode, allpayCallbackRsp.Rsp)
	}()

	md5key := config.GetInstance().GetString(AllpayMd5Key)
	if md5key == "" {
		logrus.Errorf(code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return callbackRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	allpayCallbackRsp, errCode, err = new(payment.Callback).Validate(query, md5key)
	if err != nil {
		return callbackRsp, errCode, err
	}

	callbackRsp.Status = allpayCallbackRsp.Status
	callbackRsp.OrderId = allpayCallbackRsp.OrderId

	return callbackRsp, 0, nil
}

func (allpay *Allpay) Trade(orderId, methodCode string) (tradeRsp defs.TradeRsp, errCode int, err error) {
	var allpayTradeRsp payment.TradeRsp
	defer func() {
		//记录日志
		logrus.Info("order id:%v,org:%v,method:%v,trade data:%+v",
			allpayTradeRsp.OrderId, code.AllpayOrg, methodCode, allpayTradeRsp.Rsp)
	}()

	merchant := config.GetInstance().GetString(AllpayMerchant)
	if merchant == "" {
		logrus.Errorf(code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return tradeRsp, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	acqId := config.GetInstance().GetString(AllpayAcqId)
	if acqId == "" {
		logrus.Errorf(code.AcqIdNotExistsErrMessage+",errCode:%v,err:%v", code.AcqIdNotExistsErrCode)
		return tradeRsp, code.AcqIdNotExistsErrCode, errors.New(code.AcqIdNotExistsErrMessage)
	}

	md5key := config.GetInstance().GetString(AllpayMd5Key)
	if md5key == "" {
		logrus.Errorf(code.Md5KeyNotExistsErrMessage+",errCode:%v,err:%v", code.Md5KeyNotExistsErrCode)
		return tradeRsp, code.Md5KeyNotExistsErrCode, errors.New(code.Md5KeyNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(AllpayTradeWay)
	if gateWay == "" {
		logrus.Errorf(code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return tradeRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	paymentSchema, errCode, err := allpay.getPaymentSchema(methodCode)
	if err != nil {
		return tradeRsp, errCode, err
	}

	tradeArg := payment.TradeArg{
		OrderNum:      orderId,
		MerId:         merchant,
		AcqId:         acqId,
		Md5Key:        md5key,
		PayWay:        gateWay,
		PaymentSchema: paymentSchema,
	}
	allpayTradeRsp, errCode, err = new(payment.Trade).Search(tradeArg)
	if err != nil {
		return tradeRsp, errCode, err
	}
	tradeRsp.OrderId = allpayTradeRsp.OrderId
	tradeRsp.TradeNo = allpayTradeRsp.TradeNo
	tradeRsp.Status = allpayTradeRsp.Status

	return tradeRsp, 0, nil
}
