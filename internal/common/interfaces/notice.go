package interfaces

import (
	"payment_demo/pkg/payment/alipay"
	"payment_demo/pkg/payment/allpay"
	"payment_demo/pkg/payment/common"
	"payment_demo/pkg/payment/consts"
	"payment_demo/pkg/payment/epayments"
	"payment_demo/pkg/payment/jd"
)

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
