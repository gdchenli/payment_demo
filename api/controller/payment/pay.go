package payment

import (
	"net/http"
	"payment_demo/api/controller/common"
	"payment_demo/api/validate"

	"github.com/gin-gonic/gin"
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
	o := new(validate.Order)
	ctx.ShouldBind(o)

	if errCode, err := o.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	payHandle := common.GetPayHandler(o.OrgCode)
	if payHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": common.NotSupportPaymentOrgCode, "message": common.NotSupportPaymentOrgMsg})
		return
	}

	form, errCode, err := payHandle(*o)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": form})
}
