package service

import (
	"errors"
	"fmt"
	"payment_demo/api/response"
	"payment_demo/api/validate"
	"payment_demo/internal/common/code"
	"payment_demo/pkg/config"
)

type Payment struct{}

func (payment *Payment) getConfigValue(configCodes []string, orgCode string) (payParamMap map[string]string, errCode int, err error) {
	payParamMap = make(map[string]string)
	for _, configCode := range configCodes {
		payParamMap[configCode] = config.GetInstance().GetString(orgCode + "." + configCode)
		if payParamMap[configCode] == "" {
			fmt.Println("configCode", configCode)
			return payParamMap, code.ConfigValueErrCode, errors.New(code.ConfigValueErrMessage)
		}
	}

	return payParamMap, 0, nil
}

func (payment *Payment) Sumbit(order validate.Order) (pay string, errCode int, err error) {
	//获取配置项code
	getConfigCodehandle := getConfigCodeHandler(order.OrgCode + ".payment")
	if getConfigCodehandle == nil {
		return pay, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	configCode := getConfigCodehandle()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, order.OrgCode)
	if err != nil {
		return pay, errCode, err
	}

	//支付处理
	submitHandle := getSubmitHandler(order.OrgCode, order.UserAgentType)
	if submitHandle == nil {
		return pay, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	pay, errCode, err = submitHandle(configParamMap, order)
	if err != nil {
		return pay, errCode, err
	}

	return pay, 0, nil
}

func (payment *Payment) Notify(query, orgCode, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error) {
	//获取配置项code
	getConfigCodehandle := getConfigCodeHandler(orgCode + ".notify")
	if getConfigCodehandle == nil {
		return notifyRsp, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	configCode := getConfigCodehandle()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, orgCode)
	if err != nil {
		return notifyRsp, errCode, err
	}

	//异步通知处理
	notifyHandle := getNotifyHandler(orgCode)
	if notifyHandle == nil {
		return notifyRsp, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	notifyRsp, errCode, err = notifyHandle(configParamMap, query, methodCode)
	if err != nil {
		return notifyRsp, errCode, err
	}

	return notifyRsp, 0, nil
}

func (payment *Payment) Verify(query, orgCode, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) {
	//获取配置项code
	getConfigCodehandle := getConfigCodeHandler(orgCode + ".notify")
	if getConfigCodehandle == nil {
		return verifyRsp, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	configCode := getConfigCodehandle()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, orgCode)
	if err != nil {
		return verifyRsp, errCode, err
	}

	//同步通知处理
	notifyHandle := getVerifyHandler(orgCode)
	if notifyHandle == nil {
		return verifyRsp, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	verifyRsp, errCode, err = notifyHandle(configParamMap, query, methodCode)
	if err != nil {
		return verifyRsp, errCode, err
	}

	return verifyRsp, 0, nil
}

func (payment *Payment) SearchTrade(req validate.SearchTradeReq) (searchTradeRsp response.SearchTradeRsp, errCode int, err error) {
	//获取配置项code
	getConfigCodehandle := getConfigCodeHandler(req.OrgCode + ".trade")
	if getConfigCodehandle == nil {
		return searchTradeRsp, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	configCode := getConfigCodehandle()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, req.OrgCode)
	if err != nil {
		return searchTradeRsp, errCode, err
	}

	//查询支付交易处理
	searchTradeHandle := getSeachTradeHandler(req.OrgCode)
	if searchTradeHandle == nil {
		return searchTradeRsp, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	searchTradeRsp, errCode, err = searchTradeHandle(configParamMap, req)
	if err != nil {
		return searchTradeRsp, errCode, err
	}

	return searchTradeRsp, 0, nil
}

func (payment *Payment) CloseTrade(req validate.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {

	return closeTradeRsp, 0, nil
}

func (payment *Payment) UploadLogistics(req validate.UploadLogisticsReq) (uploadLogisticsTradeRsp response.UploadLogisticsRsp, errCode int, err error) {

	return uploadLogisticsTradeRsp, 0, nil
}
