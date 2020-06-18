package payment

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"payment_demo/api/response"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/request"
	"payment_demo/pkg/config"
)

type Payment struct{}

func (payment *Payment) getConfigValue(configCodes []string, orgCode string) (payParamMap map[string]string, errCode int, err error) {
	payParamMap = make(map[string]string)
	for _, configCode := range configCodes {
		payParamMap[configCode] = config.GetInstance().GetString(orgCode + "." + configCode)
		if configCode == "private_key" || configCode == "public_key" {
			keyPath := path.Join(config.GetInstance().GetString("app_path"), payParamMap[configCode])
			fmt.Println("keyPath", keyPath)
			keyFile, err := os.Open(keyPath)
			if err != nil {
				fmt.Println("keyPath err", err)
				payParamMap[configCode] = ""
			}

			keyBytes, err := ioutil.ReadAll(keyFile)
			if err != nil {
				fmt.Println("keyBytes err", err)
				payParamMap[configCode] = ""
			}

			payParamMap[configCode] = string(keyBytes)
		}
		if payParamMap[configCode] == "" {
			fmt.Println("configCode", configCode)
			return payParamMap, code.ConfigValueErrCode, errors.New(code.ConfigValueErrMessage)
		}
	}

	return payParamMap, 0, nil
}

func (payment *Payment) Pay(order request.Order, istransfer bool) (pay string, errCode int, err error) {
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
	submitHandle := getSubmitHandler(order.OrgCode, istransfer)
	if submitHandle == nil {
		fmt.Println("submitHandle")
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
	getConfigCodehandle := getConfigCodeHandler(orgCode + ".verify")
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

func (payment *Payment) SearchTrade(req request.SearchTradeReq) (searchTradeRsp response.SearchTradeRsp, errCode int, err error) {
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

func (payment *Payment) CloseTrade(req request.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {
	//获取配置项code
	getConfigCodehandle := getConfigCodeHandler(req.OrgCode + ".close")
	if getConfigCodehandle == nil {
		return closeTradeRsp, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	configCode := getConfigCodehandle()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, req.OrgCode)
	if err != nil {
		return closeTradeRsp, errCode, err
	}

	//关闭支付交易处理
	closedTradeHandle := getCloseTradeHandler(req.OrgCode)
	if closedTradeHandle == nil {
		return closeTradeRsp, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	closeTradeRsp, errCode, err = closedTradeHandle(configParamMap, req)
	if err != nil {
		return closeTradeRsp, errCode, err
	}

	return closeTradeRsp, 0, nil
}

func (payment *Payment) UploadLogistics(req request.UploadLogisticsReq) (uploadLogisticsTradeRsp response.UploadLogisticsRsp, errCode int, err error) {
	//获取配置项code
	getConfigCodehandle := getConfigCodeHandler(req.OrgCode + ".logistics")
	if getConfigCodehandle == nil {
		return uploadLogisticsTradeRsp, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	configCode := getConfigCodehandle()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, req.OrgCode)
	if err != nil {
		return uploadLogisticsTradeRsp, errCode, err
	}

	//上传物流信息处理
	uploadLogisticsTradeHandle := getUploadLogisticsHandler(req.OrgCode)
	if uploadLogisticsTradeHandle == nil {
		return uploadLogisticsTradeRsp, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	uploadLogisticsTradeRsp, errCode, err = uploadLogisticsTradeHandle(configParamMap, req)
	if err != nil {
		return uploadLogisticsTradeRsp, errCode, err
	}

	return uploadLogisticsTradeRsp, 0, nil
}