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

	logisticsHandle := common.GetLogisticsHandler(l.OrgCode)
	if logisticsHandle == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": common.NotSupportPaymentOrgCode, "message": common.NotSupportPaymentOrgMsg})
		return
	}

	req := defs.LogisticsReq{
		OrderId:          l.OrderId,
		LogisticsNo:      l.LogisticsNo,
		LogisticsCompany: l.LogisticsCompany,
		OrgCode:          l.OrgCode,
	}
	logisticsRsp, errCode, err := logisticsHandle(req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, logisticsRsp)
}
