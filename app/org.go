package app

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type Org struct{}

func (org *Org) Router(router *gin.Engine) {
	r := router.Group("/payment")
	{

	}
	fmt.Println(r)
}
