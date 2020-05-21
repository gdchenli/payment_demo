package cashier

import (
	"azoya/nova/binding"
	"fmt"
	"net/http"
	"payment_demo/internal/common/defs"
	"payment_demo/internal/method"

	"github.com/gin-gonic/gin"
)

type Pay struct{}

func (pay *Pay) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/submit", pay.submit)
		r.Group("/notify/:org/:code", pay.notify)
		r.Group("/verify/:org/:code", pay.verify)
		r.GET("/status", pay.status)
	}
}

func (pay *Pay) submit(ctx *gin.Context) {
	var errCode int
	var err error
	var form string
	order := new(defs.Order)
	ctx.ShouldBind(order)

	if errCode, err = order.Validate(); err != nil {
		ctx.Data(http.StatusOK, "text/html", []byte(err.Error()))
		return
	}

	switch order.OrgCode {
	case "Jd":
		jdPayArg := method.JdPayArg{
			OrderId:  order.OrderId,
			TotalFee: order.TotalFee,
			Currency: order.Currency,
			UserId:   order.UserId,
		}
		form, errCode, err = new(method.Jd).Submit(jdPayArg)
	}

	if err != nil {
		fmt.Println(errCode)
		ctx.Data(http.StatusOK, binding.MIMEHTML, []byte(err.Error()))
		return
	}

	ctx.Data(http.StatusOK, "text/html", []byte(form))
}

func (pay *Pay) notify(ctx *gin.Context) {

}

func (pay *Pay) verify(ctx *gin.Context) {

}

func (pay *Pay) status(ctx *gin.Context) {

}
