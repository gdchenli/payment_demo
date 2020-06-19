package controller

import (
	"io/ioutil"
	"net/http"
	"payment_demo/api/request"
	"payment_demo/internal/service/payment"

	"github.com/gin-gonic/gin"
)

type Payment struct{}

func (p *Payment) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/pay", p.pay)                    //h5 、pc发起支付
		r.POST("/alipayminiprogram/pay", p.pay)  //支付宝小程序发起支付
		r.GET("/qrcodeimg", p.qrcode)            //二维码支付
		r.GET("/form", p.form)                   //发起支付
		r.POST("/notify/:org/:method", p.notify) //异步通知
		r.POST("/verify/:org/:method", p.verify) //同步通知
		r.GET("/trade/search", p.searchTrade)    //交易查询
		r.GET("/trade/close", p.closeTrade)      //关闭交易
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

func (p *Payment) notify(ctx *gin.Context) {
	notifyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": "1", "message": err.Error()})
		return
	}
	defer func() {
		ctx.Request.Body.Close()
	}()
	query := string(notifyBytes)

	orgCode := ctx.Param("org")
	paymentService, errCode, err := payment.New(orgCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}
	notifyRsp, errCode, err := paymentService.Notify(query, orgCode, ctx.Param("method"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": notifyRsp})
}

func (p *Payment) verify(ctx *gin.Context) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()
	if query == "" {
		query = ctx.Request.URL.Query().Encode()
	}

	orgCode := ctx.Param("org")
	paymentService, errCode, err := payment.New(orgCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}
	verifyRsp, errCode, err := paymentService.Verify(query, orgCode, ctx.Param("method"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, verifyRsp)
}

func (p *Payment) searchTrade(ctx *gin.Context) {
	arg := new(request.SearchTradeArg)
	ctx.ShouldBind(arg)

	if errCode, err := arg.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	paymentService, errCode, err := payment.New(arg.OrgCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}
	searchTradeRsp, errCode, err := paymentService.SearchTrade(*arg)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": searchTradeRsp})
}

func (p *Payment) closeTrade(ctx *gin.Context) {
	arg := new(request.CloseTradeArg)
	ctx.ShouldBind(arg)

	if errCode, err := arg.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	paymentService, errCode, err := payment.New(arg.OrgCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}
	closeTradeRsp, errCode, err := paymentService.CloseTrade(*arg)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, closeTradeRsp)
}
