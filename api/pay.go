package api

import (
	"net/http"
	"payment_demo/api/validate"
	"payment_demo/internal/service"

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

type OrderArg struct {
	OrderId       string  `form:"order_id" json:"order_id"`               //订单编号
	TotalFee      float64 `form:"total_fee" json:"total_fee"`             //金额
	Currency      string  `form:"currency" json:"currency"`               //币种
	MethodCode    string  `form:"method_code" json:"method_code"`         //支付方式
	OrgCode       string  `form:"org_code" json:"org_code"`               //支付机构
	UserId        string  `form:"user_id" json:"user_id"`                 //用户Id
	UserAgentType int     `form:"user_agent_type" json:"user_agent_type"` //环境
}

func (p *Payment) pay(ctx *gin.Context) {
	arg := new(validate.OrderArg)
	ctx.ShouldBind(arg)

	if errCode, err := arg.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	paymentService, errCode, err := service.NewPay(arg.OrgCode)
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
	arg := new(validate.OrderArg)
	ctx.ShouldBind(arg)

	if errCode, err := arg.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	paymentService, errCode, err := service.NewPay(arg.OrgCode)
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
	arg := new(validate.OrderArg)
	ctx.ShouldBind(arg)

	if errCode, err := arg.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	paymentService, errCode, err := service.NewPay(arg.OrgCode)
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
