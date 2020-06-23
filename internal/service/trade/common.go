package trade

import (
	request2 "payment_demo/api/trade/request"
	response2 "payment_demo/api/trade/response"
	"payment_demo/pkg/payment/alipay"
	"payment_demo/pkg/payment/allpay"
	"payment_demo/pkg/payment/consts"
	"payment_demo/pkg/payment/epayments"
	"payment_demo/pkg/payment/jd"
)

type Handler interface {

	//交易信息
	SearchTrade(configParamMap map[string]string, req request2.SearchTradeArg) (searchTradeRsp response2.SearchTradeRsp, errCode int, err error) //交易查询
	CloseTrade(configParamMap map[string]string, req request2.CloseTradeArg) (closeTradeRsp response2.CloseTradeRsp, errCode int, err error)     //关闭交易
	//UploadLogistics(configParamMap map[string]string, req request.UploadLogisticsReq) (uploadLogisticsRsp response.UploadLogisticsRsp, errCode int, err error) //上传物流

	//配置
	GetSearchTradeConfigCode() []string
	GetCloseTradeConfigCode() []string
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
