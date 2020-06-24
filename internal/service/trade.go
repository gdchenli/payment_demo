package service

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"payment_demo/app/request"
	"payment_demo/app/response"
	"payment_demo/internal/common/code"
	"payment_demo/pkg/config"
	"payment_demo/pkg/payment"
	"payment_demo/pkg/payment/common"
)

type Trade struct {
	Handler payment.TradeHandler
}

func NewTrade(orgCode string) (*Trade, int, error) {
	trade := new(Trade)

	trade.Handler = payment.GetTradeHandler(orgCode)
	if trade.Handler == nil {
		return trade, code.NotSupportOrgErrCode, errors.New(code.NotSupportOrgErrMessage)
	}
	return trade, 0, nil
}

func (t *Trade) getConfigValue(configCodes []string, orgCode string) (payParamMap map[string]string, errCode int, err error) {
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
func (t *Trade) SearchTrade(req request.SearchTradeArg) (rsp response.SearchTradeRsp, errCode int, err error) {
	//获取配置项code
	configCode := t.Handler.GetSearchTradeConfigCode()

	//读取配置项值
	configParamMap, errCode, err := t.getConfigValue(configCode, req.OrgCode)
	if err != nil {
		return rsp, errCode, err
	}

	searchTradeArg := common.SearchTradeArg{
		OrderId:    req.OrderId,
		MethodCode: req.MethodCode,
		OrgCode:    req.OrgCode,
		Currency:   req.Currency,
		TotalFee:   req.TotalFee,
	}
	searchTradeRsp, errCode, err := t.Handler.SearchTrade(configParamMap, searchTradeArg)
	if err != nil {
		return rsp, errCode, err
	}
	rsp.OrderId = searchTradeRsp.OrderId
	rsp.PaidAt = searchTradeRsp.PaidAt
	rsp.RmbFee = searchTradeRsp.RmbFee
	rsp.Rate = searchTradeRsp.Rate
	rsp.Status = searchTradeRsp.Status
	rsp.TradeNo = searchTradeRsp.TradeNo

	return rsp, 0, nil
}

func (t *Trade) CloseTrade(req request.CloseTradeArg) (rsp response.CloseTradeRsp, errCode int, err error) {
	//获取配置项code
	configCode := t.Handler.GetCloseTradeConfigCode()

	//读取配置项值
	configParamMap, errCode, err := t.getConfigValue(configCode, req.OrgCode)
	if err != nil {
		return rsp, errCode, err
	}

	closeTradeArg := common.CloseTradeArg{
		OrderId:    req.OrderId,
		MethodCode: req.MethodCode,
		OrgCode:    req.OrgCode,
		Currency:   req.Currency,
		TotalFee:   req.TotalFee,
	}
	//关闭支付交易处理
	closeTradeRsp, errCode, err := t.Handler.CloseTrade(configParamMap, closeTradeArg)
	if err != nil {
		return rsp, errCode, err
	}
	rsp.OrderId = closeTradeRsp.OrderId
	rsp.Status = closeTradeRsp.Status

	return rsp, 0, nil
}
