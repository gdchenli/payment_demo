package jd

import (
	"crypto"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	request2 "payment_demo/api/trade/request"
	response2 "payment_demo/api/trade/response"
	"payment_demo/pkg/curl"
	"regexp"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	CustomTradeType = "0" //交易类型
)

const (
	SearchTradeWait    = "1" //等待支付
	SearchTradeProcess = "2" //交易成功
	SearchTradeClosed  = "3" //交易关闭
	SearchTradeError   = "4" //交易失败
	SearchTradeRevoked = "5" //交撤销
	SearchTradeNotPay  = "6" //未支付
	SearchTradeRefund  = "7" //转入退款
)

const (
	SearchTradeBuildSignErrCode                    = 10501
	SearchTradeBuildSignErrMessage                 = "查询交易流水，签名生成失败"
	SearchTradeDesKeyFormatErrCode                 = 10502
	SearchTradeDesKeyFormatErrMessage              = "查询交易流水，desKey格式错误"
	SearchTradeRequestDataEncryptFailedErrCode     = 10503
	SearchTradeRequestDataEncryptFailedErrMessage  = "查询交易流水，请求数据加密失败"
	SearchTradeNetErrCode                          = 10504
	SearchTradeNetErrMessage                       = "查询交易流水,网络错误"
	SearchTradeResponseDataEncryptFormatErrCode    = 10505
	SearchTradeResponseDataEncryptFormatErrMessage = "查询交易流水,返回加密数据格式错误"
	SearchTradeResponseDataDecryptFailedErrCode    = 10506
	SearchTradeResponseDataDecryptFailedErrMessage = "查询交易流水,解密返回数据失败"
	SearchTradeResponseDataDecryptFormatErrCode    = 10507
	SearchTradeResponseDataDecryptFormatErrMessage = "查询交易流水,解密数据格式错误"
	SearchTradeResponseDataSignErrCode             = 10508
	SearchTradeResponseDataSignErrMessage          = "查询交易流水,返回数据签名校验错误"
)

type SearchWithoutSignRequest struct {
	XMLName   xml.Name `xml:"jdpay" json:"-"`
	Version   string   `xml:"version" json:"version"`     //版本
	Merchant  string   `xml:"merchant" json:"merchant"`   //商户号
	TradeNum  string   `xml:"tradeNum" json:"tradeNum"`   //订单编号
	OTradeNum string   `xml:"oTradeNum" json:"oTradeNum"` //原交易流水号
	TradeType string   `xml:"tradeType" json:"tradeType"` //交易类型
}

type SearchWithSignRequest struct {
	XMLName   xml.Name `xml:"jdpay" json:"-"`
	Version   string   `xml:"version" json:"version"`     //版本
	Merchant  string   `xml:"merchant" json:"merchant"`   //商户号
	TradeNum  string   `xml:"tradeNum" json:"tradeNum"`   //订单编号
	OTradeNum string   `xml:"oTradeNum" json:"oTradeNum"` //原交易流水号
	TradeType string   `xml:"tradeType" json:"tradeType"` //交易类型
	Sign      string   `xml:"sign" json:"sign"`           //签名
}

type SearchWithEncrypt struct {
	XMLName  xml.Name `xml:"jdpay" json:"-"`
	Version  string   `xml:"version" json:"version"`   //版本
	Merchant string   `xml:"merchant" json:"merchant"` //商户号
	Encrypt  string   `xml:"encrypt" json:"encrypt"`   //加密数据
}

type SearchResult struct {
	XMLName  xml.Name        `xml:"jdpay" json:"-"`
	Version  string          `xml:"version" json:"version"`   //版本号
	Merchant string          `xml:"merchant" json:"merchant"` //商户号
	Result   SearchResultRsp `xml:"result" json:"result"`     //交易结果
	Encrypt  string          `xml:"encrypt" json:"encrypt"`   //加密信息
}

type SearchResultRsp struct {
	Code string `xml:"code" json:"code"` //交易返回码
	Desc string `xml:"desc" json:"desc"` //返回码信息
}

type SearchDecryptRsp struct {
	XMLName   xml.Name         `xml:"jdpay" json:"-"`
	Merchant  string           `xml:"merchant" json:"merchant"`   //商户号
	TradeNum  string           `xml:"tradeNum" json:"tradeNum"`   //订单编号
	TradeType string           `xml:"tradeType" json:"tradeType"` //交易类型
	Result    SearchResultRsp  `xml:"result" json:"result"`       //交易结果
	Sign      string           `xml:"sign" json:"sign"`           //数据签名
	Amount    int64            `xml:"amount" json:"amount"`       //人民币支付总金额
	Status    string           `xml:"status" json:"status"`       //交易状态
	PayList   SearchPayListRsp `xml:"payList" json:"payList"`     //支付方式明细
}

type SearchPayListRsp struct {
	Pay []SearchPayRsp `xml:"pay" json:"pay"`
}

type SearchPayRsp struct {
	PayType   int    `xml:"payType" json:"payType"`     //支付方式
	Amount    int64  `xml:"amount" json:"amount"`       //交易金额
	Currency  string `xml:"currency" json:"currency"`   //交易币种
	TradeTime string `xml:"tradeTime" json:"tradeTime"` //交易时间
}

func (jd *Jd) SearchTrade(paramMap map[string]string, req request2.SearchTradeArg) (searchTradeRsp response2.SearchTradeRsp, errCode int, err error) {
	searchWithoutSignRequest := SearchWithoutSignRequest{
		Version:   Version,
		Merchant:  paramMap["merchant"],
		TradeNum:  req.OrderId,
		OTradeNum: "",
		TradeType: CustomTradeType,
	}
	xmlBytes, err := xml.Marshal(searchWithoutSignRequest)
	xmlStr := xml.Header + string(xmlBytes)
	xmlStr = strings.Replace(xmlStr, "\r", "", -1)
	xmlStr = strings.Replace(xmlStr, "\n", "", -1)
	xmlStr = strings.Replace(xmlStr, "\t", "", -1)
	reg, _ := regexp.Compile(">\\s+<")
	xmlStr = reg.ReplaceAllString(xmlStr, "><")
	reg, _ = regexp.Compile("\\s+\\/>")
	xmlStr = reg.ReplaceAllString(xmlStr, "/>")
	//fmt.Println("without sign xml", xmlStr)

	//生成签名
	sha256 := HaSha256(xmlStr)
	signBytes, err := SignPKCS1v15([]byte(sha256), []byte(paramMap["private_key"]), crypto.Hash(0))
	if err != nil {
		logrus.Errorf("org:jd,"+SearchTradeBuildSignErrMessage+",request:%+v,errCode:%v,err:%v", req, SearchTradeBuildSignErrCode, err.Error())
		return searchTradeRsp, SearchTradeBuildSignErrCode, errors.New(SearchTradeBuildSignErrMessage)
	}
	sign := base64.StdEncoding.EncodeToString(signBytes)
	searchWithSignRequest := SearchWithSignRequest{
		Version:   searchWithoutSignRequest.Version,
		Merchant:  searchWithoutSignRequest.Merchant,
		TradeNum:  searchWithoutSignRequest.TradeNum,
		OTradeNum: searchWithoutSignRequest.OTradeNum,
		TradeType: searchWithoutSignRequest.TradeType,
		Sign:      sign,
	}
	xmlBytes, err = xml.Marshal(searchWithSignRequest)
	xmlStr = strings.TrimRight(xml.Header, "\n") + string(xmlBytes)
	//fmt.Println("with sign xml", xmlStr)
	//searchTradeRsp.DecryptRes = xmlStr

	desKeyBytes, err := base64.StdEncoding.DecodeString(paramMap["des_key"])
	if err != nil {
		logrus.Errorf("org:jd,"+SearchTradeDesKeyFormatErrMessage+",request:%+v,errCode:%v,err:%v", req, SearchTradeDesKeyFormatErrCode, err.Error())
		return searchTradeRsp, SearchTradeDesKeyFormatErrCode, errors.New(SearchTradeDesKeyFormatErrMessage)
	}
	encryptBytes, err := TripleEcbDesEncrypt([]byte(xmlStr), desKeyBytes)
	if err != nil {
		logrus.Errorf("org:jd,"+SearchTradeRequestDataEncryptFailedErrMessage+",query:%+v,errCode:%v,err:%v", req, SearchTradeRequestDataEncryptFailedErrCode, err.Error())
		return searchTradeRsp, SearchTradeRequestDataEncryptFailedErrCode, errors.New(SearchTradeRequestDataEncryptFailedErrMessage)
	}
	reqEncrypt := DecimalByteSlice2HexString(encryptBytes)
	reqEncrypt = base64.StdEncoding.EncodeToString([]byte(reqEncrypt))
	searchWithEncrypt := SearchWithEncrypt{
		Version:  Version,
		Merchant: paramMap["merchant"],
		Encrypt:  reqEncrypt,
	}
	xmlBytes, err = xml.Marshal(searchWithEncrypt)
	xmlStr = strings.TrimRight(xml.Header, "\n") + string(xmlBytes)
	//fmt.Println("with 3des xml", xmlStr)
	//searchTradeRsp.EncryptRes = xmlStr

	var searchResult SearchResult
	playLoad := strings.NewReader(xmlStr)
	err = curl.PostXML(paramMap["trade_way"], &searchResult, playLoad)
	if err != nil {
		logrus.Errorf("org:jd,"+SearchTradeNetErrMessage+",request:%+v,errCode:%v,err:%v", req, SearchTradeNetErrCode, err.Error())
		return searchTradeRsp, SearchTradeNetErrCode, errors.New(SearchTradeNetErrMessage)
	}
	/*searchResultBytes, err := xml.Marshal(searchResult)
	if err != nil {
		logrus.Errorf("org:jd,"+SearchTradeResponseDataEncryptFormatErrMessage+",request:%+v,response:%+v,errCode:%v,err:%v", arg, searchResult, SearchTradeResponseDataEncryptFormatErrCode, err.Error())
		return searchTradeRsp, SearchTradeResponseDataEncryptFormatErrCode, errors.New(SearchTradeResponseDataEncryptFormatErrMessage)
	}*/
	//searchTradeRsp.EncryptRsp = string(searchResultBytes)

	//解密数据
	rspEncryptBytes, err := base64.StdEncoding.DecodeString(searchResult.Encrypt)
	if err != nil {
		logrus.Errorf("org:jd,"+SearchTradeResponseDataDecryptFailedErrMessage+",request:%+v,response:%+v,errCode:%v,err:%v", req, searchResult, SearchTradeResponseDataDecryptFailedErrCode, err.Error())
		return searchTradeRsp, SearchTradeResponseDataDecryptFailedErrCode, errors.New(SearchTradeResponseDataDecryptFailedErrMessage)
	}
	rspEncryptBytes, err = HexString2Bytes(string(rspEncryptBytes))
	rspDecryptBytes, err := TripleEcbDesDecrypt(rspEncryptBytes, desKeyBytes)
	if err != nil {
		logrus.Errorf("org:jd,"+SearchTradeResponseDataDecryptFailedErrMessage+",request:%+v,response:%+v,errCode:%v,err:%v", req, searchResult, SearchTradeResponseDataDecryptFailedErrCode, err.Error())
		return searchTradeRsp, SearchTradeResponseDataDecryptFailedErrCode, errors.New(SearchTradeResponseDataDecryptFailedErrMessage)
	}
	//fmt.Println("search rsp", string(rspDecrypt))
	//searchTradeRsp.DecryptRsp = string(rspDecryptBytes)

	var searchDecryptRsp SearchDecryptRsp
	err = xml.Unmarshal(rspDecryptBytes, &searchDecryptRsp)
	if err != nil {
		logrus.Errorf("org:jd,"+SearchTradeResponseDataDecryptFormatErrMessage+",request:%+v,response:%+v,errCode:%v,err:%v", req, searchResult, SearchTradeResponseDataDecryptFormatErrCode, err.Error())
		return searchTradeRsp, SearchTradeResponseDataDecryptFormatErrCode, errors.New(SearchTradeResponseDataDecryptFormatErrMessage)
	}
	searchTradeRsp.OrderId = searchDecryptRsp.TradeNum

	//签名校验
	if !checkSearchTradeSignature(searchDecryptRsp.Sign, string(rspDecryptBytes), paramMap["public_key"]) {
		logrus.Errorf("org:jd,"+SearchTradeResponseDataSignErrMessage+",request:%+v,response:%+v,errCode:%v,err:%v", req, searchResult, SearchTradeResponseDataSignErrCode)
		return searchTradeRsp, SearchTradeResponseDataSignErrCode, errors.New(SearchTradeResponseDataSignErrMessage)
	}
	switch searchDecryptRsp.Status {
	case TradeCreate:
		searchTradeRsp.Status = SearchTradeWait
	case TradePending:
		searchTradeRsp.Status = SearchTradeWait
	case TradeProcess:
		searchTradeRsp.Status = SearchTradeProcess
	case TradeClosed:
		searchTradeRsp.Status = SearchTradeClosed
	case TradeFailed:
		searchTradeRsp.Status = SearchTradeError
	}

	if searchTradeRsp.Status != TradeProcess {
		return searchTradeRsp, 0, nil
	}

	searchTradeRsp.TradeNo = searchDecryptRsp.TradeNum //该接口不会返回交易流水号，使用请求交易时的订单号
	searchTradeRsp.PaidAt = time.Now().UTC().Format(DateTimeFormatLayout)
	searchTradeRsp.RmbFee = float64(searchDecryptRsp.Amount) / 100

	return searchTradeRsp, 0, err
}

//验证查询交易结果
func checkSearchTradeSignature(sign, decryptRsp, publicKey string) bool {
	//签名字符串截取
	clipStartIndex := strings.Index(decryptRsp, "<sign>")
	clipEndIndex := strings.Index(decryptRsp, "</sign>")
	xmlStart := decryptRsp[0:clipStartIndex]
	xmlEnd := decryptRsp[clipEndIndex+7 : len(decryptRsp)]
	originXml := xmlStart + xmlEnd

	//签名校验
	if sign == "" {
		return false
	}

	signByte, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false
	}
	sha256 := HaSha256(originXml)
	verifySign := VerifyPKCS1v15([]byte(sha256), signByte, []byte(publicKey), crypto.Hash(0))
	if !verifySign {
		fmt.Println("签名校验不通过")
	}
	return verifySign
}

func (jd *Jd) GetSearchTradeConfigCode() []string {
	return []string{
		"merchant",
		"des_key", "trade_way", "private_key", "public_key",
	}
}
