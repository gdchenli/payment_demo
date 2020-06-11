package cashier

import (
	"payment_demo/internal/cashier"
	"payment_demo/internal/common/defs"
)

type PayChain func(defs.Order) (string, int, error)                               //发起支付
type NotifyChain func(string, string) (defs.NotifyRsp, int, error)                //异步通知
type CallbackChain func(string, string) (defs.CallbackRsp, int, error)            //同步通知
type TradeChain func(string, string, string, float64) (defs.TradeRsp, int, error) //交易查询
type ClosedChain func(defs.Closed) (defs.ClosedRsp, int, error)                   //关闭交易

var payMap map[string]PayChain
var notifyMap map[string]NotifyChain
var callbackMap map[string]CallbackChain
var tradeMap map[string]TradeChain
var closedMap map[string]ClosedChain

func init() {
	payMap = map[string]PayChain{
		JdOrg:        new(cashier.Jd).Pay,
		AllpayOrg:    new(cashier.Allpay).Pay,
		AlipayOrg:    new(cashier.Alipay).Pay,
		EpaymentsOrg: new(cashier.Epayments).Pay,
	}

	notifyMap = map[string]NotifyChain{
		JdOrg:        new(cashier.Jd).Notify,
		AllpayOrg:    new(cashier.Allpay).Notify,
		AlipayOrg:    new(cashier.Alipay).Notify,
		EpaymentsOrg: new(cashier.Epayments).Notify,
	}

	callbackMap = map[string]CallbackChain{
		JdOrg:        new(cashier.Jd).Callback,
		AllpayOrg:    new(cashier.Allpay).Callback,
		AlipayOrg:    new(cashier.Alipay).Callback,
		EpaymentsOrg: new(cashier.Epayments).Callback,
	}

	tradeMap = map[string]TradeChain{
		JdOrg:        new(cashier.Jd).Trade,
		AllpayOrg:    new(cashier.Allpay).Trade,
		AlipayOrg:    new(cashier.Alipay).Trade,
		EpaymentsOrg: new(cashier.Epayments).Trade,
	}

	closedMap = map[string]ClosedChain{
		JdOrg:        new(cashier.Jd).Closed,
		EpaymentsOrg: new(cashier.Epayments).Closed,
	}
}

func getPayHandler(org string) PayChain {
	return payMap[org]
}

func getNotifyHandler(org string) NotifyChain {
	return notifyMap[org]
}

func getCallbackHandler(org string) CallbackChain {
	return callbackMap[org]
}

func getTradeHandler(org string) TradeChain {
	return tradeMap[org]
}

func getClosedHandler(org string) ClosedChain {
	return closedMap[org]
}
