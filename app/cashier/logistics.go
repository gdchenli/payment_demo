package cashier

import (
	"errors"
	"net/http"
	"payment_demo/internal/common/defs"
	"payment_demo/internal/method"

	"github.com/gin-gonic/gin"
)

type Logistics struct{}

func (logistics *Logistics) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/logistics/upload", logistics.logisticsUpload) //物流信息回传
	}
}

func (logistics *Logistics) logisticsUpload(ctx *gin.Context) {
	var errCode int
	var err error
	var logisticsRsp defs.LogisticsRsp

	l := new(defs.Logistics)
	ctx.ShouldBind(l)

	switch l.OrgCode {
	case JdOrg:
		logisticsRsp, errCode, err = logistics.JdLogisticsUpload(*l)
	default:
		err = errors.New(NotSupportPaymentOrgMsg)
	}

	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, logisticsRsp)
}

func (logistics *Logistics) JdLogisticsUpload(l defs.Logistics) (logisticsRsp defs.LogisticsRsp, errCode int, err error) {
	arg := method.JdLogisticsArg{
		OrderId:          l.OrderId,
		LogisticsCompany: l.LogisticsCompany,
		LogisticsNo:      l.LogisticsNo,
	}

	return new(method.Jd).LogisticsUpload(arg)
}
