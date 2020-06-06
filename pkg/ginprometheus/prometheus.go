package ginprometheus

import (
	"bytes"

	"github.com/gin-gonic/gin"

	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	// "gopkg.in/gin-gonic/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

var defaultMetricPath = "/metrics"

//请求路径URL规则匹配函数，用于支持自定义自己的请求路径匹配方式
//例如，category/100.html 匹配成category/:id.html
type RequestCounterURLLabelMappingFn func(c *gin.Context) string

// Prometheus contains the metrics gathered by the instance and its path
type Prometheus struct {
	reqCnt               *prometheus.CounterVec
	reqDur, reqSz, resSz *prometheus.SummaryVec
	router               *gin.Engine
	listenAddress        string

	Ppg PrometheusPushGateway

	MetricsPath string

	ReqCntURLLabelMappingFn RequestCounterURLLabelMappingFn
}

// PrometheusPushGateway contains the configuration for pushing to a Prometheus pushgateway (optional)
type PrometheusPushGateway struct {

	// Push interval in seconds
	PushIntervalSeconds time.Duration

	// Push Gateway URL in format http://domain:port
	// where JOBNAME can be any string of your choice
	PushGatewayURL string

	// Local metrics URL where metrics are fetched from, this could be ommited in the future
	// if implemented using prometheus common/expfmt instead
	MetricsURL string

	// pushgateway job name, defaults to "gin"
	Job string
}

// NewPrometheus generates a new set of metrics with a certain subsystem name
func NewPrometheus() *Prometheus {

	p := &Prometheus{
		MetricsPath: defaultMetricPath,
		ReqCntURLLabelMappingFn: func(c *gin.Context) string {
			return c.Request.URL.Path // 默认返回源请求路径
		},
	}
	p.registerMetrics()

	return p
}

// SetPushGateway sends metrics to a remote pushgateway exposed on pushGatewayURL
// every pushIntervalSeconds. Metrics are fetched from metricsURL
func (p *Prometheus) SetPushGateway(pushGatewayURL, metricsURL string, pushIntervalSeconds time.Duration) {
	p.Ppg.PushGatewayURL = pushGatewayURL
	p.Ppg.MetricsURL = metricsURL
	p.Ppg.PushIntervalSeconds = pushIntervalSeconds
	p.startPushTicker()
}

// Set pushgateway job name, defaults to "gin"
func (p *Prometheus) SetPushGatewayJob(j string) {
	p.Ppg.Job = j
}

// SetListenAddress for exposing metrics on address. If not set, it will be exposed at the
// same address of the gin engine that is being used
func (p *Prometheus) SetListenAddress(address string) {
	p.listenAddress = address
	// if p.listenAddress != "" {
	// 	p.router = Default()
	// }
}

func (p *Prometheus) setMetricsPath(e *gin.Engine) {

	if p.listenAddress != "" {
		p.router.GET(p.MetricsPath, prometheusHandler())
		p.runServer()
	} else {
		e.GET(p.MetricsPath, prometheusHandler())
	}
}

func (p *Prometheus) setMetricsPathWithAuth(e *gin.Engine, accounts gin.Accounts) {

	if p.listenAddress != "" {
		p.router.GET(p.MetricsPath, gin.BasicAuth(accounts), prometheusHandler())
		p.runServer()
	} else {
		e.GET(p.MetricsPath, gin.BasicAuth(accounts), prometheusHandler())
	}

}

func (p *Prometheus) runServer() {
	if p.listenAddress != "" {
		go p.router.Run(p.listenAddress)
	}
}

func (p *Prometheus) getMetrics() []byte {
	response, _ := http.Get(p.Ppg.MetricsURL)

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)

	return body
}

func (p *Prometheus) getPushGatewayUrl() string {
	h, _ := os.Hostname()
	if p.Ppg.Job == "" {
		p.Ppg.Job = "gin"
	}
	return p.Ppg.PushGatewayURL + "/metrics/job/" + p.Ppg.Job + "/instance/" + h
}

func (p *Prometheus) sendMetricsToPushGateway(metrics []byte) {
	req, err := http.NewRequest("POST", p.getPushGatewayUrl(), bytes.NewBuffer(metrics))
	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Error("Error sending to push gatway: " + err.Error())
	}
}

func (p *Prometheus) startPushTicker() {
	ticker := time.NewTicker(time.Second * p.Ppg.PushIntervalSeconds)
	go func() {
		for range ticker.C {
			p.sendMetricsToPushGateway(p.getMetrics())
		}
	}()
}

func (p *Prometheus) registerMetrics() {

	p.reqCnt = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Subsystem: "http",
			Name:      "requests_total",
			Help:      "How many HTTP requests processed, partitioned by status code and HTTP request path.",
		},
		[]string{"code", "method", "request_path", "host"},
	)

	if err := prometheus.Register(p.reqCnt); err != nil {
		log.Info("reqCnt could not be registered: ", err)
	} else {
		log.Info("reqCnt registered.")
	}

	p.reqDur = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Subsystem: "http",
			Name:      "request_duration_seconds",
			Help:      "The HTTP request latencies in seconds.",
		},
		[]string{"request_path"},
	)

	if err := prometheus.Register(p.reqDur); err != nil {
		log.Info("reqDur could not be registered: ", err)
	} else {
		log.Info("reqDur registered.")
	}

	p.reqSz = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Subsystem: "http",
			Name:      "request_size_bytes",
			Help:      "The HTTP request sizes in bytes.",
		},
		[]string{"request_path"},
	)

	if err := prometheus.Register(p.reqSz); err != nil {
		log.Info("reqSz could not be registered: ", err)
	} else {
		log.Info("reqSz registered.")
	}

	p.resSz = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Subsystem: "http",
			Name:      "response_size_bytes",
			Help:      "The HTTP response sizes in bytes.",
		},
		[]string{"request_path"},
	)

	if err := prometheus.Register(p.resSz); err != nil {
		log.Info("resSz could not be registered: ", err)
	} else {
		log.Info("resSz registered.")
	}

}

// Use adds the middleware to a gin engine.
func (p *Prometheus) Use(e *gin.Engine) {
	e.Use(p.handlerFunc())
	p.setMetricsPath(e)
}

// UseWithAuth adds the middleware to a gin engine with BasicAuth.
func (p *Prometheus) UseWithAuth(e *gin.Engine, accounts gin.Accounts) {
	e.Use(p.handlerFunc())
	p.setMetricsPathWithAuth(e, accounts)
}

func (p *Prometheus) handlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.String() == p.MetricsPath {
			c.Next()
			return
		}

		start := time.Now()
		reqSz := make(chan int)
		go computeApproximateRequestSize(c.Request, reqSz)
		reqSize := float64(<-reqSz)

		c.Next()

		status := strconv.Itoa(c.Writer.Status())
		elapsed := float64(time.Since(start)) / float64(time.Second)
		resSz := float64(c.Writer.Size())

		requestPath := p.ReqCntURLLabelMappingFn(c)

		p.reqDur.WithLabelValues(requestPath).Observe(elapsed)
		p.reqCnt.WithLabelValues(status, c.Request.Method, requestPath, c.Request.Host).Inc()
		p.reqSz.WithLabelValues(requestPath).Observe(reqSize)
		p.resSz.WithLabelValues(requestPath).Observe(resSz)
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// From https://github.com/DanielHeckrath/gin-prometheus/blob/master/gin_prometheus.go
func computeApproximateRequestSize(r *http.Request, out chan int) {
	s := 0
	if r.URL != nil {
		s = len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// N.B. r.Form and r.MultipartForm are assumed to be included in r.URL.

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}

	out <- s
}
