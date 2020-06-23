package allpay

import (
	request2 "payment_demo/api/trade/request"
	response2 "payment_demo/api/trade/response"
)

func (allpay *Allpay) CloseTrade(configParamMap map[string]string, req request2.CloseTradeArg) (closeTradeRsp response2.CloseTradeRsp, errCode int, err error) {
	return closeTradeRsp, 0, nil
}

func (allpay *Allpay) GetCloseTradeConfigCode() []string {
	return []string{}
}
