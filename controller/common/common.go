package common

import (
	"payment_demo/internal/common/defs"
	"payment_demo/internal/service"
)

type PayChain func(defs.Order) (string, int, error)                                           //发起支付
type NotifyChain func(string, string) (defs.NotifyRsp, int, error)                            //异步通知
type VerifyChain func(string, string) (defs.VerifyRsp, int, error)                            //同步通知
type SearchTradeChain func(string, string, string, float64) (defs.SearchTradeRsp, int, error) //交易查询
type CloseTradeChain func(defs.CloseTradeReq) (defs.CloseTradeRsp, int, error)                //关闭交易
type UploadLogisticsChain func(defs.UploadLogisticsReq) (defs.UploadLogisticsRsp, int, error) //上传物流

var payMap map[string]PayChain
var notifyMap map[string]NotifyChain
var verifyMap map[string]VerifyChain
var searchTradeMap map[string]SearchTradeChain
var closeTradeMap map[string]CloseTradeChain
var uploadLogisticsMap map[string]UploadLogisticsChain

func init() {
	payMap = map[string]PayChain{
		JdOrg:        new(service.Jd).Pay,
		AllpayOrg:    new(service.Allpay).Pay,
		AlipayOrg:    new(service.Alipay).Pay,
		EpaymentsOrg: new(service.Epayments).Pay,
	}

	notifyMap = map[string]NotifyChain{
		JdOrg:        new(service.Jd).Notify,
		AllpayOrg:    new(service.Allpay).Notify,
		AlipayOrg:    new(service.Alipay).Notify,
		EpaymentsOrg: new(service.Epayments).Notify,
	}

	verifyMap = map[string]VerifyChain{
		JdOrg:        new(service.Jd).Verify,
		AllpayOrg:    new(service.Allpay).Verify,
		AlipayOrg:    new(service.Alipay).Verify,
		EpaymentsOrg: new(service.Epayments).Verify,
	}

	searchTradeMap = map[string]SearchTradeChain{
		JdOrg:        new(service.Jd).SearchTrade,
		AllpayOrg:    new(service.Allpay).SearchTrade,
		AlipayOrg:    new(service.Alipay).SearchTrade,
		EpaymentsOrg: new(service.Epayments).SearchTrade,
	}

	closeTradeMap = map[string]CloseTradeChain{
		JdOrg:        new(service.Jd).CloseTrade,
		EpaymentsOrg: new(service.Epayments).CloseTrade,
	}

	uploadLogisticsMap = map[string]UploadLogisticsChain{
		JdOrg: new(service.Jd).UploadLogistics,
	}
}

func GetPayHandler(org string) PayChain {
	return payMap[org]
}

func GetNotifyHandler(org string) NotifyChain {
	return notifyMap[org]
}

func GetVerifyHandler(org string) VerifyChain {
	return verifyMap[org]
}

func GetSearchTradeHandler(org string) SearchTradeChain {
	return searchTradeMap[org]
}

func GetCloseTradeHandler(org string) CloseTradeChain {
	return closeTradeMap[org]
}

func GetUploadLogisticsHandler(org string) UploadLogisticsChain {
	return uploadLogisticsMap[org]
}
