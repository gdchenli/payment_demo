package controller

import (
	"io/ioutil"
	"net/http"
	"payment_demo/api/validate"

	"github.com/gin-gonic/gin/binding"

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
	l := new(validate.UploadLogisticsReq)
	ctx.ShouldBind(l)

	uploadLogisticsHandle := GetUploadLogisticsHandler(l.OrgCode)
	if uploadLogisticsHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": NotSupportPaymentOrgCode, "message": NotSupportPaymentOrgMsg})
		return
	}

	logisticsRsp, errCode, err := uploadLogisticsHandle(*l)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, logisticsRsp)
}

func (payment *Payment) notify(ctx *gin.Context) {
	notifyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": NotSupportPaymentOrgCode, "message": err.Error()})
		return
	}
	defer func() {
		ctx.Request.Body.Close()
	}()
	query := string(notifyBytes)

	notifyHandle := GetNotifyHandler(ctx.Param("org"))
	if notifyHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": NotSupportPaymentOrgCode, "message": NotSupportPaymentOrgMsg})
		return
	}

	notifyRsp, errCode, err := notifyHandle(query, ctx.Param("method"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": notifyRsp})
}

func (payment *Payment) searchTrade(ctx *gin.Context) {
	t := new(validate.SearchTradeReq)
	ctx.ShouldBind(t)

	if errCode, err := t.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	searchTradeHandle := GetSearchTradeHandler(ctx.Query("org_code"))
	if searchTradeHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": NotSupportPaymentOrgCode, "message": NotSupportPaymentOrgMsg})
		return
	}

	searchTradeRsp, errCode, err := searchTradeHandle(t.OrderId, t.MethodCode, t.Currency, t.TotalFee)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, searchTradeRsp)
}

func (payment *Payment) closeTrade(ctx *gin.Context) {
	closeTradeReq := new(validate.CloseTradeReq)
	ctx.ShouldBind(closeTradeReq)

	if errCode, err := closeTradeReq.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	closeTradeHandle := GetCloseTradeHandler(ctx.Query("org_cod"))
	if closeTradeHandle == nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(NotSupportPaymentOrgMsg))
		return
	}

	closeTradeRsp, errCode, err := closeTradeHandle(*closeTradeReq)
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

	payHandle := GetPayHandler(o.OrgCode)
	if payHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": NotSupportPaymentOrgCode, "message": NotSupportPaymentOrgMsg})
		return
	}

	form, errCode, err := payHandle(*o)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": form})
}

func (payment *Payment) verify(ctx *gin.Context) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()
	if query == "" {
		query = ctx.Request.URL.Query().Encode()
	}

	verifyHandle := GetVerifyHandler(ctx.Param("org"))
	if verifyHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": NotSupportPaymentOrgCode, "message": NotSupportPaymentOrgMsg})
		return
	}

	callBackRsp, errCode, err := verifyHandle(query, ctx.Param("method"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, callBackRsp)
}
