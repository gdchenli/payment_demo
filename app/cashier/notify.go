package cashier

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"payment_demo/internal/cashier"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/defs"

	"github.com/gin-gonic/gin"
)

type Notify struct{}

func (notify *Notify) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/notify/:org/:method", notify.notify) //异步通知
	}
}

func (notify *Notify) notify(ctx *gin.Context) {
	var errCode int
	var err error
	var notifyRsp defs.NotifyRsp

	org := ctx.Param("org")
	switch org {
	case JdOrg:
		notifyRsp, errCode, err = notify.jdNotify(ctx)
	case AllpayOrg:
		notifyRsp, errCode, err = notify.allpayNotify(ctx)
	case AlipayOrg:
		notifyRsp, errCode, err = notify.alipayNotify(ctx)
	case EpaymentsOrg:
		notifyRsp, errCode, err = notify.epaymentsNotify(ctx)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
	}

	if err != nil {
		fmt.Println(errCode)
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": notifyRsp})
}

func (notify *Notify) jdNotify(ctx *gin.Context) (notifyRsp defs.NotifyRsp, errCode int, err error) {
	var notifyBytes []byte
	notifyBytes, err = ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return notifyRsp, errCode, err
	}
	defer func() {
		ctx.Request.Body.Close()
	}()
	query := string(notifyBytes)
	return new(cashier.Jd).Notify(query, code.JdMethod)
}

func (notify *Notify) allpayNotify(ctx *gin.Context) (notifyRsp defs.NotifyRsp, errCode int, err error) {
	var notifyBytes []byte
	notifyBytes, err = ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return notifyRsp, errCode, err
	}
	defer func() {
		ctx.Request.Body.Close()
	}()
	query := string(notifyBytes)

	return new(cashier.Allpay).Notify(query, ctx.Param("method"))
}

func (notify *Notify) alipayNotify(ctx *gin.Context) (notifyRsp defs.NotifyRsp, errCode int, err error) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()

	return new(cashier.Alipay).Notify(query, ctx.Param("method"))
}

func (notify *Notify) epaymentsNotify(ctx *gin.Context) (notifyRsp defs.NotifyRsp, errCode int, err error) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()

	return new(cashier.Epayments).Notify(query, ctx.Param("method"))
}
