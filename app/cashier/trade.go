package cashier

import (
	"net/http"
	"payment_demo/internal/common/defs"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Trade struct{}

func (trade *Trade) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.GET("/trade", trade.search)  //交易查询
		r.GET("/closed", trade.closed) //关闭交易
	}
}

func (trade *Trade) search(ctx *gin.Context) {
	t := new(defs.Trade)
	ctx.ShouldBind(t)

	if errCode, err := t.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	tradeHandle := getTradeHandler(ctx.Query("org_code"))
	if tradeHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": NotSupportPaymentOrgCode, "message": NotSupportPaymentOrgMsg})
		return
	}

	tradeRsp, errCode, err := tradeHandle(t.OrderId, t.MethodCode, t.Currency, t.TotalFee)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, tradeRsp)
}

func (trade *Trade) closed(ctx *gin.Context) {
	var errCode int
	var err error
	var closedRsp defs.ClosedRsp

	closed := new(defs.Closed)
	ctx.ShouldBind(closed)

	if errCode, err = closed.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	closeHandle := getClosedHandler(ctx.Query("org_code"))
	if closeHandle == nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(NotSupportPaymentOrgMsg))
		return
	}

	closedRsp, errCode, err = closeHandle(*closed)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, closedRsp)
}
