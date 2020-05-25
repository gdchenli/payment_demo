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
		form, errCode, err = pay.allpaySubmit(*order)
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
	jdpayArg := method.JdpayArg{
		OrderId:       order.OrderId,
		TotalFee:      order.TotalFee,
		Currency:      order.Currency,
		UserId:        order.UserId,
		UserAgentType: order.UserAgentType,
	}

	return new(method.Jd).Submit(jdpayArg)
}

func (pay *Pay) allpaySubmit(order defs.Order) (form string, errCode int, err error) {
	allpayArg := method.AllpayArg{
		OrderId:       order.OrderId,
		TotalFee:      order.TotalFee,
		Currency:      order.Currency,
		UserId:        order.UserId,
		UserAgentType: order.UserAgentType,
		MethodCode:    order.MethodCode,
	}

	return new(method.Allpay).Submit(allpayArg)
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
		notifyRsp, errCode, err = pay.allpayNotify(ctx)
	case AlipayOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case WechatOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case EpaymentsOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
	}

	msg := NotifyFailMsg
	if err != nil {
		fmt.Println(errCode)
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(msg))
		return
	}

	if notifyRsp.Status {
		msg = notifyRsp.Message
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

func (pay *Pay) allpayNotify(ctx *gin.Context) (notifyRsp defs.NotifyRsp, errCode int, err error) {
	var notifyBytes []byte
	notifyBytes, err = ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return notifyRsp, errCode, err
	}
	defer func() {
		ctx.Request.Body.Close()
	}()
	query := string(notifyBytes)

	return new(method.Allpay).Notify(query, ctx.Param("method"))
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
		callBackRsp, errCode, err = pay.allpayCallback(ctx)
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
	query := ctx.Request.PostForm.Encode()

	return new(method.Jd).Callback(query)
}

func (pay *Pay) allpayCallback(ctx *gin.Context) (callBackRsp defs.CallbackRsp, errCode int, err error) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()

	return new(method.Allpay).Callback(query, ctx.Param("method"))
}

func (pay *Pay) trade(ctx *gin.Context) {
	var errCode int
	var err error
	var tradeRsp defs.TradeRsp

	trade := new(defs.Trade)
	ctx.ShouldBind(trade)

	if errCode, err = trade.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	org := ctx.Query("org_code")
	switch org {
	case JdOrg:
		tradeRsp, errCode, err = pay.jdTrade(*trade)
	case AllpayOrg:
		tradeRsp, errCode, err = pay.allpayTrade(*trade)
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

func (pay *Pay) allpayTrade(trade defs.Trade) (tradeRsp defs.TradeRsp, errCode int, err error) {
	return new(method.Allpay).Trade(trade.OrderId, trade.MethodCode)
}

func (pay *Pay) closed(ctx *gin.Context) {
	var errCode int
	var err error
	var closedRsp defs.ClosedRsp

	closed := new(defs.Closed)
	ctx.ShouldBind(closed)

	if errCode, err = closed.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	org := ctx.Query("org_code")
	switch org {
	case JdOrg:
		closedRsp, errCode, err = pay.jdClosed(*closed)
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
