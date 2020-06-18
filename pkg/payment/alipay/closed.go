package alipay

import (
	"payment_demo/api/request"
	"payment_demo/api/response"
)

func (alipay *Alipay) CloseTrade(configParamMap map[string]string, req request.CloseTradeArg) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {
	return closeTradeRsp, 0, nil
}

func (alipay *Alipay) GetCloseTradeConfigCode() []string {
	return []string{}
}
