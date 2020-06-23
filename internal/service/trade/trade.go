package trade

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"payment_demo/api/trade/request"
	"payment_demo/api/trade/response"
	"payment_demo/internal/common/code"
	"payment_demo/pkg/config"
)

type Trade struct{ Handler Handler }

func New(orgCode string) (*Trade, int, error) {
	trade := new(Trade)

	trade.Handler = getHandler(orgCode)
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
func (t *Trade) SearchTrade(req request.SearchTradeArg) (searchTradeRsp response.SearchTradeRsp, errCode int, err error) {
	//获取配置项code
	configCode := t.Handler.GetSearchTradeConfigCode()

	//读取配置项值
	configParamMap, errCode, err := t.getConfigValue(configCode, req.OrgCode)
	if err != nil {
		return searchTradeRsp, errCode, err
	}

	searchTradeRsp, errCode, err = t.Handler.SearchTrade(configParamMap, req)
	if err != nil {
		return searchTradeRsp, errCode, err
	}

	return searchTradeRsp, 0, nil
}

func (t *Trade) CloseTrade(req request.CloseTradeArg) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {
	//获取配置项code
	configCode := t.Handler.GetCloseTradeConfigCode()

	//读取配置项值
	configParamMap, errCode, err := t.getConfigValue(configCode, req.OrgCode)
	if err != nil {
		return closeTradeRsp, errCode, err
	}

	//关闭支付交易处理
	closeTradeRsp, errCode, err = t.Handler.CloseTrade(configParamMap, req)
	if err != nil {
		return closeTradeRsp, errCode, err
	}

	return closeTradeRsp, 0, nil
}
