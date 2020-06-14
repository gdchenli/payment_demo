package controller

import (
	"io/ioutil"
	"net/http"
	"payment_demo/api/validate"
	"payment_demo/internal/service"

	"github.com/gin-gonic/gin"
)

type Payment struct{}

func (payment *Payment) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/order/pay", payment.Pay)              //发起支付
		r.POST("/notify/:org/:method", payment.notify) //异步通知
		r.POST("/verify/:org/:method", payment.verify) //同步通知
		r.GET("/trade/search", payment.searchTrade)    //交易查询
		r.GET("/trade/close", payment.closeTrade)      //关闭交易
		r.POST("/logistics/upload", payment.upload)    //上传物流信息
	}
}

func (payment *Payment) upload(ctx *gin.Context) {
	uploadLogisticsReq := new(validate.UploadLogisticsReq)
	ctx.ShouldBind(uploadLogisticsReq)

	logisticsRsp, errCode, err := new(service.Payment).UploadLogistics(*uploadLogisticsReq)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": logisticsRsp})
}

func (payment *Payment) notify(ctx *gin.Context) {
	notifyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": "1", "message": err.Error()})
		return
	}
	defer func() {
		ctx.Request.Body.Close()
	}()
	query := string(notifyBytes)

	notifyRsp, errCode, err := new(service.Payment).Notify(query, ctx.Param("org"), ctx.Param("method"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": notifyRsp})
}

func (payment *Payment) searchTrade(ctx *gin.Context) {
	searchTradeReq := new(validate.SearchTradeReq)
	ctx.ShouldBind(searchTradeReq)

	if errCode, err := searchTradeReq.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	searchTradeRsp, errCode, err := new(service.Payment).SearchTrade(*searchTradeReq)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": searchTradeRsp})
}

func (payment *Payment) closeTrade(ctx *gin.Context) {
	closeTradeReq := new(validate.CloseTradeReq)
	ctx.ShouldBind(closeTradeReq)

	if errCode, err := closeTradeReq.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	closeTradeRsp, errCode, err := new(service.Payment).CloseTrade(*closeTradeReq)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, closeTradeRsp)
}

//发起支付
func (payment *Payment) Pay(ctx *gin.Context) {
	o := new(validate.Order)
	ctx.ShouldBind(o)

	if errCode, err := o.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	submitRsp, errCode, err := new(service.Payment).Sumbit(*o)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": submitRsp})
}

func (payment *Payment) verify(ctx *gin.Context) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()
	if query == "" {
		query = ctx.Request.URL.Query().Encode()
	}

	verifyRsp, errCode, err := new(service.Payment).Verify(query, ctx.Param("org"), ctx.Param("method"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, verifyRsp)
}
