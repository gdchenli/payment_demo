package alipay

import (
	"payment_demo/api/response"
	"payment_demo/internal/common/request"
)

func (alipay *Alipay) CloseTrade(configParamMap map[string]string, req request.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {
	return closeTradeRsp, 0, nil
}

func (alipay *Alipay) GetCloseTradeConfigCode() []string {
	return []string{}
}
