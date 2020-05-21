package cashier

import "github.com/gin-gonic/gin"

type Pay struct{}

func (pay *Pay) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/submit", pay.submit)
		r.Group("/notify", pay.notify)
		r.Group("/verify", pay.verify)
		r.GET("/status", pay.status)
	}
}

func (pay *Pay) submit(ctx *gin.Context) {

}

func (pay *Pay) notify(ctx *gin.Context) {

}

func (pay *Pay) verify(ctx *gin.Context) {

}

func (pay *Pay) status(ctx *gin.Context) {

}
