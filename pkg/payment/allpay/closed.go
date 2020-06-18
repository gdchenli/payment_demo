package allpay

import (
	"payment_demo/api/response"
	"payment_demo/internal/common/request"
)

func (allpay *Allpay) CloseTrade(configParamMap map[string]string, req request.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {
	return closeTradeRsp, 0, nil
}

func (allpay *Allpay) GetCloseTradeConfigCode() []string {
	return []string{}
}
