package logistics

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
	"payment_demo/pkg/payment/jd"
)

type Logistics struct{}

func New() (*Logistics, int, error) {
	payment := new(Logistics)

	return payment, 0, nil
}

func (logistics *Logistics) getConfigValue(configCodes []string, orgCode string) (payParamMap map[string]string, errCode int, err error) {
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

func (logistics *Logistics) Upload(req request.UploadLogisticsArg) (uploadLogisticsTradeRsp response.UploadLogisticsRsp, errCode int, err error) {
	//获取配置项code
	jdPayment := jd.New()
	configCode := jdPayment.GetUploadLogisticsConfigCode()

	//读取配置项值
	configParamMap, errCode, err := logistics.getConfigValue(configCode, req.OrgCode)
	if err != nil {
		return uploadLogisticsTradeRsp, errCode, err
	}

	//上传物流信息处理
	uploadLogisticsTradeRsp, errCode, err = jdPayment.UploadLogistics(configParamMap, req)
	if err != nil {
		return uploadLogisticsTradeRsp, errCode, err
	}

	return uploadLogisticsTradeRsp, 0, nil
}
