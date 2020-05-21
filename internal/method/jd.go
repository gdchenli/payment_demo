package method

import (
	"errors"
	"fmt"
	"payment_demo/internal/common/config"
	"payment_demo/pkg/jd/payment"
	"strconv"
	"time"
)

type Jd struct{}

type JdPayArg struct {
	OrderId  string  `json:"order_id"`
	TotalFee float64 `json:"total_fee"`
	Currency string  `json:"currency"`
	UserId   string  `json:"user_id"`
}

func (jd *Jd) Submit(arg JdPayArg) (form string, errCode int, err error) {
	//金额转为分
	totalFee := arg.TotalFee * 100
	//金额字段类型转换
	amount, err := strconv.ParseInt(fmt.Sprintf("%.f", totalFee), 10, 64)
	if err != nil {
		return form, 10101, errors.New("金额转换异常")
	}

	date := time.Now().Format(payment.TimeLayout)
	goodsInfos := []payment.GoodsInfo{{Id: "test" + date, Name: "test" + date, Price: amount, Num: 1}}
	kjInfo := payment.KjInfo{GoodsSubmittedCustoms: "N", GoodsUnderBonded: "N"}
	payArg := payment.PayArg{
		Merchant:       config.GetInstance().GetString("jd.merchant"),
		TradeNum:       arg.OrderId,
		Amount:         amount,
		Currency:       arg.Currency,
		CallbackUrl:    config.GetInstance().GetString("jd.callback_url"),
		NotifyUrl:      config.GetInstance().GetString("jd.notify_url"),
		UserId:         arg.UserId,
		ExpireTime:     config.GetInstance().GetString("jd.expire_time"),
		SettleCurrency: config.GetInstance().GetString("jd.settle_currency"),
		GoodsInfo:      goodsInfos,
		KjInfo:         kjInfo,
	}
	form, errCode, err = new(payment.Payment).CreatePaymentForm(payArg)
	if err != nil {
		return form, errCode, err
	}

	return form, 0, nil
}
