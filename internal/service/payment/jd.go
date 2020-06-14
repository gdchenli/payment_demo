package payment

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"payment_demo/api/response"
	"payment_demo/api/validate"
	"payment_demo/internal/common/code"
	"payment_demo/pkg/config"
	consts2 "payment_demo/pkg/payment/consts"
	"strconv"
	"time"

	"github.com/gdchenli/pay/dialects/jd/payment"
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

//发起支付
func (jd *Jd) Pay(arg validate.Order) (form string, errCode int, err error) {
	//金额转为分
	totalFee := arg.TotalFee * 100
	//金额字段类型转换
	amount, err := strconv.ParseInt(fmt.Sprintf("%.f", totalFee), 10, 64)
	if err != nil {
		logrus.Errorf("org:jd,"+code.AmountFormatErrMessage+",errCode:%v,err:%v", code.AmountFormatErrCode, err.Error())
		return form, code.AmountFormatErrCode, errors.New(code.AmountFormatErrMessage)
	}

	privateKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPrivateKey))
	privateFile, err := os.Open(privateKeyPath)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PrivateKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PrivateKeyNotExistsErrCode, err.Error())
		return form, code.PrivateKeyNotExistsErrCode, errors.New(code.PrivateKeyNotExistsErrMessage)
	}
	privateKeyBytes, err := ioutil.ReadAll(privateFile)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PrivateKeyContentErrMessage+",errCode:%v,err:%v", code.PrivateKeyContentErrCode, err.Error())
		return form, code.PrivateKeyContentErrCode, errors.New(code.PrivateKeyContentErrMessage)
	}
	privateKey := string(privateKeyBytes)

	gateWay := jd.getPayWay(arg.UserAgentType)
	if gateWay == "" {
		logrus.Errorf("org:jd,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return form, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	merchant := config.GetInstance().GetString(JdMerchant)
	if merchant == "" {
		logrus.Errorf("org:jd,"+code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return form, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	desKey := config.GetInstance().GetString(JdDesKey)
	if desKey == "" {
		logrus.Errorf("org:jd,"+code.DesKeyNotExistsErrMessage+",errCode:%v,err:%v", code.DesKeyNotExistsErrCode)
		return form, code.DesKeyNotExistsErrCode, errors.New(code.DesKeyNotExistsErrMessage)
	}

	expireTime := config.GetInstance().GetString(JdExpireTime)
	if expireTime == "" {
		logrus.Errorf("org:jd,"+code.ExpireTimeNotExistsErrMessage+",errCode:%v,err:%v", code.ExpireTimeNotExistsErrCode)
		return form, code.ExpireTimeNotExistsErrCode, errors.New(code.ExpireTimeNotExistsErrMessage)
	}

	transCurrency := config.GetInstance().GetString(JdSettleCurrency)
	if transCurrency == "" {
		logrus.Errorf("org:jd,"+code.TransCurrencyNotExistsErrMessage+",errCode:%v,err:%v", code.TransCurrencyNotExistsErrCode)
		return form, code.TransCurrencyNotExistsErrCode, errors.New(code.TransCurrencyNotExistsErrMessage)
	}

	notifyUrl := config.GetInstance().GetString(JdNotifyUrl)
	if notifyUrl == "" {
		logrus.Errorf("org:jd,"+code.NotifyUrlNotExistsErrMessage+",errCode:%v,err:%v", code.NotifyUrlNotExistsErrCode)
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
	case consts2.WebUserAgentType:
		return config.GetInstance().GetString(JdPcPayWay)
	case consts2.MobileUserAgentType:
		return config.GetInstance().GetString(JdH5PayWay)
	}

	return config.GetInstance().GetString(JdPcPayWay)
}

//异步通知
func (jd *Jd) Notify(query, methodCode string) (notifyRsp response.NotifyRsp, errCode int, err error) {
	var jdNotifyRsp payment.NotifyRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,notify encrypt data:%+v,notify decrypt data:%v",
			jdNotifyRsp.OrderId, consts2.JdOrg, methodCode, jdNotifyRsp.EncryptRsp, jdNotifyRsp.DecryptRsp)
	}()

	publicKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPublicKey))
	publicKeyFile, err := os.Open(publicKeyPath)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PublicKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExistsErrCode, err.Error())
		return notifyRsp, code.PublicKeyNotExistsErrCode, errors.New(code.PublicKeyNotExistsErrMessage)
	}
	publicKeyBytes, err := ioutil.ReadAll(publicKeyFile)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PublicKeyContentErrMessage+",errCode:%v,err:%v", code.PublicKeyContentErrCode, err.Error())
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
	fmt.Printf("jdNotifyRsp：%+v\n", jdNotifyRsp)

	notifyRsp.OrderId = jdNotifyRsp.OrderId
	notifyRsp.Status = jdNotifyRsp.Status
	notifyRsp.TradeNo = jdNotifyRsp.TradeNo
	notifyRsp.Message = "success"
	notifyRsp.RmbFee = jdNotifyRsp.RmbFee

	return notifyRsp, 0, nil
}

//同步通知
func (jd *Jd) Verify(query, methodCode string) (verifyRsp response.VerifyRsp, errCode int, err error) {
	var jdCallbackRsp payment.CallbackRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,callback encrypt data:%+v,callback decrypt data:%v",
			jdCallbackRsp.OrderId, consts2.JdOrg, consts2.JdMethod, jdCallbackRsp.EncryptRsp, jdCallbackRsp.DecryptRsp)
	}()

	publicKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPublicKey))
	file, err := os.Open(publicKeyPath)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PublicKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExistsErrCode, err.Error())
		return verifyRsp, code.PublicKeyNotExistsErrCode, errors.New(code.PublicKeyNotExistsErrMessage)
	}
	publicKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PublicKeyContentErrMessage+",errCode:%v,err:%v", code.PublicKeyContentErrCode, err.Error())
		return verifyRsp, code.PublicKeyContentErrCode, errors.New(code.PublicKeyContentErrMessage)
	}
	publicKey := string(publicKeyBytes)

	callbackArg := payment.CallbackArg{
		PublicKey: publicKey,
		DesKey:    config.GetInstance().GetString(JdDesKey),
	}
	jdCallbackRsp, errCode, err = new(payment.Callback).Validate(query, callbackArg)
	if err != nil {
		return verifyRsp, errCode, err
	}

	verifyRsp.Status = jdCallbackRsp.Status
	verifyRsp.OrderId = jdCallbackRsp.OrderId

	return verifyRsp, 0, err
}

//查询交易
func (jd *Jd) SearchTrade(orderId, methodCode, currency string, totalFee float64) (searchTradeRsp response.SearchTradeRsp, errCode int, err error) {
	var jdTradeRsp payment.TradeRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,trade search request encrypt data:%+v,trade search request decrypt data:%v"+
			",trade search response search encrypt data:%v,trade search response search decrypt data:%v",
			jdTradeRsp.OrderId, consts2.JdOrg, methodCode, jdTradeRsp.EncryptRes, jdTradeRsp.DecryptRes, jdTradeRsp.EncryptRsp, jdTradeRsp.DecryptRsp)

	}()

	privateKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPrivateKey))
	privateFile, err := os.Open(privateKeyPath)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PrivateKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PrivateKeyNotExistsErrCode, err.Error())
		return searchTradeRsp, code.PrivateKeyNotExistsErrCode, errors.New(code.PrivateKeyNotExistsErrMessage)
	}
	privateKeyBytes, err := ioutil.ReadAll(privateFile)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PrivateKeyContentErrMessage+",errCode:%v,err:%v", code.PrivateKeyContentErrCode, err.Error())
		return searchTradeRsp, code.PrivateKeyContentErrCode, errors.New(code.PrivateKeyContentErrMessage)
	}
	privateKey := string(privateKeyBytes)

	publicKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPublicKey))
	file, err := os.Open(publicKeyPath)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PublicKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExistsErrCode, err.Error())
		return searchTradeRsp, code.PublicKeyNotExistsErrCode, errors.New(code.PublicKeyNotExistsErrMessage)
	}
	publicKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PublicKeyContentErrMessage+",errCode:%v,err:%v", code.PublicKeyContentErrCode, err.Error())
		return searchTradeRsp, code.PublicKeyContentErrCode, errors.New(code.PublicKeyContentErrMessage)
	}
	publicKey := string(publicKeyBytes)

	merchant := config.GetInstance().GetString(JdMerchant)
	if merchant == "" {
		logrus.Errorf("org:jd,"+code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return searchTradeRsp, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	desKey := config.GetInstance().GetString(JdDesKey)
	if desKey == "" {
		logrus.Errorf("org:jd,"+code.DesKeyNotExistsErrMessage+",errCode:%v,err:%v", code.DesKeyNotExistsErrCode)
		return searchTradeRsp, code.DesKeyNotExistsErrCode, errors.New(code.DesKeyNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(JdTradeWay)
	if gateWay == "" {
		logrus.Errorf("org:jd,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return searchTradeRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	tradeArg := payment.TradeArg{
		Merchant:   merchant,
		TradeNum:   orderId,
		DesKey:     desKey,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		GateWay:    gateWay,
		TotalFee:   totalFee,
	}
	jdTradeRsp, errCode, err = new(payment.Trade).Search(tradeArg)
	if err != nil {
		return searchTradeRsp, errCode, err
	}

	searchTradeRsp.OrderId = jdTradeRsp.OrderId
	searchTradeRsp.Status = jdTradeRsp.Status
	searchTradeRsp.TradeNo = jdTradeRsp.TradeNo
	searchTradeRsp.RmbFee = jdTradeRsp.RmbFee

	return searchTradeRsp, 0, nil
}

//关闭交易
func (jd *Jd) CloseTrade(arg validate.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {
	var jdClosedRsp payment.ClosedRsp
	defer func() {
		//记录日志
		logrus.Info("order id:%v,org:%v,method:%v,closed trade request encrypt data:%+v,closed trade request decrypt data:%v"+
			",closed trade response encrypt data:%v,closed trade response decrypt data:%v",
			jdClosedRsp.OrderId, consts2.JdOrg, consts2.JdMethod, jdClosedRsp.EncryptRes, jdClosedRsp.DecryptRes, jdClosedRsp.EncryptRsp, jdClosedRsp.DecryptRsp)
	}()

	//金额转为分
	totalFee := arg.TotalFee * 100
	//金额字段类型转换
	amount, err := strconv.ParseInt(fmt.Sprintf("%.f", totalFee), 10, 64)
	if err != nil {
		logrus.Errorf("org:jd,"+code.AmountFormatErrMessage+",errCode:%v,err:%v", code.AmountFormatErrCode, err.Error())
		return closeTradeRsp, code.AmountFormatErrCode, errors.New(code.AmountFormatErrMessage)
	}

	privateKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPrivateKey))
	privateFile, err := os.Open(privateKeyPath)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PrivateKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PrivateKeyNotExistsErrCode, err.Error())
		return closeTradeRsp, code.PrivateKeyNotExistsErrCode, errors.New(code.PrivateKeyNotExistsErrMessage)
	}
	privateKeyBytes, err := ioutil.ReadAll(privateFile)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PrivateKeyContentErrMessage+",errCode:%v,err:%v", code.PrivateKeyContentErrCode, err.Error())
		return closeTradeRsp, code.PrivateKeyContentErrCode, errors.New(code.PrivateKeyContentErrMessage)
	}
	privateKey := string(privateKeyBytes)

	publicKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPublicKey))
	file, err := os.Open(publicKeyPath)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PublicKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExistsErrCode, err.Error())
		return closeTradeRsp, code.PublicKeyNotExistsErrCode, errors.New(code.PublicKeyNotExistsErrMessage)
	}
	publicKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PublicKeyContentErrMessage+",errCode:%v,err:%v", code.PublicKeyContentErrCode, err.Error())
		return closeTradeRsp, code.PublicKeyContentErrCode, errors.New(code.PublicKeyContentErrMessage)
	}
	publicKey := string(publicKeyBytes)

	merchant := config.GetInstance().GetString(JdMerchant)
	if merchant == "" {
		logrus.Errorf("org:jd,"+code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return closeTradeRsp, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	desKey := config.GetInstance().GetString(JdDesKey)
	if desKey == "" {
		logrus.Errorf("org:jd,"+code.DesKeyNotExistsErrMessage+",errCode:%v,err:%v", code.DesKeyNotExistsErrCode)
		return closeTradeRsp, code.DesKeyNotExistsErrCode, errors.New(code.DesKeyNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(JdClosedWay)
	if gateWay == "" {
		logrus.Errorf("org:jd,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return closeTradeRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
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
		return closeTradeRsp, errCode, err
	}

	closeTradeRsp.OrderId = jdClosedRsp.OrderId
	closeTradeRsp.Status = jdClosedRsp.Status

	return closeTradeRsp, 0, nil
}

//物流信息上传
func (jd *Jd) UploadLogistics(arg validate.UploadLogisticsReq) (uploadLogisticsRsp response.UploadLogisticsRsp, errCode int, err error) {
	var jdLogisticsRsp payment.LogisticsRsp
	defer func() {
		//记录日志
		logrus.Info("order id:%v,org:%v,method:%v,upload logistics request encrypt data:%+v,upload logistics request decrypt data:%v"+
			",upload logistics response encrypt data:%v,upload logistics response decrypt data:%v",
			jdLogisticsRsp.OrderId, consts2.JdOrg, consts2.JdMethod, jdLogisticsRsp.EncryptRes, jdLogisticsRsp.DecryptRes, jdLogisticsRsp.EncryptRsp, jdLogisticsRsp.DecryptRsp)
	}()

	privateKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPrivateKey))
	privateFile, err := os.Open(privateKeyPath)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PrivateKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PrivateKeyNotExistsErrCode, err.Error())
		return uploadLogisticsRsp, code.PrivateKeyNotExistsErrCode, errors.New(code.PrivateKeyNotExistsErrMessage)
	}
	privateKeyBytes, err := ioutil.ReadAll(privateFile)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PrivateKeyContentErrMessage+",errCode:%v,err:%v", code.PrivateKeyContentErrCode, err.Error())
		return uploadLogisticsRsp, code.PrivateKeyContentErrCode, errors.New(code.PrivateKeyContentErrMessage)
	}
	privateKey := string(privateKeyBytes)

	publicKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPublicKey))
	file, err := os.Open(publicKeyPath)
	if err != nil {
		logrus.Errorf("org:jd,"+code.PublicKeyNotExistsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExistsErrCode, err.Error())
		return uploadLogisticsRsp, code.PublicKeyNotExistsErrCode, errors.New(code.PublicKeyNotExistsErrMessage)
	}
	publicKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf(code.PublicKeyContentErrMessage+",errCode:%v,err:%v", code.PublicKeyContentErrCode, err.Error())
		return uploadLogisticsRsp, code.PublicKeyContentErrCode, errors.New(code.PublicKeyContentErrMessage)
	}
	publicKey := string(publicKeyBytes)

	merchant := config.GetInstance().GetString(JdMerchant)
	if merchant == "" {
		logrus.Errorf("org:jd,"+code.MerchantNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNotExistsErrCode)
		return uploadLogisticsRsp, code.MerchantNotExistsErrCode, errors.New(code.MerchantNotExistsErrMessage)
	}

	desKey := config.GetInstance().GetString(JdDesKey)
	if desKey == "" {
		logrus.Errorf("org:jd,"+code.DesKeyNotExistsErrMessage+",errCode:%v,err:%v", code.DesKeyNotExistsErrCode)
		return uploadLogisticsRsp, code.DesKeyNotExistsErrCode, errors.New(code.DesKeyNotExistsErrMessage)
	}

	gateWay := config.GetInstance().GetString(JdLogisticsWay)
	if gateWay == "" {
		logrus.Errorf("org:jd,"+code.GateWayNotExistsErrMessage+",errCode:%v,err:%v", code.GateWayNotExistsErrCode)
		return uploadLogisticsRsp, code.GateWayNotExistsErrCode, errors.New(code.GateWayNotExistsErrMessage)
	}

	merchantName := config.GetInstance().GetString(JdMerchantName)
	if merchantName == "" {
		logrus.Errorf("org:jd,"+code.MerchantNameNotExistsErrMessage+",errCode:%v,err:%v", code.MerchantNameNotExistsErrCode)
		return uploadLogisticsRsp, code.MerchantNameNotExistsErrCode, errors.New(code.MerchantNameNotExistsErrMessage)
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
	jdLogisticsRsp, errCode, err = new(payment.Logistics).Upload(logisticsArg)
	if err != nil {
		return uploadLogisticsRsp, errCode, err
	}

	uploadLogisticsRsp.Status = jdLogisticsRsp.Status
	uploadLogisticsRsp.OrderId = jdLogisticsRsp.OrderId

	return uploadLogisticsRsp, 0, nil
}
