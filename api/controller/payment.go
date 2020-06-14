package controller

import (
	"io/ioutil"
	"net/http"
	"payment_demo/api/validate"
	"payment_demo/internal/service/payment"

	"github.com/gin-gonic/gin"
)

type Payment struct{}

func (p *Payment) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/order/pay", p.Pay)              //发起支付
		r.POST("/notify/:org/:method", p.notify) //异步通知
		r.POST("/verify/:org/:method", p.verify) //同步通知
		r.GET("/trade/search", p.searchTrade)    //交易查询
		r.GET("/trade/close", p.closeTrade)      //关闭交易
		r.POST("/logistics/upload", p.upload)    //上传物流信息
	}
}

func (p *Payment) upload(ctx *gin.Context) {
	uploadLogisticsReq := new(validate.UploadLogisticsReq)
	ctx.ShouldBind(uploadLogisticsReq)

	logisticsRsp, errCode, err := new(payment.Payment).UploadLogistics(*uploadLogisticsReq)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": logisticsRsp})
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

	notifyRsp, errCode, err := new(payment.Payment).Notify(query, ctx.Param("org"), ctx.Param("method"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": notifyRsp})
}

func (p *Payment) searchTrade(ctx *gin.Context) {
	searchTradeReq := new(validate.SearchTradeReq)
	ctx.ShouldBind(searchTradeReq)

	if errCode, err := searchTradeReq.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	searchTradeRsp, errCode, err := new(payment.Payment).SearchTrade(*searchTradeReq)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": searchTradeRsp})
}

func (p *Payment) closeTrade(ctx *gin.Context) {
	closeTradeReq := new(validate.CloseTradeReq)
	ctx.ShouldBind(closeTradeReq)

	if errCode, err := closeTradeReq.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	closeTradeRsp, errCode, err := new(payment.Payment).CloseTrade(*closeTradeReq)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, closeTradeRsp)
}

//发起支付
func (p *Payment) Pay(ctx *gin.Context) {
	o := new(validate.Order)
	ctx.ShouldBind(o)

	if errCode, err := o.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	submitRsp, errCode, err := new(payment.Payment).Sumbit(*o)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": submitRsp})
}

func (p *Payment) verify(ctx *gin.Context) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()
	if query == "" {
		query = ctx.Request.URL.Query().Encode()
	}

	verifyRsp, errCode, err := new(payment.Payment).Verify(query, ctx.Param("org"), ctx.Param("method"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, verifyRsp)
}
