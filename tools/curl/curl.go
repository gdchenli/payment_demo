package curl

import (
	"azoya/nova"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

// GetJSON executes HTTP GET against specified url and tried to parse
// the response into out object.
func GetJSON(c *nova.Context, operation, url string, out interface{}) error {
	span, _ := opentracing.StartSpanFromContext(c.Context(), "HTTP GET: "+operation)
	defer span.Finish()

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err = span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var reader io.ReadCloser
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		reader, _ = gzip.NewReader(res.Body)
		defer reader.Close()
	default:
		reader = res.Body
	}

	if res.StatusCode >= 400 {
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	decoder := json.NewDecoder(reader)
	return decoder.Decode(out)
}

// GetJSON executes HTTP GET against specified url and tried to parse
// the response into out object.
func GetJSONWithHeader(c *nova.Context, params map[string]string, operation, url string, out interface{}) error {
	span, _ := opentracing.StartSpanFromContext(c.Context(), "HTTP GET: "+operation)
	defer span.Finish()

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	carrier := opentracing.HTTPHeadersCarrier(req.Header)
	err = span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, carrier)
	if err != nil {
		return err
	}

	for k, v := range params {
		req.Header.Set(k, v)
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}
	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(out)
}

// PostJSON executes HTTP GET against specified url and tried to parse
// the response into out object.
func PostJSON(c *nova.Context, url string, out interface{}, body io.Reader) error {

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return errors.New(string(body))
	}

	decoder := json.NewDecoder(res.Body)
	return decoder.Decode(out)
}

// PostXML executes HTTP POST against specified url and tried to parse
// the response into out object.
func PostXML(url string, out interface{}, body io.Reader) error {

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, body)
	req.Header.Add("content-type", "application/xml; charset=utf-8")
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	//状态码异常
	if resp.StatusCode != http.StatusOK {
		fmt.Println("resp.StatusCode", resp.StatusCode)
		return errors.New(CurlErrorMsg)
	}

	//读取内容
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("respBytes", string(respBytes))

	return xml.Unmarshal(respBytes, out)

}

const CurlErrorMsg = "状态码异常"

type Curl struct {
	RequestUrl    string
	RequestMethod string
	RequestData   string
	HeaderMap     map[string]string
}

//Get请求
func (curl *Curl) RequestGet() (curlResult string, e error) {
	defer func() {
		if r := recover(); r != nil {
			logrus.Infof("CURL 异常错误 %v\n", r)
		}
	}()

	timeout := time.Duration(30 * time.Second) //超时时间30s

	//发起请求
	buf := new(bytes.Buffer)
	client := &http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest(curl.RequestMethod, curl.RequestUrl, buf)
	if len(curl.HeaderMap) > 0 {
		for key, val := range curl.HeaderMap {
			req.Header.Set(key, val)
		}
	}
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//状态码异常
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(CurlErrorMsg)
	}

	//读取内容
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	//返回结果
	return string(responseBytes), nil

}
