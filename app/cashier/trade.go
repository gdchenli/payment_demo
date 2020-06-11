package cashier

import (
	"net/http"
	"payment_demo/internal/cashier"
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
	var errCode int
	var err error
	var tradeRsp defs.TradeRsp

	t := new(defs.Trade)
	ctx.ShouldBind(t)

	if errCode, err = t.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	org := ctx.Query("org_code")
	payMethod := getPayMethod(org)
	if payMethod == nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(NotSupportPaymentOrgMsg))
		return
	}

	tradeRsp, errCode, err = payMethod.Trade(t.OrderId, t.MethodCode, t.Currency, t.TotalFee)
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

	org := ctx.Query("org_code")

	payMethod := getPayMethod(org)
	if payMethod == nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(NotSupportPaymentOrgMsg))
		return
	}

	closedRsp, errCode, err = payMethod.Closed(*closed)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, closedRsp)
}

func (trade *Trade) jdClosed(closed defs.Closed) (closedRsp defs.ClosedRsp, errCode int, err error) {
	return new(cashier.Jd).Closed(closed)
}
