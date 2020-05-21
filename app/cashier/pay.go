package cashier

import (
	"azoya/nova/binding"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"payment_demo/internal/common/defs"
	"payment_demo/internal/method"

	"github.com/gin-gonic/gin"
)

const (
	NotSupportPaymentOrg = "不支持该支付机构"
	JdOrg                = "jd"
	AllpayOrg            = "allpay"
	WechatOrg            = "wechat"
	EpaymentsOrg         = "epayments"
	AlipayOrg            = "alipay"
	NotifySuccessMsg     = "success"
	NotifyFailMsg        = "fail"
)

type Pay struct{}

func (pay *Pay) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/submit", pay.submit)
		r.POST("/notify/:org/:method", pay.notify)
		r.POST("/verify/:org/:method", pay.verify)
		r.GET("/status", pay.status)
	}
}

//发起支付
func (pay *Pay) submit(ctx *gin.Context) {
	var errCode int
	var err error
	var form string
	order := new(defs.Order)
	ctx.ShouldBind(order)

	if errCode, err = order.Validate(); err != nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(err.Error()))
		return
	}

	switch order.OrgCode {
	case JdOrg:
		form, errCode, err = pay.jdSubmit(*order)
	case AllpayOrg:
		err = errors.New(NotSupportPaymentOrg)
	case AlipayOrg:
		err = errors.New(NotSupportPaymentOrg)
	case WechatOrg:
		err = errors.New(NotSupportPaymentOrg)
	case EpaymentsOrg:
		err = errors.New(NotSupportPaymentOrg)
	default:
		err = errors.New(NotSupportPaymentOrg)
	}

	if err != nil {
		fmt.Println(errCode)
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(err.Error()))
		return
	}

	ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(form))
}

func (pay *Pay) jdSubmit(order defs.Order) (form string, errCode int, err error) {
	jdPayArg := method.JdPayArg{
		OrderId:  order.OrderId,
		TotalFee: order.TotalFee,
		Currency: order.Currency,
		UserId:   order.UserId,
	}
	return new(method.Jd).Submit(jdPayArg)
}

func (pay *Pay) notify(ctx *gin.Context) {
	var errCode int
	var err error
	var notifyRsp defs.NotifyRsp

	org := ctx.Param("org")
	switch org {
	case JdOrg:
		notifyRsp, errCode, err = pay.jdNotify(ctx)
	case AllpayOrg:
		err = errors.New(NotSupportPaymentOrg)
	case AlipayOrg:
		err = errors.New(NotSupportPaymentOrg)
	case WechatOrg:
		err = errors.New(NotSupportPaymentOrg)
	case EpaymentsOrg:
		err = errors.New(NotSupportPaymentOrg)
	default:
		err = errors.New(NotSupportPaymentOrg)
	}

	msg := NotifySuccessMsg
	if err != nil || !notifyRsp.Status {
		fmt.Println(errCode)
		msg = NotifyFailMsg
	}

	ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(msg))
}

func (pay *Pay) jdNotify(ctx *gin.Context) (notifyRsp defs.NotifyRsp, errCode int, err error) {
	var notifyBytes []byte
	notifyBytes, err = ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return notifyRsp, errCode, err
	}
	defer func() {
		ctx.Request.Body.Close()
	}()
	query := string(notifyBytes)
	return new(method.Jd).Notify(query)
}

func (pay *Pay) verify(ctx *gin.Context) {

}

func (pay *Pay) status(ctx *gin.Context) {

}
