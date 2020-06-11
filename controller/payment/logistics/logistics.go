package logistics

import (
	"net/http"
	"payment_demo/controller/payment/common"
	"payment_demo/internal/common/defs"

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
	l := new(Req)
	ctx.ShouldBind(l)

	uploadLogisticsHandle := common.GetUploadLogisticsHandler(l.OrgCode)
	if uploadLogisticsHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": common.NotSupportPaymentOrgCode, "message": common.NotSupportPaymentOrgMsg})
		return
	}

	req := defs.UploadLogisticsReq{
		OrderId:          l.OrderId,
		LogisticsNo:      l.LogisticsNo,
		LogisticsCompany: l.LogisticsCompany,
		OrgCode:          l.OrgCode,
	}
	logisticsRsp, errCode, err := uploadLogisticsHandle(req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, logisticsRsp)
}
