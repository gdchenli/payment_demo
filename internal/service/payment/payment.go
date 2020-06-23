package payment

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"payment_demo/api/payment/request"
	"payment_demo/internal/common/code"
	"payment_demo/pkg/config"
	"payment_demo/pkg/payment/epayments"
	"payment_demo/pkg/payment/jd"
)

type Payment struct {
	Handler Handler
}

func New(orgCode string) (*Payment, int, error) {
	payment := new(Payment)

	payment.Handler = getHandler(orgCode)
	if payment.Handler == nil {
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
	configCode := payment.Handler.GetPayConfigCode()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, order.OrgCode)
	if err != nil {
		return pay, errCode, err
	}

	//处理配置映射

	//支付处理
	pay, errCode, err = payment.Handler.CreatePayUrl(configParamMap, order)
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
