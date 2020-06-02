package cashier

import (
	"fmt"
	"net/http"
	"payment_demo/internal/common/defs"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

const (
	JdOrg        = "jd"
	AllpayOrg    = "allpay"
	WechatOrg    = "wechat"
	EpaymentsOrg = "epayments"
	AlipayOrg    = "alipay"
)

const (
	NotSupportPaymentOrgCode = "10101"
	NotSupportPaymentOrgMsg  = "不支持该支付机构"
	NotifyFailMsg            = "fail"
)

type Pay struct{}

func (pay *Pay) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/order/submit", pay.orderSubmit) //发起支付
		r.POST("/amp/submit", pay.ampSubmit)     //支付宝小程序，支付参数
		r.POST("/order/qrcode", pay.qrCode)      //二维码
	}
}

func (pay *Pay) qrCode(ctx *gin.Context) {
	var errCode int
	var err error
	var payStr string
	order := new(defs.Order)
	ctx.ShouldBind(order)
	order.UserAgentType = 7

	if errCode, err = order.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	payMethod := getPayMethod(order.OrgCode)
	if payMethod == nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(NotSupportPaymentOrgMsg))
		return
	}
	payStr, errCode, err = payMethod.OrderQrCode(*order)

	if errCode, err = order.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": payStr})
}

//发起支付
func (pay *Pay) ampSubmit(ctx *gin.Context) {
	var errCode int
	var err error
	var payStr string
	order := new(defs.Order)
	ctx.ShouldBind(order)
	order.UserAgentType = 7

	if errCode, err = order.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	payMethod := getPayMethod(order.OrgCode)
	if payMethod == nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(NotSupportPaymentOrgMsg))
		return
	}

	payStr, errCode, err = payMethod.AmpSubmit(*order)
	if errCode, err = order.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": payStr})
}

//发起支付
func (pay *Pay) orderSubmit(ctx *gin.Context) {
	var errCode int
	var err error
	var form string
	order := new(defs.Order)
	ctx.ShouldBind(order)

	if errCode, err = order.Validate(); err != nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(err.Error()))
		return
	}

	payMethod := getPayMethod(order.OrgCode)
	if payMethod == nil {
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(NotSupportPaymentOrgMsg))
		return
	}

	form, errCode, err = payMethod.OrderSubmit(*order)
	if err != nil {
		fmt.Println(errCode)
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(err.Error()))
		return
	}

	ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(form))
}
