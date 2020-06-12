package verify

import (
	"net/http"
	"payment_demo/controller/common"

	"github.com/gin-gonic/gin"
)

type Verify struct{}

func (verify *Verify) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/verify/:org/:method", verify.verify) //同步通知
	}
}

func (verify *Verify) verify(ctx *gin.Context) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()
	if query == "" {
		query = ctx.Request.URL.Query().Encode()
	}

	verifyHandle := common.GetVerifyHandler(ctx.Param("org"))
	if verifyHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": common.NotSupportPaymentOrgCode, "message": common.NotSupportPaymentOrgMsg})
		return
	}

	callBackRsp, errCode, err := verifyHandle(query, ctx.Param("method"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, callBackRsp)
}
