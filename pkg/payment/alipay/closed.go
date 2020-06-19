package alipay

import (
	"payment_demo/api/payment/request"
	"payment_demo/api/payment/response"
)

func (alipay *Alipay) CloseTrade(configParamMap map[string]string, req request.CloseTradeArg) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {
	return closeTradeRsp, 0, nil
}

func (alipay *Alipay) GetCloseTradeConfigCode() []string {
	return []string{}
}
