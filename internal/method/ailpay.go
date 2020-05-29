package method

import "payment_demo/internal/common/defs"

type Alipay struct{}

func (alipay *Alipay) OrderSubmit(arg defs.Order) (form string, errCode int, err error) {
	return form, 0, nil
}

func (alipay *Alipay) Notify(query, methodCode string) (notifyRsp defs.NotifyRsp, errCode int, err error) {
	return notifyRsp, 0, nil
}

func (alipay *Alipay) Callback(query, methodCode string) (callbackRsp defs.CallbackRsp, errCode int, err error) {
	return callbackRsp, 0, nil
}

func (alipay *Alipay) Trade(orderId, methodCode string) (tradeRsp defs.TradeRsp, errCode int, err error) {
	return tradeRsp, 0, nil
}
