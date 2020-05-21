package cart

import "github.com/gin-gonic/gin"

type Checkout struct{}

func (checkout *Checkout) Router(router *gin.Engine) {
	r := router.Group("/cart")
	{
		r.POST("/checkout", checkout.checkout)
	}
}

func (checkout *Checkout) index(ctx *gin.Context) {

}

func (checkout *Checkout) checkout(ctx *gin.Context) {

}
