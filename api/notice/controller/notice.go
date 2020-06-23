package controller

import (
	"io/ioutil"
	"net/http"
	"payment_demo/internal/service/notice"

	"github.com/gin-gonic/gin"
)

type Notice struct{}

func (n *Notice) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{
		r.POST("/notify/:org/:method", n.notify) //异步通知
		r.POST("/verify/:org/:method", n.verify) //同步通知
	}
}

func (n *Notice) notify(ctx *gin.Context) {
	notifyBytes, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": "1", "message": err.Error()})
		return
	}
	defer func() {
		ctx.Request.Body.Close()
	}()
	query := string(notifyBytes)

	orgCode := ctx.Param("org")
	noticeService, errCode, err := notice.New(orgCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}
	notifyRsp, errCode, err := noticeService.Notify(query, orgCode, ctx.Param("method"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "message": "success", "data": notifyRsp})
}

func (n *Notice) verify(ctx *gin.Context) {
	ctx.Request.ParseForm()
	query := ctx.Request.PostForm.Encode()
	if query == "" {
		query = ctx.Request.URL.Query().Encode()
	}

	orgCode := ctx.Param("org")
	noticeService, errCode, err := notice.New(orgCode)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}
	verifyRsp, errCode, err := noticeService.Verify(query, orgCode, ctx.Param("method"))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": errCode, "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, verifyRsp)
}
