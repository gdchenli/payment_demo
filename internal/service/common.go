package service

import (
	"payment_demo/api/response"
	"payment_demo/api/validate"
	"payment_demo/internal/common/consts"
	"payment_demo/pkg/payment/alipay"
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
type CloseTradeHandler func(configParamMap map[string]string, arg validate.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error)                     //关闭交易
type UploadLogisticsHandler func(configParamMap map[string]string, arg validate.UploadLogisticsReq) (uploadLogisticsRsp response.UploadLogisticsRsp, errCode int, err error) //上传物流

//配置
type ConfigCodeHandler func() []string

//发起支付
var submitMap map[string]SumbitHandler
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

	configCodeMap = map[string]ConfigCodeHandler{
		consts.AlipayOrg + ".payment": alipayPayment.GetConfigCode, //发起支付配置
		consts.AlipayOrg + ".notify":  alipayNotify.GetConfigCode,  //异步通知配置
		consts.AlipayOrg + ".verify":  alipayVerify.GetConfigCode,  //同步通知
		consts.AlipayOrg + ".trade":   alipayTrade.GetConfigCode,   //交易查询
	}

	submitMap = map[string]SumbitHandler{
		consts.AlipayOrg: alipayPayment.CreatePayUrl,
	}

	ampSubmitMap = map[string]SumbitHandler{
		consts.AlipayOrg: alipayPayment.CreateAmpPayStr,
	}

	appSubmitMap = map[string]AppSumbitHandler{
		consts.AlipayOrg: alipayPayment.CreateAppPayStr,
	}

	notifyMap = map[string]NotifyHandler{
		consts.AlipayOrg: alipayNotify.Validate,
	}

	verifyMap = map[string]VerifyHandler{
		consts.AlipayOrg: alipayVerify.Validate,
	}

	searchTradeMap = map[string]SearchTradeHandler{
		consts.AlipayOrg: alipayTrade.Search,
	}

}

//发起支付
func getSubmitHandler(orgCode string, userAgentType int) SumbitHandler {
	switch userAgentType {
	case consts.WebUserAgentType, consts.MobileUserAgentType:
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

//配置读取
func getConfigCodeHandler(orgCode string) ConfigCodeHandler {
	return configCodeMap[orgCode]
}
