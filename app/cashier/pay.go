package cashier

import (
	"azoya/nova/binding"
	"errors"
	"fmt"
	"net/http"
	"payment_demo/internal/common/defs"
	"payment_demo/internal/method"

	"github.com/gin-gonic/gin"
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
	}
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
	switch order.OrgCode {
	case JdOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case AllpayOrg:
		payStr, errCode, err = pay.allpayAmpSubmit(*order)
	case AlipayOrg:
		payStr, errCode, err = pay.alipayAmpSubmit(*order)
	case WechatOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	case EpaymentsOrg:
		err = errors.New(NotSupportPaymentOrgMsg)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
	}

	if errCode, err = order.Validate(); err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": payStr})

}

func (pay *Pay) allpayAmpSubmit(order defs.Order) (form string, errCode int, err error) {
	return new(method.Allpay).AmpSubmit(order)
}

func (pay *Pay) alipayAmpSubmit(order defs.Order) (form string, errCode int, err error) {
	return new(method.Alipay).AmpSubmit(order)
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

	var payMethod method.PayMethod
	switch order.OrgCode {
	case JdOrg:
		payMethod = new(method.Jd)
	case AllpayOrg:
		payMethod = new(method.Allpay)
	case AlipayOrg:
		payMethod = new(method.Alipay)
	case EpaymentsOrg:
		payMethod = new(method.Epayments)
	default:
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
