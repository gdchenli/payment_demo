package alipay

import (
	request2 "payment_demo/api/trade/request"
	response2 "payment_demo/api/trade/response"
)

func (alipay *Alipay) CloseTrade(configParamMap map[string]string, req request2.CloseTradeArg) (closeTradeRsp response2.CloseTradeRsp, errCode int, err error) {
	return closeTradeRsp, 0, nil
}

func (alipay *Alipay) GetCloseTradeConfigCode() []string {
	return []string{}
}
