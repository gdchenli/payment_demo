package method

import "payment_demo/internal/common/defs"

type Epayments struct{}

func (e *Epayments) OrderSubmit(order defs.Order) (string, int, error) {
	panic("implement me")
}

func (e *Epayments) Notify(query, methodCode string) (defs.NotifyRsp, int, error) {
	panic("implement me")
}

func (e *Epayments) Callback(query, methodCode string) (defs.CallbackRsp, int, error) {
	panic("implement me")
}

func (e *Epayments) Trade(orderId, methodCode string) (defs.TradeRsp, int, error) {
	panic("implement me")
}
