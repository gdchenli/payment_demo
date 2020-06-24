package payment

import (
	"payment_demo/pkg/payment/alipay"
	"payment_demo/pkg/payment/allpay"
	"payment_demo/pkg/payment/common"
	"payment_demo/pkg/payment/consts"
	"payment_demo/pkg/payment/epayments"
	"payment_demo/pkg/payment/jd"
)

type PayHandler interface {
	//发起支付
	CreatePayUrl(configParamMap map[string]string, order common.OrderArg) (url string, errCode int, err error) //pc、h5、支付宝小程序
	/*WmpSumbit(configParamMap map[string]string, order common.Order) (wmRsp common.WmpRsp, errCode int, err error)  //微信小程序
	AppSumbit(configParamMap map[string]string, order common.Order) (appRsp common.AppRsp, errCode int, err error) //App*/

	//配置
	GetPayConfigCode() []string
}

func GetPayHandler(orgCode string) PayHandler {
	switch orgCode {
	case consts.AlipayOrg:
		return alipay.New()
	case consts.EpaymentsOrg:
		return epayments.New()
	case consts.AllpayOrg:
		return allpay.New()
	case consts.JdOrg:
		return jd.New()
	default:
		return nil
	}
}

type NotificeHandler interface {
	//支付通知
	Notify(configParamMap map[string]string, query, methodCode string) (notifyRsp common.NotifyRsp, errCode int, err error) //异步通知
	Verify(configParamMap map[string]string, query, methodCode string) (verifyRsp common.VerifyRsp, errCode int, err error) //同步通知

	//配置
	GetNotifyConfigCode() []string
	GetVerifyConfigCode() []string
}

func GetNoticeHandler(orgCode string) NotificeHandler {
	switch orgCode {
	case consts.AlipayOrg:
		return alipay.New()
	case consts.EpaymentsOrg:
		return epayments.New()
	case consts.AllpayOrg:
		return allpay.New()
	case consts.JdOrg:
		return jd.New()
	default:
		return nil
	}
}

type TradeHandler interface {

	//交易信息
	SearchTrade(configParamMap map[string]string, req common.SearchTradeArg) (searchTradeRsp common.SearchTradeRsp, errCode int, err error) //交易查询
	CloseTrade(configParamMap map[string]string, req common.CloseTradeArg) (closeTradeRsp common.CloseTradeRsp, errCode int, err error)     //关闭交易
	//UploadLogistics(configParamMap map[string]string, req request.UploadLogisticsReq) (uploadLogisticsRsp response.UploadLogisticsRsp, errCode int, err error) //上传物流

	//配置
	GetSearchTradeConfigCode() []string
	GetCloseTradeConfigCode() []string
}

func GetTradeHandler(orgCode string) TradeHandler {
	switch orgCode {
	case consts.AlipayOrg:
		return alipay.New()
	case consts.EpaymentsOrg:
		return epayments.New()
	case consts.AllpayOrg:
		return allpay.New()
	case consts.JdOrg:
		return jd.New()
	default:
		return nil
	}
}
