package notify

import (
	"io/ioutil"
	"net/http"
	"payment_demo/controller/common"

	"github.com/gin-gonic/gin"
)

type Notify struct{}

func (notify *Notify) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/notify/:org/:method", notify.notify) //异步通知
	}
}

func (notify *Notify) notify(ctx *gin.Context) {
	notifyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": common.NotSupportPaymentOrgCode, "message": err.Error()})
		return
	}
	defer func() {
		ctx.Request.Body.Close()
	}()
	query := string(notifyBytes)

	notifyHandle := common.GetNotifyHandler(ctx.Param("org"))
	if notifyHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": common.NotSupportPaymentOrgCode, "message": common.NotSupportPaymentOrgMsg})
		return
	}

	notifyRsp, errCode, err := notifyHandle(query, ctx.Param("method"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": notifyRsp})
}
