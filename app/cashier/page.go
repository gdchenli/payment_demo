package cashier

import "github.com/gin-gonic/gin"

type Page struct{}

func (page *Page) Router(router *gin.Engine) {
	r := router.Group("/cashier")
	{
		r.POST("/cart", page.order)
		r.Group("/success", page.success)
	}
}

func (page *Page) order(ctx *gin.Context) {

}

func (page *Page) success(ctx *gin.Context) {

}
