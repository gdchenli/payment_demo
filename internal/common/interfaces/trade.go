package interfaces

import (
	"payment_demo/pkg/payment/alipay"
	"payment_demo/pkg/payment/allpay"
	"payment_demo/pkg/payment/common"
	"payment_demo/pkg/payment/consts"
	"payment_demo/pkg/payment/epayments"
	"payment_demo/pkg/payment/jd"
)

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
