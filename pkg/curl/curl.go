package curl

import (
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var requestTimeout = 15 * time.Second

// PostJSON executes HTTP GET against specified url and tried to parse
// the response into out object.
func PostJSON(url string, out interface{}, body io.Reader) error {

	client := &http.Client{Timeout: requestTimeout}

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

	client := &http.Client{Timeout: requestTimeout}

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
		return errors.New("状态码异常")
	}

	//读取内容
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	fmt.Println("respBytes", string(respBytes))

	return xml.Unmarshal(respBytes, out)

}

func GetJSONReturnByte(url string) (out []byte, err error) {
	client := &http.Client{Timeout: requestTimeout}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return out, err
	}

	res, err := client.Do(req)
	if err != nil {
		return out, err
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

	out, err = ioutil.ReadAll(reader)
	if err != nil {
		return out, err
	}

	if res.StatusCode >= http.StatusBadRequest {
		return []byte{}, errors.New(string(out))
	}

	return out, nil
}
