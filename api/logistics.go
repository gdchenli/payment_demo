package api

import (
	"net/http"
	"payment_demo/api/validate"
	"payment_demo/internal/service"

	"github.com/gin-gonic/gin"
)

type Logistics struct{}

func (l *Logistics) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/logistics/upload", l.uploadLogistics) //上传物流信息
	}
}

func (l *Logistics) uploadLogistics(ctx *gin.Context) {
	arg := new(validate.UploadLogisticsArg)
	ctx.ShouldBind(arg)

	logisticsUpload, errCode, err := service.NewLogictics()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	logisticsRsp, errCode, err := logisticsUpload.Upload(*arg)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"message": err.Error(), "code": errCode})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": logisticsRsp})
}
