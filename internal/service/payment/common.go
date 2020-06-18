package payment

import (
	"payment_demo/api/response"
	"payment_demo/internal/common/request"
	"payment_demo/pkg/payment/alipay"
	"payment_demo/pkg/payment/allpay"
	"payment_demo/pkg/payment/consts"
	"payment_demo/pkg/payment/epayments"
	"payment_demo/pkg/payment/jd"
)

//发起支付
type SumbitHandler func(configParamMap map[string]string, order request.Order) (url string, errCode int, err error)                //pc、h5、支付宝小程序
type WmpSumbitHandler func(configParamMap map[string]string, order request.Order) (wmRsp response.WmpRsp, errCode int, err error)  //微信小程序
type AppSumbitHandler func(configParamMap map[string]string, order request.Order) (appRsp response.AppRsp, errCode int, err error) //App

//支付通知
type NotifyHandler func(configParamMap map[string]string, query, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error) //异步通知
type VerifyHandler func(configParamMap map[string]string, query, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) //同步通知

//交易信息
type SearchTradeHandler func(configParamMap map[string]string, req request.SearchTradeReq) (searchTradeRsp response.SearchTradeRsp, errCode int, err error)                 //交易查询
type CloseTradeHandler func(configParamMap map[string]string, req request.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error)                     //关闭交易
type UploadLogisticsHandler func(configParamMap map[string]string, req request.UploadLogisticsReq) (uploadLogisticsRsp response.UploadLogisticsRsp, errCode int, err error) //上传物流

//配置
type ConfigCodeHandler func() []string

//发起支付
var submitMap map[string]SumbitHandler
var transferSubmitMap map[string]SumbitHandler
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
	alipayPayment := new(alipay.Alipay)
	allpayPayment := new(allpay.Allpay)
	epaymentsPayment := new(epayments.Epayments)
	jdPayment := new(jd.Jd)

	configCodeMap = map[string]ConfigCodeHandler{
		consts.AlipayOrg + ".payment": alipayPayment.GetPayConfigCode,         //发起支付配置
		consts.AlipayOrg + ".notify":  alipayPayment.GetNotifyConfigCode,      //异步通知配置
		consts.AlipayOrg + ".verify":  alipayPayment.GetVerifyConfigCode,      //同步通知
		consts.AlipayOrg + ".trade":   alipayPayment.GetSearchTradeConfigCode, //交易查询

		consts.AllpayOrg + ".payment": allpayPayment.GetPayConfigCode,         //发起支付配置
		consts.AllpayOrg + ".notify":  allpayPayment.GetNotifyConfigCode,      //异步通知配置
		consts.AllpayOrg + ".verify":  allpayPayment.GetVerifyConfigCode,      //同步通知
		consts.AllpayOrg + ".trade":   allpayPayment.GetSearchTradeConfigCode, //交易查询

		consts.EpaymentsOrg + ".payment": epaymentsPayment.GetPayConfigCode,         //发起支付配置
		consts.EpaymentsOrg + ".notify":  epaymentsPayment.GetNotifyConfigCode,      //异步通知配置
		consts.EpaymentsOrg + ".verify":  epaymentsPayment.GetVerifyConfigCode,      //同步通知
		consts.EpaymentsOrg + ".trade":   epaymentsPayment.GetSearchTradeConfigCode, //交易查询

		consts.JdOrg + ".payment":   jdPayment.GetPayConfigCode,             //发起支付配置
		consts.JdOrg + ".notify":    jdPayment.GetNotifyConfigCode,          //异步通知配置
		consts.JdOrg + ".verify":    jdPayment.GetVerifyConfigCode,          //同步通知
		consts.JdOrg + ".trade":     jdPayment.GetSearchTradeConfigCode,     //交易查询
		consts.JdOrg + ".close":     jdPayment.GetCloseTradeConfigCode,      //交易关闭
		consts.JdOrg + ".logistics": jdPayment.GetUploadLogisticsConfigCode, //物流上传
	}

	submitMap = map[string]SumbitHandler{
		consts.AlipayOrg:    alipayPayment.CreatePayUrl,
		consts.AllpayOrg:    allpayPayment.CreatePayUrl,
		consts.EpaymentsOrg: epaymentsPayment.CreatePayUrl,
		consts.JdOrg:        jdPayment.CreatePayUrl,
	}

	transferSubmitMap = map[string]SumbitHandler{
		consts.EpaymentsOrg: epaymentsPayment.CreateQrCode,
		consts.JdOrg:        jdPayment.CreatePayForm,
	}

	ampSubmitMap = map[string]SumbitHandler{
		consts.AlipayOrg: alipayPayment.CreateAmpPayStr,
		consts.AllpayOrg: allpayPayment.CreateAmpPayStr,
	}

	appSubmitMap = map[string]AppSumbitHandler{
		consts.AlipayOrg: alipayPayment.CreateAppPayStr,
	}

	notifyMap = map[string]NotifyHandler{
		consts.AlipayOrg:    alipayPayment.Notify,
		consts.AllpayOrg:    allpayPayment.Notify,
		consts.EpaymentsOrg: epaymentsPayment.Notify,
		consts.JdOrg:        jdPayment.Notify,
	}

	verifyMap = map[string]VerifyHandler{
		consts.AlipayOrg:    alipayPayment.Verify,
		consts.AllpayOrg:    allpayPayment.Verify,
		consts.EpaymentsOrg: epaymentsPayment.Verify,
		consts.JdOrg:        jdPayment.Verify,
	}

	searchTradeMap = map[string]SearchTradeHandler{
		consts.AlipayOrg:    alipayPayment.SearchTrade,
		consts.AllpayOrg:    allpayPayment.SearchTrade,
		consts.EpaymentsOrg: epaymentsPayment.SearchTrade,
		consts.JdOrg:        jdPayment.SearchTrade,
	}

	closeTradeMap = map[string]CloseTradeHandler{
		consts.JdOrg: jdPayment.CloseTrade,
	}

	uploadLogisticsMap = map[string]UploadLogisticsHandler{
		consts.JdOrg: jdPayment.Upload,
	}

}

//发起支付
func getSubmitHandler(orgCode string, istransfer bool) SumbitHandler {
	if istransfer {
		return transferSubmitMap[orgCode]
	}

	return submitMap[orgCode]
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

type PayHandler interface {
	//发起支付
	Sumbit(configParamMap map[string]string, order request.Order) (url string, errCode int, err error)                //pc、h5、支付宝小程序
	WmpSumbit(configParamMap map[string]string, order request.Order) (wmRsp response.WmpRsp, errCode int, err error)  //微信小程序
	AppSumbit(configParamMap map[string]string, order request.Order) (appRsp response.AppRsp, errCode int, err error) //App

	//支付通知
	Notify(configParamMap map[string]string, query, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error) //异步通知
	Verify(configParamMap map[string]string, query, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) //同步通知

	//交易信息
	SearchTrade(configParamMap map[string]string, req request.SearchTradeReq) (searchTradeRsp response.SearchTradeRsp, errCode int, err error)                 //交易查询
	CloseTrade(configParamMap map[string]string, req request.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error)                     //关闭交易
	UploadLogistics(configParamMap map[string]string, req request.UploadLogisticsReq) (uploadLogisticsRsp response.UploadLogisticsRsp, errCode int, err error) //上传物流

	//配置
	GetPayConfigCode() []string
	GetNotifyConfigCode() []string
	GetVerfiyConfigCode() []string
	GetSearchCodeConfigCode() []string
	GetCloseCodeConfigCode() []string
}
