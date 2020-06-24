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
)

type Notice struct {
	Handler payment.NotificeHandler
}

func NewNotice(orgCode string) (*Notice, int, error) {
	notice := new(Notice)

	notice.Handler = payment.GetNoticeHandler(orgCode)
	if notice.Handler == nil {
		return notice, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	return notice, 0, nil
}

func (n *Notice) getConfigValue(configCodes []string, orgCode string) (payParamMap map[string]string, errCode int, err error) {
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

func (n *Notice) Notify(query, orgCode, methodCode string) (rsp validate.NotifyRsp, errCode int, err error) {
	//获取配置项code
	configCode := n.Handler.GetNotifyConfigCode()

	//读取配置项值
	configParamMap, errCode, err := n.getConfigValue(configCode, orgCode)
	if err != nil {
		return rsp, errCode, err
	}

	//异步通知处理
	notifyRsp, errCode, err := n.Handler.Notify(configParamMap, query, methodCode)
	if err != nil {
		return rsp, errCode, err
	}

	rsp.OrderId = notifyRsp.OrderId
	rsp.Status = notifyRsp.Status
	rsp.TradeNo = notifyRsp.TradeNo
	rsp.Message = notifyRsp.Message
	rsp.Rate = notifyRsp.Rate
	rsp.RmbFee = notifyRsp.RmbFee
	rsp.PaidAt = notifyRsp.PaidAt

	return rsp, 0, nil
}

func (n *Notice) Verify(query, orgCode, methodCode string) (rsp validate.VerifyRsp, errCode int, err error) {
	//获取配置项code
	configCode := n.Handler.GetVerifyConfigCode()

	//读取配置项值
	configParamMap, errCode, err := n.getConfigValue(configCode, orgCode)
	if err != nil {
		return rsp, errCode, err
	}

	//同步通知处理
	verifyRsp, errCode, err := n.Handler.Verify(configParamMap, query, methodCode)
	if err != nil {
		return rsp, errCode, err
	}
	rsp.Status = verifyRsp.Status
	rsp.OrderId = verifyRsp.OrderId

	return rsp, 0, nil
}
