package payment

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"payment_demo/api/request"
	"payment_demo/api/response"
	"payment_demo/internal/common/code"
	"payment_demo/pkg/config"
	"payment_demo/pkg/payment/epayments"
	"payment_demo/pkg/payment/jd"
)

type Payment struct {
	OrgHandler OrgHandler
}

func New(orgCode string) (*Payment, int, error) {
	payment := new(Payment)

	payment.OrgHandler = getOrgHandler(orgCode)
	if payment.OrgHandler == nil {
		return payment, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	return payment, 0, nil
}

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

func (payment *Payment) Pay(order request.OrderArg) (pay string, errCode int, err error) {
	//获取配置项code
	configCode := payment.OrgHandler.GetPayConfigCode()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, order.OrgCode)
	if err != nil {
		return pay, errCode, err
	}

	//处理配置映射

	//支付处理
	pay, errCode, err = payment.OrgHandler.CreatePayUrl(configParamMap, order)
	if err != nil {
		return pay, errCode, err
	}

	return pay, 0, nil
}

func (payment *Payment) PayQrCode(order request.OrderArg) (pay string, errCode int, err error) {
	//获取配置项code
	epaymentsPayment := epayments.New()
	configCode := epaymentsPayment.GetPayConfigCode()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, order.OrgCode)
	if err != nil {
		return pay, errCode, err
	}

	//处理配置code映射

	//支付处理
	pay, errCode, err = epaymentsPayment.CreateQrCode(configParamMap, order)
	if err != nil {
		return pay, errCode, err
	}

	return pay, 0, nil
}

func (payment *Payment) PayForm(order request.OrderArg) (pay string, errCode int, err error) {
	//获取配置项code
	jdPayment := jd.New()
	configCode := jdPayment.GetPayConfigCode()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, order.OrgCode)
	if err != nil {
		return pay, errCode, err
	}

	//处理配置code映射

	//支付处理
	pay, errCode, err = jdPayment.CreatePayForm(configParamMap, order)
	if err != nil {
		return pay, errCode, err
	}

	return pay, 0, nil
}

func (payment *Payment) Notify(query, orgCode, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error) {
	//获取配置项code
	configCode := payment.OrgHandler.GetNotifyConfigCode()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, orgCode)
	if err != nil {
		return notifyRsp, errCode, err
	}

	//处理配置映射

	//异步通知处理
	notifyRsp, errCode, err = payment.OrgHandler.Notify(configParamMap, query, methodCode)
	if err != nil {
		return notifyRsp, errCode, err
	}

	return notifyRsp, 0, nil
}

func (payment *Payment) Verify(query, orgCode, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) {
	//获取配置项code
	configCode := payment.OrgHandler.GetVerifyConfigCode()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, orgCode)
	if err != nil {
		return verifyRsp, errCode, err
	}

	//处理配置code映射

	//同步通知处理
	verifyRsp, errCode, err = payment.OrgHandler.Verify(configParamMap, query, methodCode)
	if err != nil {
		return verifyRsp, errCode, err
	}

	return verifyRsp, 0, nil
}

func (payment *Payment) SearchTrade(req request.SearchTradeArg) (searchTradeRsp response.SearchTradeRsp, errCode int, err error) {
	//获取配置项code
	configCode := payment.OrgHandler.GetVerifyConfigCode()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, req.OrgCode)
	if err != nil {
		return searchTradeRsp, errCode, err
	}

	//处理配置code映射

	searchTradeRsp, errCode, err = payment.OrgHandler.SearchTrade(configParamMap, req)
	if err != nil {
		return searchTradeRsp, errCode, err
	}

	return searchTradeRsp, 0, nil
}

func (payment *Payment) CloseTrade(req request.CloseTradeArg) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {
	//获取配置项code
	configCode := payment.OrgHandler.GetCloseTradeConfigCode()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, req.OrgCode)
	if err != nil {
		return closeTradeRsp, errCode, err
	}

	//处理配置code映射

	//关闭支付交易处理
	closeTradeRsp, errCode, err = payment.OrgHandler.CloseTrade(configParamMap, req)
	if err != nil {
		return closeTradeRsp, errCode, err
	}

	return closeTradeRsp, 0, nil
}

func (payment *Payment) UploadLogistics(req request.UploadLogisticsArg) (uploadLogisticsTradeRsp response.UploadLogisticsRsp, errCode int, err error) {
	//获取配置项code
	jdPayment := jd.New()
	configCode := jdPayment.GetUploadLogisticsConfigCode()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, req.OrgCode)
	if err != nil {
		return uploadLogisticsTradeRsp, errCode, err
	}

	//处理配置code映射

	//上传物流信息处理
	uploadLogisticsTradeRsp, errCode, err = jdPayment.UploadLogistics(configParamMap, req)
	if err != nil {
		return uploadLogisticsTradeRsp, errCode, err
	}

	return uploadLogisticsTradeRsp, 0, nil
}
