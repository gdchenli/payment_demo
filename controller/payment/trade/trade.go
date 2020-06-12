package trade

import (
	"net/http"
	"payment_demo/controller/common"
	"payment_demo/internal/common/defs"

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
	t := new(SearchReq)
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
	close := new(CloseReq)
	ctx.ShouldBind(close)

	if errCode, err := close.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	closeTradeHandle := common.GetCloseTradeHandler(ctx.Query("org_cod"))
	if closeTradeHandle == nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(common.NotSupportPaymentOrgMsg))
		return
	}

	closeTradeReq := defs.CloseTradeReq{
		OrderId:    close.OrderId,
		TotalFee:   close.TotalFee,
		Currency:   close.Currency,
		MethodCode: close.MethodCode,
		OrgCode:    close.OrgCode,
	}
	closeTradeRsp, errCode, err := closeTradeHandle(closeTradeReq)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, closeTradeRsp)
}
