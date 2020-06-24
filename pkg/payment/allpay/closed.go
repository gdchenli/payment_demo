package allpay

import (
	"payment_demo/pkg/payment/common"
)

func (allpay *Allpay) CloseTrade(configParamMap map[string]string, req common.CloseTradeArg) (closeTradeRsp common.CloseTradeRsp, errCode int, err error) {
	return closeTradeRsp, 0, nil
}

func (allpay *Allpay) GetCloseTradeConfigCode() []string {
	return []string{}
}
