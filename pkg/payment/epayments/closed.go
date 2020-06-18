package epayments

import (
	"payment_demo/api/response"
	"payment_demo/internal/common/request"
)

func (epayments *Epayments) CloseTrade(configParamMap map[string]string, req request.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {
	return closeTradeRsp, 0, nil
}

func (epayments *Epayments) GetCloseTradeConfigCode() []string {
	return []string{}
}
