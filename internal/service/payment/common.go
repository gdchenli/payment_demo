package payment

import (
	"fmt"
	"payment_demo/api/response"
	"payment_demo/api/validate"
	"payment_demo/pkg/payment/alipay"
	"payment_demo/pkg/payment/allpay"
	"payment_demo/pkg/payment/consts"
	"payment_demo/pkg/payment/epayments"
	"payment_demo/pkg/payment/jd"
)

//发起支付
type SumbitHandler func(configParamMap map[string]string, order validate.Order) (url string, errCode int, err error)                //pc、h5、支付宝小程序
type WmpSumbitHandler func(configParamMap map[string]string, order validate.Order) (wmRsp response.WmpRsp, errCode int, err error)  //微信小程序
type AppSumbitHandler func(configParamMap map[string]string, order validate.Order) (appRsp response.AppRsp, errCode int, err error) //App

//支付通知
type NotifyHandler func(configParamMap map[string]string, query, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error) //异步通知
type VerifyHandler func(configParamMap map[string]string, query, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) //同步通知

//交易信息
type SearchTradeHandler func(configParamMap map[string]string, req validate.SearchTradeReq) (searchTradeRsp response.SearchTradeRsp, errCode int, err error)                 //交易查询
type CloseTradeHandler func(configParamMap map[string]string, req validate.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error)                     //关闭交易
type UploadLogisticsHandler func(configParamMap map[string]string, req validate.UploadLogisticsReq) (uploadLogisticsRsp response.UploadLogisticsRsp, errCode int, err error) //上传物流

//配置
type ConfigCodeHandler func() []string

//发起支付
var submitMap map[string]SumbitHandler
var qrCodeSubmitMap map[string]SumbitHandler
var ampSubmitMap map[string]SumbitHandler
var wmpSubmitMap map[string]WmpSumbitHandler
var appSubmitMap map[string]AppSumbitHandler

//支付通知
var notifyMap map[string]NotifyHandler
var verifyMap map[string]VerifyHandler

//交易信息
var searchTradeMap map[string]SearchTradeHandler
var closeTradeMap map[string]CloseTradeHandler
var uploadLogisticsMap map[string]UploadLogisticsHandler

//配置
var configCodeMap map[string]ConfigCodeHandler

func init() {
	alipayPayment := new(alipay.Payment)
	alipayNotify := new(alipay.Notify)
	alipayVerify := new(alipay.Verify)
	alipayTrade := new(alipay.Trade)

	allpayPayment := new(allpay.Payment)
	allpayNotify := new(allpay.Notify)
	allpayVerify := new(allpay.Verify)
	allpayTrade := new(allpay.Trade)

	epaymentsPayment := new(epayments.Payment)
	epaymentsNotify := new(epayments.Notify)
	epaymentsVerify := new(epayments.Verify)
	epaymentsTrade := new(epayments.Trade)

	jdPayment := new(jd.Payment)
	jdNotify := new(jd.Notify)
	jdVerify := new(jd.Verify)
	jdTrade := new(jd.Trade)
	jdClose := new(jd.Close)
	jdLogistics := new(jd.Logistics)

	configCodeMap = map[string]ConfigCodeHandler{
		consts.AlipayOrg + ".payment": alipayPayment.GetConfigCode, //发起支付配置
		consts.AlipayOrg + ".notify":  alipayNotify.GetConfigCode,  //异步通知配置
		consts.AlipayOrg + ".verify":  alipayVerify.GetConfigCode,  //同步通知
		consts.AlipayOrg + ".trade":   alipayTrade.GetConfigCode,   //交易查询

		consts.AllpayOrg + ".payment": allpayPayment.GetConfigCode, //发起支付配置
		consts.AllpayOrg + ".notify":  allpayNotify.GetConfigCode,  //异步通知配置
		consts.AllpayOrg + ".verify":  allpayVerify.GetConfigCode,  //同步通知
		consts.AllpayOrg + ".trade":   allpayTrade.GetConfigCode,   //交易查询

		consts.EpaymentsOrg + ".payment": epaymentsPayment.GetConfigCode, //发起支付配置
		consts.EpaymentsOrg + ".notify":  epaymentsNotify.GetConfigCode,  //异步通知配置
		consts.EpaymentsOrg + ".verify":  epaymentsVerify.GetConfigCode,  //同步通知
		consts.EpaymentsOrg + ".trade":   epaymentsTrade.GetConfigCode,   //交易查询

		consts.JdOrg + ".payment":   jdPayment.GetConfigCode,   //发起支付配置
		consts.JdOrg + ".notify":    jdNotify.GetConfigCode,    //异步通知配置
		consts.JdOrg + ".verify":    jdVerify.GetConfigCode,    //同步通知
		consts.JdOrg + ".trade":     jdTrade.GetConfigCode,     //交易查询
		consts.JdOrg + ".close":     jdClose.GetConfigCode,     //交易关闭
		consts.JdOrg + ".logistics": jdLogistics.GetConfigCode, //物流上传
	}

	submitMap = map[string]SumbitHandler{
		consts.AlipayOrg:    alipayPayment.CreatePayUrl,
		consts.AllpayOrg:    allpayPayment.CreatePayUrl,
		consts.EpaymentsOrg: epaymentsPayment.CreatePayUrl,
		consts.JdOrg:        jdPayment.CreatePayUrl,
	}

	qrCodeSubmitMap = map[string]SumbitHandler{
		consts.EpaymentsOrg: epaymentsPayment.CreateQrCode,
	}

	ampSubmitMap = map[string]SumbitHandler{
		consts.AlipayOrg: alipayPayment.CreateAmpPayStr,
		consts.AllpayOrg: allpayPayment.CreateAmpPayStr,
	}

	appSubmitMap = map[string]AppSumbitHandler{
		consts.AlipayOrg: alipayPayment.CreateAppPayStr,
	}

	notifyMap = map[string]NotifyHandler{
		consts.AlipayOrg:    alipayNotify.Validate,
		consts.AllpayOrg:    allpayNotify.Validate,
		consts.EpaymentsOrg: epaymentsNotify.Validate,
		consts.JdOrg:        jdNotify.Validate,
	}

	verifyMap = map[string]VerifyHandler{
		consts.AlipayOrg:    alipayVerify.Validate,
		consts.AllpayOrg:    allpayVerify.Validate,
		consts.EpaymentsOrg: epaymentsVerify.Validate,
		consts.JdOrg:        jdVerify.Validate,
	}

	searchTradeMap = map[string]SearchTradeHandler{
		consts.AlipayOrg:    alipayTrade.Search,
		consts.AllpayOrg:    allpayTrade.Search,
		consts.EpaymentsOrg: epaymentsTrade.Search,
		consts.JdOrg:        jdTrade.Search,
	}

	closeTradeMap = map[string]CloseTradeHandler{
		consts.JdOrg: jdClose.Trade,
	}

	uploadLogisticsMap = map[string]UploadLogisticsHandler{
		consts.JdOrg: jdLogistics.Upload,
	}

}

//发起支付
func getSubmitHandler(orgCode string, userAgentType int) SumbitHandler {
	fmt.Printf("%v\n", submitMap)
	switch userAgentType {
	case consts.WebUserAgentType:
		if orgCode == consts.EpaymentsOrg {
			return qrCodeSubmitMap[orgCode]
		}
		return submitMap[orgCode]
	case consts.MobileUserAgentType:
		return submitMap[orgCode]
	case consts.AlipayMiniProgramUserAgentType:
		return ampSubmitMap[orgCode]
	}

	return nil
}

func getAppSubmitHandler(orgCode string) AppSumbitHandler {
	return appSubmitMap[orgCode]
}

//异步通知
func getNotifyHandler(orgCode string) NotifyHandler {
	return notifyMap[orgCode]
}

//同步通知
func getVerifyHandler(orgCode string) VerifyHandler {
	return verifyMap[orgCode]
}

//交易查询
func getSeachTradeHandler(orgCode string) SearchTradeHandler {
	return searchTradeMap[orgCode]
}

//关闭交易
func getCloseTradeHandler(orgCode string) CloseTradeHandler {
	return closeTradeMap[orgCode]
}

//上传物流信息
func getUploadLogisticsHandler(orgCode string) UploadLogisticsHandler {
	return uploadLogisticsMap[orgCode]
}

//配置读取
func getConfigCodeHandler(orgCode string) ConfigCodeHandler {
	return configCodeMap[orgCode]
}
