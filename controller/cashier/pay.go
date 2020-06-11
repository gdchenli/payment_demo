package cashier

import (
	"net/http"
	"payment_demo/internal/common/defs"

	"github.com/gin-gonic/gin"
)

const (
	JdOrg        = "jd"
	AllpayOrg    = "allpay"
	WechatOrg    = "wechat"
	EpaymentsOrg = "epayments"
	AlipayOrg    = "alipay"
)

type Pay struct{}

func (pay *Pay) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/order/pay", pay.Pay) //发起支付
	}
}

//发起支付
func (pay *Pay) Pay(ctx *gin.Context) {
	order := new(defs.Order)
	ctx.ShouldBind(order)

	if errCode, err := order.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	payHandle := getPayHandler(order.OrgCode)
	if payHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": NotSupportPaymentOrgCode, "message": NotSupportPaymentOrgMsg})
		return
	}

	form, errCode, err := payHandle(*order)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": form})
}
