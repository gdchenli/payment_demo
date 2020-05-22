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
	JdOrg        = "jd"
	AllpayOrg    = "allpay"
	WechatOrg    = "wechat"
	EpaymentsOrg = "epayments"
	AlipayOrg    = "alipay"
)

const (
	NotSupportPaymentOrgMsg = "不支持该支付机构"
	NotifySuccessMsg        = "success"
	NotifyFailMsg           = "fail"
)

type Pay struct{}

func (pay *Pay) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/submit", pay.submit)
		r.POST("/notify/:org/:method", pay.notify)
		r.POST("/callback/:org/:method", pay.callback)
		r.GET("/trade", pay.trade)
		r.GET("/closed", pay.closed)
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
		err = errors.New(NotSupportPaymentOrgMsg)
	case AlipayOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case WechatOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case EpaymentsOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
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
		err = errors.New(NotSupportPaymentOrgMsg)
	case AlipayOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case WechatOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case EpaymentsOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
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

func (pay *Pay) callback(ctx *gin.Context) {
	var errCode int
	var err error
	var callBackRsp defs.CallbackRsp

	org := ctx.Param("org")
	switch org {
	case JdOrg:
		callBackRsp, errCode, err = pay.jdCallback(ctx)
	case AllpayOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case AlipayOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case WechatOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case EpaymentsOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, callBackRsp)
}

func (pay *Pay) jdCallback(ctx *gin.Context) (callBackRsp defs.CallbackRsp, errCode int, err error) {
	ctx.Request.ParseForm()
	//sign := base64.StdEncoding.EncodeToString([]byte(ctx.Request.PostForm.Get("sign")))
	//ctx.Request.PostForm.Set("sign", sign)
	query := ctx.Request.PostForm.Encode()

	return new(method.Jd).Callback(query)
}

func (pay *Pay) trade(ctx *gin.Context) {
	var errCode int
	var err error
	var trade defs.Trade
	var tradeRsp defs.TradeRsp

	ctx.ShouldBind(trade)

	if errCode, err = trade.Validate(); err != nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(err.Error()))
		return
	}

	org := ctx.Param("org")
	switch org {
	case JdOrg:
		tradeRsp, errCode, err = pay.jdTrade(trade)
	case AllpayOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case AlipayOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case WechatOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case EpaymentsOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, tradeRsp)
}

func (pay *Pay) jdTrade(trade defs.Trade) (tradeRsp defs.TradeRsp, errCode int, err error) {
	return new(method.Jd).Trade(trade.OrderId)
}

func (pay *Pay) closed(ctx *gin.Context) {
	var errCode int
	var err error
	var closed defs.Closed
	var closedRsp defs.ClosedRsp

	ctx.ShouldBind(closed)

	if errCode, err = closed.Validate(); err != nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(err.Error()))
		return
	}

	org := ctx.Param("org")
	switch org {
	case JdOrg:
		closedRsp, errCode, err = pay.jdClosed(closed)
	case AllpayOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case AlipayOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case WechatOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case EpaymentsOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, closedRsp)
}

func (pay *Pay) jdClosed(closed defs.Closed) (closedRsp defs.ClosedRsp, errCode int, err error) {
	closedArg := method.JdClosedArg{
		OrderId:  closed.OrderId,
		Currency: closed.Currency,
		TotalFee: closed.TotalFee,
	}
	return new(method.Jd).Closed(closedArg)
}
