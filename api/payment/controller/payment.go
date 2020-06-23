package controller

import (
	"net/http"
	"payment_demo/api/payment/request"
	"payment_demo/internal/service/payment"

	"github.com/gin-gonic/gin"
)

type Payment struct{}

func (p *Payment) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/pay", p.pay)                   //h5 、pc发起支付
		r.POST("/alipayminiprogram/pay", p.pay) //支付宝小程序发起支付
		r.GET("/qrcodeimg", p.qrcode)           //二维码支付
		r.GET("/form", p.form)                  //发起支付
	}
}

func (p *Payment) pay(ctx *gin.Context) {
	arg := new(request.OrderArg)
	ctx.ShouldBind(arg)

	if errCode, err := arg.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	paymentService, errCode, err := payment.New(arg.OrgCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	payRsp, errCode, err := paymentService.Pay(*arg)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": payRsp})
}

func (p *Payment) qrcode(ctx *gin.Context) {
	arg := new(request.OrderArg)
	ctx.ShouldBind(arg)

	if errCode, err := arg.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	paymentService, errCode, err := payment.New(arg.OrgCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}
	submitRsp, errCode, err := paymentService.PayQrCode(*arg)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": submitRsp})
}

func (p *Payment) form(ctx *gin.Context) {
	arg := new(request.OrderArg)
	ctx.ShouldBind(arg)

	if errCode, err := arg.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	paymentService, errCode, err := payment.New(arg.OrgCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}
	submitRsp, errCode, err := paymentService.PayForm(*arg)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": submitRsp})
}
