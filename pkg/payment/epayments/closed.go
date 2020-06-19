package epayments

import (
	"payment_demo/api/payment/request"
	"payment_demo/api/payment/response"
)

func (epayments *Epayments) CloseTrade(configParamMap map[string]string, req request.CloseTradeArg) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {
	return closeTradeRsp, 0, nil
}

func (epayments *Epayments) GetCloseTradeConfigCode() []string {
	return []string{}
}
