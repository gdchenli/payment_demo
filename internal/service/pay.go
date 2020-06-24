package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"payment_demo/api/validate"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/interfaces"
	"payment_demo/pkg/config"
	"payment_demo/pkg/payment/common"
	"payment_demo/pkg/payment/epayments"
	"payment_demo/pkg/payment/jd"
)

type Payment struct {
	Handler interfaces.PayHandler
}

func NewPay(orgCode string) (*Payment, int, error) {
	payment := new(Payment)

	payment.Handler = interfaces.GetPayHandler(orgCode)
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

func (payment *Payment) Pay(order validate.OrderArg) (pay string, errCode int, err error) {
	//获取配置项code
	configCode := payment.Handler.GetPayConfigCode()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, order.OrgCode)
	if err != nil {
		return pay, errCode, err
	}

	//支付处理
	pay, errCode, err = payment.Handler.CreatePayUrl(configParamMap, payment.getOrderArg(order))
	if err != nil {
		return pay, errCode, err
	}

	return pay, 0, nil
}

func (payment *Payment) getOrderArg(arg validate.OrderArg) common.OrderArg {
	return common.OrderArg{
		OrderId:       arg.OrderId,
		Currency:      arg.Currency,
		MethodCode:    arg.MethodCode,
		OrgCode:       arg.OrgCode,
		UserId:        arg.UserId,
		UserAgentType: arg.UserAgentType}
}

func (payment *Payment) PayQrCode(order validate.OrderArg) (pay string, errCode int, err error) {
	//获取配置项code
	epaymentsPayment := epayments.New()
	configCode := epaymentsPayment.GetPayConfigCode()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, order.OrgCode)
	if err != nil {
		return pay, errCode, err
	}

	//支付处理
	pay, errCode, err = epaymentsPayment.CreateQrCode(configParamMap, payment.getOrderArg(order))
	if err != nil {
		return pay, errCode, err
	}

	return pay, 0, nil
}

func (payment *Payment) PayForm(order validate.OrderArg) (pay string, errCode int, err error) {
	//获取配置项code
	jdPayment := jd.New()
	configCode := jdPayment.GetPayConfigCode()

	//读取配置项值
	configParamMap, errCode, err := payment.getConfigValue(configCode, order.OrgCode)
	if err != nil {
		return pay, errCode, err
	}

	//支付处理
	pay, errCode, err = jdPayment.CreatePayForm(configParamMap, payment.getOrderArg(order))
	if err != nil {
		return pay, errCode, err
	}

	return pay, 0, nil
}
