package method

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/config"
	"payment_demo/internal/common/defs"
	"payment_demo/pkg/jd/payment"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	JdMerchant       = "jd.merchant"
	JdCallBackUrl    = "jd.callback_url"
	JdNotifyUrl      = "jd.notify_url"
	JdExpireTime     = "jd.expire_time"
	JdSettleCurrency = "jd.settle_currency"
	JdDesKey         = "jd.des_key"
	JdPcPayWay       = "jd.pc_pay_way"
	JdH5PayWay       = "jd.h5_pay_way"
	JdPrivateKey     = "jd.private_key"
	JdPublicKey      = "jd.public_key"
	JdTradeWay       = "jd.trade_way"
	JdClosedWay      = "jd.closed_way"
	JdLogisticsWay   = "jd.logistics_way"
	JdMerchantName   = "jd.merchant_name"
)

type Jd struct{}

func (jd *Jd) OrderQrCode(arg defs.Order) (form string, errCode int, err error) {
	return form, code.NotSupportPaymentMethodErrCode, errors.New(code.NotSupportPaymentMethodErrMessage)
}

func (jd *Jd) AmpSubmit(arg defs.Order) (form string, errCode int, err error) {
	return form, code.NotSupportPaymentMethodErrCode, errors.New(code.NotSupportPaymentMethodErrMessage)
}

//发起支付
func (jd *Jd) OrderSubmit(arg defs.Order) (form string, errCode int, err error) {
	//金额转为分
	totalFee := arg.TotalFee * 100
	//金额字段类型转换
	amount, err := strconv.ParseInt(fmt.Sprintf("%.f", totalFee), 10, 64)
	if err != nil {
		logrus.Errorf(code.AmountFormatErrMessage+",errCode:%v,err:%v", code.AmountFormatErrCode, err.Error())
		return form, code.AmountFormatErrCode, errors.New(code.AmountFormatErrMessage)
	}

	privateKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPrivateKey))
	privateFile, err := os.Open(privateKeyPath)
	if err != nil {
		logrus.Errorf(code.PrivateKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PrivateKeyNotExistsErrCode, err.Error())
		return form, code.PrivateKeyNotExistsErrCode, errors.New(code.PrivateKeyNotExistsErrMessage)
	}
	privateKeyBytes, err := ioutil.ReadAll(privateFile)
	if err != nil {
		logrus.Errorf(code.PrivateKeyContentErrMessage+",errCode:%v,err:%v", code.PrivateKeyContentErrCode, err.Error())
		return form, code.PrivateKeyContentErrCode, errors.New(code.PrivateKeyContentErrMessage)
	}
	privateKey := string(privateKeyBytes)

	gateWay := jd.getPayWay(arg.UserAgentType)
	if gateWay == "" {
		logrus.Errorf(code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return form, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	merchant := config.GetInstance().GetString(JdMerchant)
	if merchant == "" {
		logrus.Errorf(code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return form, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	desKey := config.GetInstance().GetString(JdDesKey)
	if desKey == "" {
		logrus.Errorf(code.DesKeyNotExistsErrMessage+",errCode:%v,err:%v", code.DesKeyNotExistsErrCode)
		return form, code.DesKeyNotExistsErrCode, errors.New(code.DesKeyNotExistsErrMessage)
	}

	expireTime := config.GetInstance().GetString(JdExpireTime)
	if expireTime == "" {
		logrus.Errorf(code.ExpireTimeNotExistsErrMessage+",errCode:%v,err:%v", code.ExpireTimeNotExistsErrCode)
		return form, code.ExpireTimeNotExistsErrCode, errors.New(code.ExpireTimeNotExistsErrMessage)
	}

	transCurrency := config.GetInstance().GetString(JdSettleCurrency)
	if transCurrency == "" {
		logrus.Errorf(code.TransCurrencyNotExistsErrMessage+",errCode:%v,err:%v", code.TransCurrencyNotExistsErrCode)
		return form, code.TransCurrencyNotExistsErrCode, errors.New(code.TransCurrencyNotExistsErrMessage)
	}

	notifyUrl := config.GetInstance().GetString(JdNotifyUrl)
	if notifyUrl == "" {
		logrus.Errorf(code.NotifyUrlNotExistsErrMessage+",errCode:%v,err:%v", code.NotifyUrlNotExistsErrCode)
		return form, code.NotifyUrlNotExistsErrCode, errors.New(code.NotifyUrlNotExistsErrMessage)
	}

	callbackUrl := config.GetInstance().GetString(JdCallBackUrl)
	if callbackUrl == "" {
		logrus.Errorf(code.CallbackUrlNotExistsErrMessage+",errCode:%v,err:%v", code.CallbackUrlNotExistsErrCode)
		return form, code.CallbackUrlNotExistsErrCode, errors.New(code.CallbackUrlNotExistsErrMessage)
	}

	date := time.Now().Format(payment.TimeLayout)
	goodsInfos := []payment.GoodsInfo{{Id: "test" + date, Name: "test" + date, Price: amount, Num: 1}}
	kjInfo := payment.KjInfo{GoodsSubmittedCustoms: "N", GoodsUnderBonded: "N"}
	payArg := payment.PayArg{
		Merchant:      merchant,
		TradeNum:      arg.OrderId,
		Amount:        amount,
		Currency:      arg.Currency,
		CallbackUrl:   callbackUrl,
		NotifyUrl:     notifyUrl,
		UserId:        arg.UserId,
		ExpireTime:    expireTime,
		TransCurrency: transCurrency,
		GoodsInfo:     goodsInfos,
		KjInfo:        kjInfo,
		DesKey:        desKey,
		PrivateKey:    privateKey,
		PayWay:        gateWay,
		TradeName:     arg.OrderId,
	}
	form, errCode, err = new(payment.Payment).CreateForm(payArg)
	if err != nil {
		return form, errCode, err
	}

	return form, 0, nil
}

func (jd *Jd) getPayWay(userAgentType int) string {
	switch userAgentType {
	case code.WebUserAgentType:
		return config.GetInstance().GetString(JdPcPayWay)
	case code.MobileUserAgentType:
		return config.GetInstance().GetString(JdH5PayWay)
	}

	return config.GetInstance().GetString(JdPcPayWay)
}

//异步通知
func (jd *Jd) Notify(query, methodCode string) (notifyRsp defs.NotifyRsp, errCode int, err error) {
	var jdNotifyRsp payment.NotifyRsp
	defer func() {
		//记录日志
		logrus.Info("order id:%v,org:%v,method:%v,notify encrypt data:%+v,notify decrypt data:%v",
			jdNotifyRsp.OrderId, code.JdOrg, methodCode, jdNotifyRsp.EncryptRsp, jdNotifyRsp.DecryptRsp)
	}()

	publicKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPublicKey))
	publicKeyFile, err := os.Open(publicKeyPath)
	if err != nil {
		logrus.Errorf(code.PublicKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExistsErrCode, err.Error())
		return notifyRsp, code.PublicKeyNotExistsErrCode, errors.New(code.PublicKeyNotExistsErrMessage)
	}
	publicKeyBytes, err := ioutil.ReadAll(publicKeyFile)
	if err != nil {
		logrus.Errorf(code.PublicKeyContentErrMessage+",errCode:%v,err:%v", code.PublicKeyContentErrCode, err.Error())
		return notifyRsp, code.PublicKeyContentErrCode, errors.New(code.PublicKeyContentErrMessage)
	}
	publicKey := string(publicKeyBytes)
	notifyArg := payment.NotifyArg{
		PublicKey: publicKey,
		DesKey:    config.GetInstance().GetString(JdDesKey),
	}

	jdNotifyRsp, errCode, err = new(payment.Notify).Validate(query, notifyArg)
	if err != nil {
		return notifyRsp, errCode, err
	}

	notifyRsp.OrderId = jdNotifyRsp.OrderId
	notifyRsp.Status = jdNotifyRsp.Status
	notifyRsp.TradeNo = jdNotifyRsp.TradeNo
	notifyRsp.Message = "success"

	return notifyRsp, 0, nil
}

//同步通知
func (jd *Jd) Callback(query, methodCode string) (callbackRsp defs.CallbackRsp, errCode int, err error) {
	var jdCallbackRsp payment.CallbackRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,callback encrypt data:%+v,callback decrypt data:%v",
			jdCallbackRsp.OrderId, code.JdOrg, code.JdMethod, jdCallbackRsp.EncryptRsp, jdCallbackRsp.DecryptRsp)
	}()

	publicKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPublicKey))
	file, err := os.Open(publicKeyPath)
	if err != nil {
		logrus.Errorf(code.PublicKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExistsErrCode, err.Error())
		return callbackRsp, code.PublicKeyNotExistsErrCode, errors.New(code.PublicKeyNotExistsErrMessage)
	}
	publicKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf(code.PublicKeyContentErrMessage+",errCode:%v,err:%v", code.PublicKeyContentErrCode, err.Error())
		return callbackRsp, code.PublicKeyContentErrCode, errors.New(code.PublicKeyContentErrMessage)
	}
	publicKey := string(publicKeyBytes)

	callbackArg := payment.CallbackArg{
		PublicKey: publicKey,
		DesKey:    config.GetInstance().GetString(JdDesKey),
	}
	jdCallbackRsp, errCode, err = new(payment.Callback).Validate(query, callbackArg)
	if err != nil {
		return callbackRsp, errCode, err
	}

	callbackRsp.Status = jdCallbackRsp.Status
	callbackRsp.OrderId = jdCallbackRsp.OrderId

	return callbackRsp, 0, err
}

//查询交易
func (jd *Jd) Trade(orderId, methodCode string) (tradeRsp defs.TradeRsp, errCode int, err error) {
	var jdTradeRsp payment.TradeRsp
	defer func() {
		//记录日志
		logrus.Info("order id:%v,org:%v,method:%v,trade search request encrypt data:%+v,trade search request decrypt data:%v"+
			",trade search response search encrypt data:%v,trade search response search decrypt data:%v",
			jdTradeRsp.OrderId, code.JdOrg, methodCode, jdTradeRsp.EncryptRes, jdTradeRsp.DecryptRes, jdTradeRsp.EncryptRsp, jdTradeRsp.DecryptRsp)

	}()

	privateKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPrivateKey))
	privateFile, err := os.Open(privateKeyPath)
	if err != nil {
		logrus.Errorf(code.PrivateKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PrivateKeyNotExistsErrCode, err.Error())
		return tradeRsp, code.PrivateKeyNotExistsErrCode, errors.New(code.PrivateKeyNotExistsErrMessage)
	}
	privateKeyBytes, err := ioutil.ReadAll(privateFile)
	if err != nil {
		logrus.Errorf(code.PrivateKeyContentErrMessage+",errCode:%v,err:%v", code.PrivateKeyContentErrCode, err.Error())
		return tradeRsp, code.PrivateKeyContentErrCode, errors.New(code.PrivateKeyContentErrMessage)
	}
	privateKey := string(privateKeyBytes)

	publicKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPublicKey))
	file, err := os.Open(publicKeyPath)
	if err != nil {
		logrus.Errorf(code.PublicKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExistsErrCode, err.Error())
		return tradeRsp, code.PublicKeyNotExistsErrCode, errors.New(code.PublicKeyNotExistsErrMessage)
	}
	publicKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf(code.PublicKeyContentErrMessage+",errCode:%v,err:%v", code.PublicKeyContentErrCode, err.Error())
		return tradeRsp, code.PublicKeyContentErrCode, errors.New(code.PublicKeyContentErrMessage)
	}
	publicKey := string(publicKeyBytes)

	merchant := config.GetInstance().GetString(JdMerchant)
	if merchant == "" {
		logrus.Errorf(code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return tradeRsp, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	desKey := config.GetInstance().GetString(JdDesKey)
	if desKey == "" {
		logrus.Errorf(code.DesKeyNotExistsErrMessage+",errCode:%v,err:%v", code.DesKeyNotExistsErrCode)
		return tradeRsp, code.DesKeyNotExistsErrCode, errors.New(code.DesKeyNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(JdTradeWay)
	if gateWay == "" {
		logrus.Errorf(code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return tradeRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	tradeArg := payment.TradeArg{
		Merchant:   merchant,
		TradeNum:   orderId,
		DesKey:     desKey,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		GateWay:    gateWay,
	}
	jdTradeRsp, errCode, err = new(payment.Trade).Search(tradeArg)
	if err != nil {
		return tradeRsp, errCode, err
	}

	tradeRsp.OrderId = jdTradeRsp.OrderId
	tradeRsp.Status = jdTradeRsp.Status
	tradeRsp.TradeNo = jdTradeRsp.TradeNo

	return tradeRsp, 0, nil
}

type JdClosedArg struct {
	OrderId  string
	Currency string
	TotalFee float64
}

//关闭交易
func (jd *Jd) Closed(arg JdClosedArg) (closedRsp defs.ClosedRsp, errCode int, err error) {
	var jdClosedRsp payment.ClosedRsp
	defer func() {
		//记录日志
		logrus.Info("order id:%v,org:%v,method:%v,closed trade request encrypt data:%+v,closed trade request decrypt data:%v"+
			",closed trade response encrypt data:%v,closed trade response decrypt data:%v",
			jdClosedRsp.OrderId, code.JdOrg, code.JdMethod, jdClosedRsp.EncryptRes, jdClosedRsp.DecryptRes, jdClosedRsp.EncryptRsp, jdClosedRsp.DecryptRsp)
	}()

	//金额转为分
	totalFee := arg.TotalFee * 100
	//金额字段类型转换
	amount, err := strconv.ParseInt(fmt.Sprintf("%.f", totalFee), 10, 64)
	if err != nil {
		logrus.Errorf(code.AmountFormatErrMessage+",errCode:%v,err:%v", code.AmountFormatErrCode, err.Error())
		return closedRsp, code.AmountFormatErrCode, errors.New(code.AmountFormatErrMessage)
	}

	privateKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPrivateKey))
	privateFile, err := os.Open(privateKeyPath)
	if err != nil {
		logrus.Errorf(code.PrivateKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PrivateKeyNotExistsErrCode, err.Error())
		return closedRsp, code.PrivateKeyNotExistsErrCode, errors.New(code.PrivateKeyNotExistsErrMessage)
	}
	privateKeyBytes, err := ioutil.ReadAll(privateFile)
	if err != nil {
		logrus.Errorf(code.PrivateKeyContentErrMessage+",errCode:%v,err:%v", code.PrivateKeyContentErrCode, err.Error())
		return closedRsp, code.PrivateKeyContentErrCode, errors.New(code.PrivateKeyContentErrMessage)
	}
	privateKey := string(privateKeyBytes)

	publicKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPublicKey))
	file, err := os.Open(publicKeyPath)
	if err != nil {
		logrus.Errorf(code.PublicKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExistsErrCode, err.Error())
		return closedRsp, code.PublicKeyNotExistsErrCode, errors.New(code.PublicKeyNotExistsErrMessage)
	}
	publicKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf(code.PublicKeyContentErrMessage+",errCode:%v,err:%v", code.PublicKeyContentErrCode, err.Error())
		return closedRsp, code.PublicKeyContentErrCode, errors.New(code.PublicKeyContentErrMessage)
	}
	publicKey := string(publicKeyBytes)

	merchant := config.GetInstance().GetString(JdMerchant)
	if merchant == "" {
		logrus.Errorf(code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return closedRsp, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	desKey := config.GetInstance().GetString(JdDesKey)
	if desKey == "" {
		logrus.Errorf(code.DesKeyNotExistsErrMessage+",errCode:%v,err:%v", code.DesKeyNotExistsErrCode)
		return closedRsp, code.DesKeyNotExistsErrCode, errors.New(code.DesKeyNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(JdClosedWay)
	if gateWay == "" {
		logrus.Errorf(code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return closedRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	closedArg := payment.ClosedArg{
		Merchant:   merchant,
		TradeNum:   arg.OrderId + "jd",
		OTradeNum:  arg.OrderId,
		Amount:     amount,
		Currency:   arg.Currency,
		DesKey:     desKey,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		GateWay:    gateWay,
	}
	jdClosedRsp, errCode, err = new(payment.Closed).Trade(closedArg)
	if err != nil {
		return closedRsp, errCode, err
	}

	closedRsp.OrderId = jdClosedRsp.OrderId
	closedRsp.Status = jdClosedRsp.Status

	return closedRsp, 0, nil
}

type JdLogisticsArg struct {
	OrderId          string `json:"order_id"`          //订单编号
	LogisticsNo      string `json:"logistics_no"`      //物流单号
	LogisticsCompany string `json:"logistics_company"` //物流公司名称
}

//物流信息上传
func (jd *Jd) LogisticsUpload(arg JdLogisticsArg) (logisticsRsp defs.LogisticsRsp, errCode int, err error) {
	privateKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPrivateKey))
	privateFile, err := os.Open(privateKeyPath)
	if err != nil {
		logrus.Errorf(code.PrivateKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PrivateKeyNotExistsErrCode, err.Error())
		return logisticsRsp, code.PrivateKeyNotExistsErrCode, errors.New(code.PrivateKeyNotExistsErrMessage)
	}
	privateKeyBytes, err := ioutil.ReadAll(privateFile)
	if err != nil {
		logrus.Errorf(code.PrivateKeyContentErrMessage+",errCode:%v,err:%v", code.PrivateKeyContentErrCode, err.Error())
		return logisticsRsp, code.PrivateKeyContentErrCode, errors.New(code.PrivateKeyContentErrMessage)
	}
	privateKey := string(privateKeyBytes)

	publicKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPublicKey))
	file, err := os.Open(publicKeyPath)
	if err != nil {
		logrus.Errorf(code.PublicKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExistsErrCode, err.Error())
		return logisticsRsp, code.PublicKeyNotExistsErrCode, errors.New(code.PublicKeyNotExistsErrMessage)
	}
	publicKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf(code.PublicKeyContentErrMessage+",errCode:%v,err:%v", code.PublicKeyContentErrCode, err.Error())
		return logisticsRsp, code.PublicKeyContentErrCode, errors.New(code.PublicKeyContentErrMessage)
	}
	publicKey := string(publicKeyBytes)

	merchant := config.GetInstance().GetString(JdMerchant)
	if merchant == "" {
		logrus.Errorf(code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return logisticsRsp, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	desKey := config.GetInstance().GetString(JdDesKey)
	if desKey == "" {
		logrus.Errorf(code.DesKeyNotExistsErrMessage+",errCode:%v,err:%v", code.DesKeyNotExistsErrCode)
		return logisticsRsp, code.DesKeyNotExistsErrCode, errors.New(code.DesKeyNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(JdLogisticsWay)
	if gateWay == "" {
		logrus.Errorf(code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return logisticsRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	merchantName := config.GetInstance().GetString(JdMerchantName)
	if merchantName == "" {
		logrus.Errorf(code.MerchantNameNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNameNotExistsErrCode)
		return logisticsRsp, code.MerchantNameNotExistsErrCode, errors.New(code.MerchantNameNotExistsErrMessage)
	}

	logisticsArg := payment.LogisticsArg{
		OrderId:          arg.OrderId,
		LogisticsNo:      arg.LogisticsNo,
		LogisticsCompany: arg.LogisticsCompany,
		Merchant:         merchant,
		PrivateKey:       privateKey,
		PublicKey:        publicKey,
		DesKey:           desKey,
		MerchantName:     merchantName,
		GateWay:          gateWay,
	}
	jdLogisticsRsp, errCode, err := new(payment.Logistics).Upload(logisticsArg)
	if err != nil {
		return logisticsRsp, errCode, err
	}

	logisticsRsp.Status = jdLogisticsRsp.Status
	logisticsRsp.OrderId = jdLogisticsRsp.OrderId

	return logisticsRsp, 0, nil
}
