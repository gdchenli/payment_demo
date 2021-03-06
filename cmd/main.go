package main

import (
	"fmt"
	"net/http"
	"payment_demo/api"
	"payment_demo/pkg/config"
	"payment_demo/pkg/ginprometheus"
	"payment_demo/pkg/grace"
	"payment_demo/pkg/log"
	"payment_demo/pkg/recovery"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	Port                = "server.http_port"
	Mode                = "server.run_mod"
	ReadTimeout         = "server.read_timeout"
	WriteTimeout        = "server.write_timeout"
	MonitorStatus       = "monitor.status"
	MetricsAuthStatus   = "metrics.auth_status"
	MetricsAuthUser     = "metrics.auth_user"
	MetricsAuthPassword = "metrics.auth_password"
)

func init() {
	log.Init()
}

func main() {
	gin.SetMode(config.GetInstance().GetString(Mode))
	router := gin.Default()
	//router.StaticFS("/assets", http.Dir("../../assets"))

	router.Use(recovery.Recovery(RecoveryHandler))
	registerMonitor(router)
	registerRouter(router)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.GetInstance().GetInt(Port)),
		Handler:      router,
		ReadTimeout:  config.GetInstance().GetDuration(ReadTimeout) * time.Second,
		WriteTimeout: config.GetInstance().GetDuration(WriteTimeout) * time.Second,
	}

	//通过不同OS用不同的方式构建约束
	//windows采用 golang 自带http包
	//linux 和 darwin 采用 facebookgo gracehttp包
	err := grace.Serve(server)
	if err != nil {
		panic(err)
	}
}

func registerRouter(router *gin.Engine) {
	new(api.Payment).Router(router)
	new(api.Logistics).Router(router)
}

func registerMonitor(router *gin.Engine) {
	//监控
	if !config.GetInstance().GetBool(MonitorStatus) {
		return
	}

	p := ginprometheus.NewPrometheus()

	//指标验证
	if !config.GetInstance().GetBool(MetricsAuthStatus) {
		p.Use(router)
		return
	}
	p.UseWithAuth(router, gin.Accounts{config.GetInstance().GetString(MetricsAuthUser): config.GetInstance().GetString(MetricsAuthPassword)})
}

//全局panic recovery
func RecoveryHandler(c *gin.Context, err interface{}) {
	c.AbortWithStatus(http.StatusBadGateway)
}
