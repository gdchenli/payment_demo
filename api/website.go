package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Website struct{}

func (website *Website) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{

	}
	fmt.Println(r)
}
