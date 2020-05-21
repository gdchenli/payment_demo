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

	privateKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString("jd.private_key"))
	file, err := os.Open(privateKeyPath)
	if err != nil {
		logrus.Errorf(code.PrivateKeyNotExitsErrMessage+",errCode:%v,err:%v", code.PrivateKeyNotExitsErrCode, err.Error())
		return form, code.PrivateKeyNotExitsErrCode, errors.New(code.PrivateKeyNotExitsErrMessage)
	}
	privateKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf(code.PrivateKeyContentErrMessage+",errCode:%v,err:%v", code.PrivateKeyContentErrCode, err.Error())
		return form, code.PrivateKeyContentErrCode, errors.New(code.PrivateKeyContentErrMessage)
	}
	privateKey := string(privateKeyBytes)

	date := time.Now().Format(payment.TimeLayout)
	goodsInfos := []payment.GoodsInfo{{Id: "test" + date, Name: "test" + date, Price: amount, Num: 1}}
	kjInfo := payment.KjInfo{GoodsSubmittedCustoms: "N", GoodsUnderBonded: "N"}
	payArg := payment.PayArg{
		Merchant:      config.GetInstance().GetString("jd.merchant"),
		TradeNum:      arg.OrderId,
		Amount:        amount,
		Currency:      arg.Currency,
		CallbackUrl:   config.GetInstance().GetString("jd.callback_url"),
		NotifyUrl:     config.GetInstance().GetString("jd.notify_url"),
		UserId:        arg.UserId,
		ExpireTime:    config.GetInstance().GetString("jd.expire_time"),
		TransCurrency: config.GetInstance().GetString("jd.settle_currency"),
		GoodsInfo:     goodsInfos,
		KjInfo:        kjInfo,
		DesKey:        config.GetInstance().GetString("jd.des_key"),
		PrivateKey:    privateKey,
		PayWay:        config.GetInstance().GetString("jd.pc_way"),
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
	publicKeyPath := path.Join(config.GetInstance().GetString("app_path"), config.GetInstance().GetString("jd.public_key"))
	file, err := os.Open(publicKeyPath)
	if err != nil {
		logrus.Errorf(code.PublicKeyNotExitsErrMessage+",errCode:%v,err:%v", code.PublicKeyNotExitsErrCode, err.Error())
		return notifyRsp, code.PublicKeyNotExitsErrCode, errors.New(code.PublicKeyNotExitsErrMessage)
	}
	publicKeyBytes, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf(code.PublicKeyContentErrMessage+",errCode:%v,err:%v", code.PublicKeyContentErrCode, err.Error())
		return notifyRsp, code.PublicKeyContentErrCode, errors.New(code.PublicKeyContentErrMessage)
	}
	publicKey := string(publicKeyBytes)
	notifyArg := payment.NotifyArg{
		PublicKey: publicKey,
		DesKey:    config.GetInstance().GetString("jd.des_key"),
	}

	jdNotifyRsp, errCode, err := new(payment.Notify).Validate(query, notifyArg)
	if err != nil {
		return notifyRsp, errCode, err
	}

	notifyRsp.OrderId = jdNotifyRsp.OrderId
	notifyRsp.Status = jdNotifyRsp.Status

	return notifyRsp, 0, nil
}
