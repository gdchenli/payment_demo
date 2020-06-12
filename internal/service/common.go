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
type NotifyHandler func(query, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error) //异步通知
type VerifyHandler func(query, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) //同步通知

//交易信息
type SearchTradeHandler func(orderId, methodCode, curreny string, totalFee float64) (searchTradeRsp response.SearchTradeRsp, errCode int, err error) //交易查询
type CloseTradeHandler func(arg validate.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error)                               //关闭交易
type UploadLogisticsHandler func(arg validate.UploadLogisticsReq) (uploadLogisticsRsp response.UploadLogisticsRsp, errCode int, err error)           //上传物流

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
	submitMap = map[string]SumbitHandler{
		consts.AlipayOrg: alipayPayment.CreatePayUrl,
	}

	ampSubmitMap = map[string]SumbitHandler{
		consts.AlipayOrg: alipayPayment.CreateAmpPayStr,
	}

	appSubmitMap = map[string]AppSumbitHandler{
		consts.AlipayOrg: alipayPayment.CreateAppPayStr,
	}

	configCodeMap = map[string]ConfigCodeHandler{
		consts.AlipayOrg: alipayPayment.GetConfigCode,
	}

	notifyMap = map[string]NotifyHandler{
		consts.AlipayOrg: new(alipay.Notify).Validate,
	}
}

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

func getNotifyHandler(orgCode string) NotifyHandler {
	return notifyMap[orgCode]
}

func getConfigCodeHandler(orgCode string) ConfigCodeHandler {
	return configCodeMap[orgCode]
}
