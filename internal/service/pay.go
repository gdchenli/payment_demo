package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"payment_demo/api/validate"
	"payment_demo/internal/common/code"
	"payment_demo/pkg/config"
	"payment_demo/pkg/payment"
	"payment_demo/pkg/payment/common"
	"payment_demo/pkg/payment/epayments"
	"payment_demo/pkg/payment/jd"
)

type Pay struct {
	Handler payment.PayHandler
}

func NewPay(orgCode string) (*Pay, int, error) {
	pay := new(Pay)

	pay.Handler = payment.GetPayHandler(orgCode)
	if pay.Handler == nil {
		return pay, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	return pay, 0, nil
}

func (pay *Pay) getConfigValue(configCodes []string, orgCode string) (payParamMap map[string]string, errCode int, err error) {
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

func (pay *Pay) Pay(order validate.OrderArg) (payStr string, errCode int, err error) {
	//获取配置项code
	configCode := pay.Handler.GetPayConfigCode()

	//读取配置项值
	configParamMap, errCode, err := pay.getConfigValue(configCode, order.OrgCode)
	if err != nil {
		return payStr, errCode, err
	}

	//支付处理
	payStr, errCode, err = pay.Handler.CreatePayUrl(configParamMap, pay.getOrderArg(order))
	if err != nil {
		return payStr, errCode, err
	}

	return payStr, 0, nil
}

func (pay *Pay) getOrderArg(arg validate.OrderArg) common.OrderArg {
	return common.OrderArg{
		OrderId:       arg.OrderId,
		Currency:      arg.Currency,
		MethodCode:    arg.MethodCode,
		OrgCode:       arg.OrgCode,
		UserId:        arg.UserId,
		UserAgentType: arg.UserAgentType}
}

func (pay *Pay) PayQrCode(order validate.OrderArg) (payStr string, errCode int, err error) {
	//获取配置项code
	epaymentsPayment := epayments.New()
	configCode := epaymentsPayment.GetPayConfigCode()

	//读取配置项值
	configParamMap, errCode, err := pay.getConfigValue(configCode, order.OrgCode)
	if err != nil {
		return payStr, errCode, err
	}

	//支付处理
	payStr, errCode, err = epaymentsPayment.CreateQrCode(configParamMap, pay.getOrderArg(order))
	if err != nil {
		return payStr, errCode, err
	}

	return payStr, 0, nil
}

func (pay *Pay) PayForm(order validate.OrderArg) (payStr string, errCode int, err error) {
	//获取配置项code
	jdPayment := jd.New()
	configCode := jdPayment.GetPayConfigCode()

	//读取配置项值
	configParamMap, errCode, err := pay.getConfigValue(configCode, order.OrgCode)
	if err != nil {
		return payStr, errCode, err
	}

	//支付处理
	payStr, errCode, err = jdPayment.CreatePayForm(configParamMap, pay.getOrderArg(order))
	if err != nil {
		return payStr, errCode, err
	}

	return payStr, 0, nil
}
