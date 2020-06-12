package service

import (
	"payment_demo/api/response"
	"payment_demo/api/validate"
	"payment_demo/pkg/config"
)

type Payment struct{}

func (payment *Payment) getConfigValue(configCodes []string) (payParamMap map[string]string, errCode int, err error) {
	payParamMap = make(map[string]string)
	for _, configCode := range configCodes {
		payParamMap[configCode] = config.GetInstance().GetString("alipay." + configCode)
	}

	return payParamMap, 0, nil
}

func (payment *Payment) Sumbit(order validate.Order) (pay interface{}, errCode int, err error) {
	//获取配置项code
	getConfigCodehandle := getConfigCodeHandler(order.OrgCode)
	if getConfigCodehandle == nil {
		return pay, errCode, err
	}

	configCode := getConfigCodehandle()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode)
	if err != nil {
		return pay, errCode, err
	}

	//支付处理
	submitHandle := getSubmitHandler(order.OrgCode, order.UserAgentType)
	if submitHandle == nil {
		return pay, errCode, err
	}

	pay, errCode, err = submitHandle(configParamMap, order)
	if err != nil {
		return pay, errCode, err
	}

	return pay, 0, nil
}

func (payment *Payment) Notify(query, orgCode, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error) {

	return notifyRsp, 0, nil
}

func (payment *Payment) Verify(query, orgCode, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) {

	return verifyRsp, 0, nil
}

func (payment *Payment) SearchTrade(req validate.SearchTradeReq, orgCode, methodCode string) (searchTradeRsp response.SearchTradeRsp, errCode int, err error) {

	return searchTradeRsp, 0, nil
}

func (payment *Payment) CloseTrade(req validate.CloseTradeReq, orgCode, methodCode string) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {

	return closeTradeRsp, 0, nil
}

func (payment *Payment) UploadLogistics(req validate.UploadLogisticsReq, orgCode, methodCode string) (uploadLogisticsTradeRsp response.UploadLogisticsRsp, errCode int, err error) {

	return uploadLogisticsTradeRsp, 0, nil
}
