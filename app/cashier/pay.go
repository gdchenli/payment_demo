package cashier

import (
	"azoya/nova/binding"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"payment_demo/internal/common/code"
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
		r.POST("/order/submit", pay.orderSubmit)         //发起支付
		r.POST("/amp/submit", pay.ampSubmit)             //支付宝小程序，支付参数
		r.POST("/notify/:org/:method", pay.notify)       //异步通知
		r.POST("/callback/:org/:method", pay.callback)   //同步通知
		r.GET("/trade", pay.trade)                       //交易查询
		r.GET("/closed", pay.closed)                     //关闭交易
		r.POST("/logistics/upload", pay.logisticsUpload) //物流信息回传
	}
}

//发起支付
func (pay *Pay) ampSubmit(ctx *gin.Context) {
	var errCode int
	var err error
	var payStr string
	order := new(defs.Order)
	ctx.ShouldBind(order)
	order.UserAgentType = 7

	if errCode, err = order.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}
	switch order.OrgCode {
	case JdOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case AllpayOrg:
		payStr, errCode, err = pay.allpayAmpSubmit(*order)
	case AlipayOrg:
		payStr, errCode, err = pay.alipayAmpSubmit(*order)
	case WechatOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case EpaymentsOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
	}

	if errCode, err = order.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": payStr})

}

func (pay *Pay) allpayAmpSubmit(order defs.Order) (form string, errCode int, err error) {
	return new(method.Allpay).AmpSubmit(order)
}

func (pay *Pay) alipayAmpSubmit(order defs.Order) (form string, errCode int, err error) {
	return new(method.Alipay).AmpSubmit(order)
}

//发起支付
func (pay *Pay) orderSubmit(ctx *gin.Context) {
	var errCode int
	var err error
	var form string
	order := new(defs.Order)
	ctx.ShouldBind(order)

	if errCode, err = order.Validate(); err != nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(err.Error()))
		return
	}

	var payMethod method.PayMethod
	switch order.OrgCode {
	case JdOrg:
		payMethod = new(method.Jd)
	case AllpayOrg:
		payMethod = new(method.Allpay)
	case AlipayOrg:
		payMethod = new(method.Alipay)
	case EpaymentsOrg:
		payMethod = new(method.Epayments)
	default:
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(err.Error()))
		return
	}

	form, errCode, err = payMethod.OrderSubmit(*order)
	if err != nil {
		fmt.Println(errCode)
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(err.Error()))
		return
	}

	ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(form))
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
	return new(method.Jd).Notify(query, code.JdMethod)
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

	return new(method.Jd).Callback(query, code.JdMethod)
}

func (pay *Pay) allpayCallback(ctx *gin.Context) (callBackRsp defs.CallbackRsp, errCode int, err error) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()
	if query == "" {
		query = ctx.Request.URL.Query().Encode()
	}

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
	return new(method.Jd).Trade(trade.OrderId, code.JdMethod)
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

func (pay *Pay) logisticsUpload(ctx *gin.Context) {
	var errCode int
	var err error
	var logisticsRsp defs.LogisticsRsp

	logistics := new(defs.Logistics)
	ctx.ShouldBind(logistics)

	switch logistics.OrgCode {
	case JdOrg:
		logisticsRsp, errCode, err = pay.JdLogisticsUpload(logistics)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, logisticsRsp)
}

func (pay *Pay) JdLogisticsUpload(logistics *defs.Logistics) (logisticsRsp defs.LogisticsRsp, errCode int, err error) {
	arg := method.JdLogisticsArg{
		OrderId:          logistics.OrderId,
		LogisticsCompany: logistics.LogisticsCompany,
		LogisticsNo:      logistics.LogisticsNo,
	}

	return new(method.Jd).LogisticsUpload(arg)
}
