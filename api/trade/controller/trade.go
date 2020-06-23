package controller

import (
	"net/http"
	request2 "payment_demo/api/trade/request"
	"payment_demo/internal/service/trade"

	"github.com/gin-gonic/gin"
)

type Trade struct{}

func (t *Trade) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.GET("/trade/search", t.searchTrade) //交易查询
		r.GET("/trade/close", t.closeTrade)   //关闭交易
	}
}
func (t *Trade) searchTrade(ctx *gin.Context) {
	arg := new(request2.SearchTradeArg)
	ctx.ShouldBind(arg)

	if errCode, err := arg.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	tradeService, errCode, err := trade.New(arg.OrgCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}
	searchTradeRsp, errCode, err := tradeService.SearchTrade(*arg)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": searchTradeRsp})
}

func (t *Trade) closeTrade(ctx *gin.Context) {
	arg := new(request2.CloseTradeArg)
	ctx.ShouldBind(arg)

	if errCode, err := arg.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	tradeService, errCode, err := trade.New(arg.OrgCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}
	closeTradeRsp, errCode, err := tradeService.CloseTrade(*arg)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, closeTradeRsp)
}
