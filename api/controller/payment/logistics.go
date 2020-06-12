package payment

import (
	"net/http"
	"payment_demo/api/controller/common"
	"payment_demo/api/validate/payment"

	"github.com/gin-gonic/gin"
)

type Logistics struct{}

func (logistics *Logistics) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/logistics/upload", logistics.upload) //上传物流信息
	}
}

func (logistics *Logistics) upload(ctx *gin.Context) {
	l := new(payment.UploadLogisticsReq)
	ctx.ShouldBind(l)

	uploadLogisticsHandle := common.GetUploadLogisticsHandler(l.OrgCode)
	if uploadLogisticsHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": common.NotSupportPaymentOrgCode, "message": common.NotSupportPaymentOrgMsg})
		return
	}

	logisticsRsp, errCode, err := uploadLogisticsHandle(*l)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, logisticsRsp)
}
