package cashier

import (
	"errors"
	"net/http"
	"payment_demo/internal/common/code"
	"payment_demo/internal/common/defs"
	"payment_demo/internal/method"

	"github.com/gin-gonic/gin"
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
	switch org {
	case JdOrg:
		tradeRsp, errCode, err = trade.jdTrade(*t)
	case AllpayOrg:
		tradeRsp, errCode, err = trade.allpayTrade(*t)
	case AlipayOrg:
		tradeRsp, errCode, err = trade.alipayTrade(*t)
	case EpaymentsOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, tradeRsp)
}

func (trade *Trade) jdTrade(t defs.Trade) (tradeRsp defs.TradeRsp, errCode int, err error) {
	return new(method.Jd).Trade(t.OrderId, code.JdMethod)
}

func (trade *Trade) allpayTrade(t defs.Trade) (tradeRsp defs.TradeRsp, errCode int, err error) {
	return new(method.Allpay).Trade(t.OrderId, t.MethodCode)
}
func (trade *Trade) alipayTrade(t defs.Trade) (tradeRsp defs.TradeRsp, errCode int, err error) {
	return new(method.Alipay).Trade(t.OrderId, t.MethodCode)
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
	switch org {
	case JdOrg:
		closedRsp, errCode, err = trade.jdClosed(*closed)
	case AllpayOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case AlipayOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case WechatOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case EpaymentsOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, closedRsp)
}

func (trade *Trade) jdClosed(closed defs.Closed) (closedRsp defs.ClosedRsp, errCode int, err error) {
	closedArg := method.JdClosedArg{
		OrderId:  closed.OrderId,
		Currency: closed.Currency,
		TotalFee: closed.TotalFee,
	}
	return new(method.Jd).Closed(closedArg)
}
