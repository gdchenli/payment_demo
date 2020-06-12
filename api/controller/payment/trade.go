package payment

import (
	"net/http"
	"payment_demo/api/controller/common"
	"payment_demo/api/validate"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Trade struct{}

func (trade *Trade) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.GET("/trade/search", trade.search) //交易查询
		r.GET("/trade/close", trade.close)   //关闭交易
	}
}

func (trade *Trade) search(ctx *gin.Context) {
	t := new(validate.SearchTradeReq)
	ctx.ShouldBind(t)

	if errCode, err := t.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	searchTradeHandle := common.GetSearchTradeHandler(ctx.Query("org_code"))
	if searchTradeHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": common.NotSupportPaymentOrgCode, "message": common.NotSupportPaymentOrgMsg})
		return
	}

	searchTradeRsp, errCode, err := searchTradeHandle(t.OrderId, t.MethodCode, t.Currency, t.TotalFee)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, searchTradeRsp)
}

func (trade *Trade) close(ctx *gin.Context) {
	closeTradeReq := new(validate.CloseTradeReq)
	ctx.ShouldBind(closeTradeReq)

	if errCode, err := closeTradeReq.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	closeTradeHandle := common.GetCloseTradeHandler(ctx.Query("org_cod"))
	if closeTradeHandle == nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(common.NotSupportPaymentOrgMsg))
		return
	}

	closeTradeRsp, errCode, err := closeTradeHandle(*closeTradeReq)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, closeTradeRsp)
}
