package cashier

import (
	"errors"
	"net/http"
	"payment_demo/internal/cashier"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/defs"

	"github.com/gin-gonic/gin"
)

type Callback struct{}

func (callback *Callback) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/callback/:org/:method", callback.callback) //同步通知
	}
}

func (callback *Callback) callback(ctx *gin.Context) {
	var errCode int
	var err error
	var callBackRsp defs.CallbackRsp

	org := ctx.Param("org")
	switch org {
	case JdOrg:
		callBackRsp, errCode, err = callback.jdCallback(ctx)
	case AllpayOrg:
		callBackRsp, errCode, err = callback.allpayCallback(ctx)
	case AlipayOrg:
		callBackRsp, errCode, err = callback.alipayCallback(ctx)
	case EpaymentsOrg:
		callBackRsp, errCode, err = callback.epaymentsCallback(ctx)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, callBackRsp)
}

func (callback *Callback) jdCallback(ctx *gin.Context) (callBackRsp defs.CallbackRsp, errCode int, err error) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()

	return new(cashier.Jd).Callback(query, code.JdMethod)
}

func (callback *Callback) allpayCallback(ctx *gin.Context) (callBackRsp defs.CallbackRsp, errCode int, err error) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()
	if query == "" {
		query = ctx.Request.URL.Query().Encode()
	}

	return new(cashier.Allpay).Callback(query, ctx.Param("method"))
}

func (callback *Callback) alipayCallback(ctx *gin.Context) (callBackRsp defs.CallbackRsp, errCode int, err error) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()
	if query == "" {
		query = ctx.Request.URL.Query().Encode()
	}

	return new(cashier.Alipay).Callback(query, ctx.Param("method"))
}

func (callback *Callback) epaymentsCallback(ctx *gin.Context) (callBackRsp defs.CallbackRsp, errCode int, err error) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()
	if query == "" {
		query = ctx.Request.URL.Query().Encode()
	}

	return new(cashier.Epayments).Callback(query, ctx.Param("method"))
}
