package payment

import (
	"payment_demo/api/payment/request"
	"payment_demo/pkg/payment/alipay"
	"payment_demo/pkg/payment/allpay"
	"payment_demo/pkg/payment/consts"
	"payment_demo/pkg/payment/epayments"
	"payment_demo/pkg/payment/jd"
)

type Handler interface {
	//发起支付
	CreatePayUrl(configParamMap map[string]string, order request.OrderArg) (url string, errCode int, err error) //pc、h5、支付宝小程序
	/*WmpSumbit(configParamMap map[string]string, order request.Order) (wmRsp response.WmpRsp, errCode int, err error)  //微信小程序
	AppSumbit(configParamMap map[string]string, order request.Order) (appRsp response.AppRsp, errCode int, err error) //App*/

	//配置
	GetPayConfigCode() []string
}

func getHandler(orgCode string) Handler {
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
