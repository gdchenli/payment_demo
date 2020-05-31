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
	return new(method.Jd).Notify(query, code.JdMethod)
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

	return new(method.Allpay).Notify(query, ctx.Param("method"))
}

func (notify *Notify) alipayNotify(ctx *gin.Context) (notifyRsp defs.NotifyRsp, errCode int, err error) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()

	return new(method.Alipay).Notify(query, ctx.Param("method"))
}
