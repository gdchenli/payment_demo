package common

import (
	"payment_demo/internal/common/defs"
	"payment_demo/internal/service/payment"
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
		JdOrg:        new(payment.Jd).Pay,
		AllpayOrg:    new(payment.Allpay).Pay,
		AlipayOrg:    new(payment.Alipay).Pay,
		EpaymentsOrg: new(payment.Epayments).Pay,
	}

	notifyMap = map[string]NotifyChain{
		JdOrg:        new(payment.Jd).Notify,
		AllpayOrg:    new(payment.Allpay).Notify,
		AlipayOrg:    new(payment.Alipay).Notify,
		EpaymentsOrg: new(payment.Epayments).Notify,
	}

	verifyMap = map[string]VerifyChain{
		JdOrg:        new(payment.Jd).Verify,
		AllpayOrg:    new(payment.Allpay).Verify,
		AlipayOrg:    new(payment.Alipay).Verify,
		EpaymentsOrg: new(payment.Epayments).Verify,
	}

	searchTradeMap = map[string]SearchTradeChain{
		JdOrg:        new(payment.Jd).SearchTrade,
		AllpayOrg:    new(payment.Allpay).SearchTrade,
		AlipayOrg:    new(payment.Alipay).SearchTrade,
		EpaymentsOrg: new(payment.Epayments).SearchTrade,
	}

	closeTradeMap = map[string]CloseTradeChain{
		JdOrg:        new(payment.Jd).CloseTrade,
		EpaymentsOrg: new(payment.Epayments).CloseTrade,
	}

	uploadLogisticsMap = map[string]UploadLogisticsChain{
		JdOrg: new(payment.Jd).UploadLogistics,
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
