package payment

import (
	"payment_demo/api/response"
	"payment_demo/internal/common/request"
	"payment_demo/pkg/payment/alipay"
	"payment_demo/pkg/payment/allpay"
	"payment_demo/pkg/payment/consts"
	"payment_demo/pkg/payment/jd"
)

type OrgHandler interface {
	//发起支付
	CreatePayUrl(configParamMap map[string]string, order request.Order) (url string, errCode int, err error) //pc、h5、支付宝小程序
	/*WmpSumbit(configParamMap map[string]string, order request.Order) (wmRsp response.WmpRsp, errCode int, err error)  //微信小程序
	AppSumbit(configParamMap map[string]string, order request.Order) (appRsp response.AppRsp, errCode int, err error) //App*/

	//支付通知
	Notify(configParamMap map[string]string, query, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error) //异步通知
	Verify(configParamMap map[string]string, query, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) //同步通知

	//交易信息
	SearchTrade(configParamMap map[string]string, req request.SearchTradeReq) (searchTradeRsp response.SearchTradeRsp, errCode int, err error) //交易查询
	CloseTrade(configParamMap map[string]string, req request.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error)     //关闭交易
	//UploadLogistics(configParamMap map[string]string, req request.UploadLogisticsReq) (uploadLogisticsRsp response.UploadLogisticsRsp, errCode int, err error) //上传物流

	//配置
	GetPayConfigCode() []string
	GetNotifyConfigCode() []string
	GetVerifyConfigCode() []string
	GetSearchTradeConfigCode() []string
	GetCloseTradeConfigCode() []string
}

func getOrgHandler(orgCode string) OrgHandler {
	switch orgCode {
	case consts.AlipayOrg:
		return alipay.New()
	case consts.EpaymentsOrg:
		return allpay.New()
	case consts.AllpayOrg:
		return allpay.New()
	case consts.JdOrg:
		return jd.New()
	default:
		return nil
	}
}
