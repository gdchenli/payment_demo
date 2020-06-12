package controller

import (
	"payment_demo/api/response"
	"payment_demo/api/validate"
	"payment_demo/internal/service/payment"
)

const (
	NotSupportPaymentOrgCode = "10101"
	NotSupportPaymentOrgMsg  = "不支持该支付机构"
)

const (
	JdOrg        = "jd"
	AllpayOrg    = "allpay"
	EpaymentsOrg = "epayments"
	AlipayOrg    = "alipay"
)

type PayHandler func(arg validate.Order) (form string, errCode int, err error)                                                                       //发起支付
type NotifyHandler func(query, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error)                                             //异步通知
type VerifyHandler func(query, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error)                                             //同步通知
type SearchTradeHandler func(orderId, methodCode, curreny string, totalFee float64) (searchTradeRsp response.SearchTradeRsp, errCode int, err error) //交易查询
type CloseTradeHandler func(arg validate.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error)                               //关闭交易
type UploadLogisticsHandler func(arg validate.UploadLogisticsReq) (uploadLogisticsRsp response.UploadLogisticsRsp, errCode int, err error)           //上传物流

var payMap map[string]PayHandler
var notifyMap map[string]NotifyHandler
var verifyMap map[string]VerifyHandler
var searchTradeMap map[string]SearchTradeHandler
var closeTradeMap map[string]CloseTradeHandler
var uploadLogisticsMap map[string]UploadLogisticsHandler

func init() {
	payMap = map[string]PayHandler{
		JdOrg:        new(payment.Jd).Pay,
		AllpayOrg:    new(payment.Allpay).Pay,
		AlipayOrg:    new(payment.Alipay).Pay,
		EpaymentsOrg: new(payment.Epayments).Pay,
	}

	notifyMap = map[string]NotifyHandler{
		JdOrg:        new(payment.Jd).Notify,
		AllpayOrg:    new(payment.Allpay).Notify,
		AlipayOrg:    new(payment.Alipay).Notify,
		EpaymentsOrg: new(payment.Epayments).Notify,
	}

	verifyMap = map[string]VerifyHandler{
		JdOrg:        new(payment.Jd).Verify,
		AllpayOrg:    new(payment.Allpay).Verify,
		AlipayOrg:    new(payment.Alipay).Verify,
		EpaymentsOrg: new(payment.Epayments).Verify,
	}

	searchTradeMap = map[string]SearchTradeHandler{
		JdOrg:        new(payment.Jd).SearchTrade,
		AllpayOrg:    new(payment.Allpay).SearchTrade,
		AlipayOrg:    new(payment.Alipay).SearchTrade,
		EpaymentsOrg: new(payment.Epayments).SearchTrade,
	}

	closeTradeMap = map[string]CloseTradeHandler{
		JdOrg:        new(payment.Jd).CloseTrade,
		EpaymentsOrg: new(payment.Epayments).CloseTrade,
	}

	uploadLogisticsMap = map[string]UploadLogisticsHandler{
		JdOrg: new(payment.Jd).UploadLogistics,
	}
}

func GetPayHandler(org string) PayHandler {
	return payMap[org]
}

func GetNotifyHandler(org string) NotifyHandler {
	return notifyMap[org]
}

func GetVerifyHandler(org string) VerifyHandler {
	return verifyMap[org]
}

func GetSearchTradeHandler(org string) SearchTradeHandler {
	return searchTradeMap[org]
}

func GetCloseTradeHandler(org string) CloseTradeHandler {
	return closeTradeMap[org]
}

func GetUploadLogisticsHandler(org string) UploadLogisticsHandler {
	return uploadLogisticsMap[org]
}
