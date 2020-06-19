package allpay

import (
	"payment_demo/api/payment/request"
	"payment_demo/api/payment/response"
)

func (allpay *Allpay) CloseTrade(configParamMap map[string]string, req request.CloseTradeArg) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {
	return closeTradeRsp, 0, nil
}

func (allpay *Allpay) GetCloseTradeConfigCode() []string {
	return []string{}
}
