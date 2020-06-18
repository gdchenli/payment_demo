package allpay

import (
	"payment_demo/api/request"
	"payment_demo/api/response"
)

func (allpay *Allpay) CloseTrade(configParamMap map[string]string, req request.CloseTradeArg) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {
	return closeTradeRsp, 0, nil
}

func (allpay *Allpay) GetCloseTradeConfigCode() []string {
	return []string{}
}
