package cashier

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Callback struct{}

func (callback *Callback) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/callback/:org/:method", callback.callback) //同步通知
	}
}

func (callback *Callback) callback(ctx *gin.Context) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()
	if query == "" {
		query = ctx.Request.URL.Query().Encode()
	}

	callbackHandle := getCallbackHandler(ctx.Param("org"))
	if callbackHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": NotSupportPaymentOrgCode, "message": NotSupportPaymentOrgMsg})
		return
	}

	callBackRsp, errCode, err := callbackHandle(query, ctx.Param("method"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, callBackRsp)
}
