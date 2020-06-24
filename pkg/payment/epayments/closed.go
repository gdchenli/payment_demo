package epayments

import (
	"payment_demo/pkg/payment/common"
)

func (epayments *Epayments) CloseTrade(configParamMap map[string]string, req common.CloseTradeArg) (closeTradeRsp common.CloseTradeRsp, errCode int, err error) {
	return closeTradeRsp, 0, nil
}

func (epayments *Epayments) GetCloseTradeConfigCode() []string {
	return []string{}
}
