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
	JdPrivateKey     = "jd.private_key"
	JdPublicKey      = "jd.public_key"
	JdTradeWay       = "jd.trade_way"
	JdClosedWay      = "jd.closed_way"
)

type Jd struct{}

type JdPayArg struct {
	OrderId  string  `json:"order_id"`
	TotalFee float64 `json:"total_fee"`
	Currency string  `json:"currency"`
	UserId   string  `json:"user_id"`
}

//发起支付
func (jd *Jd) Submit(arg JdPayArg) (form string, errCode int, err error) {
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
		logrus.Errorf(code.PrivateKeyNotExitsErrMessage+",errCode:%v,err:%v", code.PrivateKeyNotExitsErrCode, err.Error())
		return form, code.PrivateKeyNotExitsErrCode, errors.New(code.PrivateKeyNotExitsErrMessage)
	}
	privateKeyBytes, err := ioutil.ReadAll(privateFile)
	if err != nil {
		logrus.Errorf(code.PrivateKeyContentErrMessage+",errCode:%v,err:%v", code.PrivateKeyContentErrCode, err.Error())
		return form, code.PrivateKeyContentErrCode, errors.New(code.PrivateKeyContentErrMessage)
	}
	privateKey := string(privateKeyBytes)

	date := time.Now().Format(payment.TimeLayout)
	goodsInfos := []payment.GoodsInfo{{Id: "test" + date, Name: "test" + date, Price: amount, Num: 1}}
	kjInfo := payment.KjInfo{GoodsSubmittedCustoms: "N", GoodsUnderBonded: "N"}
	payArg := payment.PayArg{
		Merchant:      config.GetInstance().GetString(JdMerchant),
		TradeNum:      arg.OrderId,
		Amount:        amount,
		Currency:      arg.Currency,
		CallbackUrl:   config.GetInstance().GetString(JdCallBackUrl),
		NotifyUrl:     config.GetInstance().GetString(JdNotifyUrl),
		UserId:        arg.UserId,
		ExpireTime:    config.GetInstance().GetString(JdExpireTime),
		TransCurrency: config.GetInstance().GetString(JdSettleCurrency),
		GoodsInfo:     goodsInfos,
		KjInfo:        kjInfo,
		DesKey:        config.GetInstance().GetString(JdDesKey),
		PrivateKey:    privateKey,
		PayWay:        config.GetInstance().GetString(JdPcPayWay),
		BizTp:         "100006",
		TradeName:     arg.OrderId,
	}
	form, errCode, err = new(payment.Payment).CreatePaymentForm(payArg)
	if err != nil {
		return form, errCode, err
	}

	return form, 0, nil
}

//异步通知
func (jd *Jd) Notify(query string) (notifyRsp defs.NotifyRsp, errCode int, err error) {
	var jdNotifyRsp payment.NotifyRsp
	defer func() {
		//记录日志
		logrus.Info("order id:%v,org:%v,method:%v,notify encrypt data:%+v,notify decrypt data:%v",
			jdNotifyRsp.OrderId, "jd", "jd_payment", jdNotifyRsp.EncryptRsp, jdNotifyRsp.DecryptRsp)
	}()

	publicKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPublicKey))
	publicKeyFile, err := os.Open(publicKeyPath)
	if err != nil {
		logrus.Errorf(code.PublicKeyNotExitsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExitsErrCode, err.Error())
		return notifyRsp, code.PublicKeyNotExitsErrCode, errors.New(code.PublicKeyNotExitsErrMessage)
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

	return notifyRsp, 0, nil
}

//同步通知
func (jd *Jd) Callback(query string) (callbackRsp defs.CallbackRsp, errCode int, err error) {
	var jdCallbackRsp payment.CallbackRsp
	defer func() {
		//记录日志
		logrus.Infof("order id:%v,org:%v,method:%v,callback encrypt data:%+v,callback decrypt data:%v",
			jdCallbackRsp.OrderId, "jd", "jd_payment", jdCallbackRsp.EncryptRsp, jdCallbackRsp.DecryptRsp)
	}()

	publicKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPublicKey))
	file, err := os.Open(publicKeyPath)
	if err != nil {
		logrus.Errorf(code.PublicKeyNotExitsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExitsErrCode, err.Error())
		return callbackRsp, code.PublicKeyNotExitsErrCode, errors.New(code.PublicKeyNotExitsErrMessage)
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
func (jd *Jd) Trade(orderId string) (tradeRsp defs.TradeRsp, errCode int, err error) {
	var jdTradeRsp payment.TradeRsp
	defer func() {
		//记录日志
		logrus.Info("order id:%v,org:%v,method:%v,trade search request encrypt data:%+v,trade search request decrypt data:%v"+
			",trade search response search encrypt data:%v,trade search response search decrypt data:%v",
			jdTradeRsp.OrderId, "jd", "jd_payment", jdTradeRsp.EncryptRes, jdTradeRsp.DecryptRes, jdTradeRsp.EncryptRsp, jdTradeRsp.DecryptRsp)

	}()

	privateKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString(JdPrivateKey))
	privateFile, err := os.Open(privateKeyPath)
	if err != nil {
		logrus.Errorf(code.PrivateKeyNotExitsErrMessage+",errCode:%v,err:%v", code.PrivateKeyNotExitsErrCode, err.Error())
		return tradeRsp, code.PrivateKeyNotExitsErrCode, errors.New(code.PrivateKeyNotExitsErrMessage)
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
		logrus.Errorf(code.PublicKeyNotExitsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExitsErrCode, err.Error())
		return tradeRsp, code.PublicKeyNotExitsErrCode, errors.New(code.PublicKeyNotExitsErrMessage)
	}
	publicKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf(code.PublicKeyContentErrMessage+",errCode:%v,err:%v", code.PublicKeyContentErrCode, err.Error())
		return tradeRsp, code.PublicKeyContentErrCode, errors.New(code.PublicKeyContentErrMessage)
	}
	publicKey := string(publicKeyBytes)

	tradeArg := payment.TradeArg{
		Merchant:   config.GetInstance().GetString(JdMerchant),
		TradeNum:   orderId,
		DesKey:     config.GetInstance().GetString(JdDesKey),
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		GateWay:    config.GetInstance().GetString(JdTradeWay),
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
			jdClosedRsp.OrderId, "jd", "jd_payment", jdClosedRsp.EncryptRes, jdClosedRsp.DecryptRes, jdClosedRsp.EncryptRsp, jdClosedRsp.DecryptRsp)
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
		logrus.Errorf(code.PrivateKeyNotExitsErrMessage+",errCode:%v,err:%v", code.PrivateKeyNotExitsErrCode, err.Error())
		return closedRsp, code.PrivateKeyNotExitsErrCode, errors.New(code.PrivateKeyNotExitsErrMessage)
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
		logrus.Errorf(code.PublicKeyNotExitsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExitsErrCode, err.Error())
		return closedRsp, code.PublicKeyNotExitsErrCode, errors.New(code.PublicKeyNotExitsErrMessage)
	}
	publicKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf(code.PublicKeyContentErrMessage+",errCode:%v,err:%v", code.PublicKeyContentErrCode, err.Error())
		return closedRsp, code.PublicKeyContentErrCode, errors.New(code.PublicKeyContentErrMessage)
	}
	publicKey := string(publicKeyBytes)

	closedArg := payment.ClosedArg{
		Merchant:   config.GetInstance().GetString(JdMerchant),
		TradeNum:   arg.OrderId + "jd",
		OTradeNum:  arg.OrderId,
		Amount:     amount,
		Currency:   arg.Currency,
		DesKey:     config.GetInstance().GetString(JdDesKey),
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		GateWay:    config.GetInstance().GetString(JdClosedWay),
	}
	jdClosedRsp, errCode, err = new(payment.Closed).Trade(closedArg)
	if err != nil {
		return closedRsp, errCode, err
	}

	closedRsp.OrderId = jdClosedRsp.OrderId
	closedRsp.Status = jdClosedRsp.Status

	return closedRsp, 0, nil
}
