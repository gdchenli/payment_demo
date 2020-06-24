package alipay

import (
	"payment_demo/pkg/payment/common"
)

func (alipay *Alipay) CloseTrade(configParamMap map[string]string, req common.CloseTradeArg) (closeTradeRsp common.CloseTradeRsp, errCode int, err error) {
	return closeTradeRsp, 0, nil
}

func (alipay *Alipay) GetCloseTradeConfigCode() []string {
	return []string{}
}
