package jd

import (
	"crypto"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"payment_demo/api/response"
	"payment_demo/internal/common/request"
	"payment_demo/pkg/curl"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	ClosedTradeSuccessStatus = 1
)

const (
	CloseTradeBuildSignErrCode                    = 10401
	CloseTradeBuildSignErrMessage                 = "查询交易流水，签名生成失败"
	CloseTradeDesKeyFormatErrCode                 = 10402
	CloseTradeDesKeyFormatErrMessage              = "查询交易流水，desKey格式错误"
	CloseTradeRequestDataEncryptFailedErrCode     = 10403
	CloseTradeRequestDataEncryptFailedErrMessage  = "查询交易流水，请求数据加密失败"
	CloseTradeNetErrCode                          = 10404
	CloseTradeNetErrMessage                       = "关闭交易流水,网络错误"
	CloseTradeResponseDataEncryptFormatErrCode    = 10405
	CloseTradeResponseDataEncryptFormatErrMessage = "关闭交易流水,返回加密数据格式错误"
	CloseTradeResponseDataDecryptFailedErrCode    = 10406
	CloseTradeResponseDataDecryptFailedErrMessage = "关闭交易流水,解密返回数据失败"
	CloseTradeResponseDataDecryptFormatErrCode    = 10407
	CloseTradeResponseDataDecryptFormatErrMessage = "关闭交易流水,解密数据格式错误"
	CloseTradeResponseDataSignErrCode             = 10408
	CloseTradeResponseDataSignErrMessage          = "关闭交易流水,返回数据签名校验错误"
)

type Close struct{}

type CloseWithoutSignRequest struct {
	XMLName   xml.Name `xml:"jdpay" json:"-"`
	Version   string   `xml:"version" json:"version"`     //版本
	Merchant  string   `xml:"merchant" json:"merchant"`   //商户号
	TradeNum  string   `xml:"tradeNum" json:"tradeNum"`   //订单编号
	OTradeNum string   `xml:"oTradeNum" json:"oTradeNum"` //原交易流水号
	Amount    int64    `xml:"amount" json:"amount"`       //金额
	Currency  string   `xml:"currency" json:"currency"`   //币种
}

type CloseWithSignRequest struct {
	XMLName   xml.Name `xml:"jdpay" json:"-"`
	Version   string   `xml:"version" json:"version"`     //版本
	Merchant  string   `xml:"merchant" json:"merchant"`   //商户号
	TradeNum  string   `xml:"tradeNum" json:"tradeNum"`   //订单编号
	OTradeNum string   `xml:"oTradeNum" json:"oTradeNum"` //原交易流水号
	Amount    int64    `xml:"amount" json:"amount"`       //金额
	Currency  string   `xml:"currency" json:"currency"`   //币种
	Sign      string   `xml:"sign" json:"sign"`           //签名
}

type CloseWithEncrypt struct {
	XMLName  xml.Name `xml:"jdpay" json:"-"`
	Version  string   `xml:"version" json:"version"`   //版本
	Merchant string   `xml:"merchant" json:"merchant"` //商户号
	Encrypt  string   `xml:"encrypt" json:"encrypt"`   //加密数据
}

type CloseResult struct {
	XMLName  xml.Name            `xml:"jdpay" json:"-"`
	Version  string              `xml:"version" json:"version"`   //版本号
	Merchant string              `xml:"merchant" json:"merchant"` //商户号
	Result   CloseResultResponse `xml:"result" json:"result"`     //交易结果
	Encrypt  string              `xml:"encrypt" json:"encrypt"`   //加密信息
}

type CloseResultResponse struct {
	Code string `xml:"code" json:"code"` //交易返回码
	Desc string `xml:"desc" json:"desc"` //返回码信息
}

type CloseDecryptRsp struct {
	XMLName   xml.Name       `xml:"jdpay" json:"-"`
	Merchant  string         `xml:"merchant" json:"merchant"`   //商户号
	TradeNum  string         `xml:"tradeNum" json:"tradeNum"`   //订单编号
	TradeType string         `xml:"tradeType" json:"tradeType"` //交易类型
	Result    CloseResultRsp `xml:"result" json:"result"`       //交易结果
	Sign      string         `xml:"sign" json:"sign"`           //数据签名
	OTradeNum string         `xml:"oTradeNum" json:"oTradeNum"` //原交易流水号
	Amount    int64          `xml:"amount" json:"amount"`       //人民币支付总金额
	Currency  string         `xml:"currency" json:"currency"`   //交易币种
	TradeTime string         `xml:"tradeTime" json:"tradeTime"` //交易时间
	Status    int            `xml:"status" json:"status"`       //交易状态
}

type CloseResultRsp struct {
	Code string `xml:"code" json:"code"` //交易返回码
	Desc string `xml:"desc" json:"desc"` //返回码信息
}

func (jd *Jd) CloseTrade(configParamMap map[string]string, req request.CloseTradeReq) (closeTradeRsp response.CloseTradeRsp, errCode int, err error) {
	totalFeeStr := fmt.Sprintf("%.f", req.TotalFee*100)
	totalFee, _ := strconv.ParseInt(totalFeeStr, 10, 64)
	closedWithoutSignRequest := CloseWithoutSignRequest{
		Version:   Version,
		Merchant:  configParamMap["merchant"],
		TradeNum:  req.OrderId + "jd",
		OTradeNum: req.OrderId,
		Amount:    totalFee,
		Currency:  req.Currency,
	}

	xmlBytes, err := xml.Marshal(closedWithoutSignRequest)
	xmlStr := xml.Header + string(xmlBytes)
	xmlStr = strings.Replace(xmlStr, "\r", "", -1)
	xmlStr = strings.Replace(xmlStr, "\n", "", -1)
	xmlStr = strings.Replace(xmlStr, "\t", "", -1)
	reg, _ := regexp.Compile(">\\s+<")
	xmlStr = reg.ReplaceAllString(xmlStr, "><")
	reg, _ = regexp.Compile("\\s+\\/>")
	xmlStr = reg.ReplaceAllString(xmlStr, "/>")
	fmt.Println("without sign xml", xmlStr)

	//生成签名
	sha256 := HaSha256(xmlStr)
	signBytes, err := SignPKCS1v15([]byte(sha256), []byte(configParamMap["private_key"]), crypto.Hash(0))
	if err != nil {
		logrus.Errorf(CloseTradeBuildSignErrMessage+",request:%+v,errCode:%v,err:%v", req, CloseTradeBuildSignErrCode, err.Error())
		return closeTradeRsp, CloseTradeBuildSignErrCode, errors.New(CloseTradeBuildSignErrMessage)
	}
	sign := base64.StdEncoding.EncodeToString(signBytes)
	closedWithSignRequest := CloseWithSignRequest{
		Version:   closedWithoutSignRequest.Version,
		Merchant:  closedWithoutSignRequest.Merchant,
		TradeNum:  closedWithoutSignRequest.TradeNum,
		OTradeNum: closedWithoutSignRequest.OTradeNum,
		Amount:    totalFee,
		Currency:  req.Currency,
		Sign:      sign,
	}
	xmlBytes, err = xml.Marshal(closedWithSignRequest)
	xmlStr = strings.TrimRight(xml.Header, "\n") + string(xmlBytes)
	//closeTradeRsp.DecryptRes = xmlStr

	desKeyBytes, err := base64.StdEncoding.DecodeString(configParamMap["des_key"])
	if err != nil {
		logrus.Errorf(CloseTradeDesKeyFormatErrMessage+",request:%+v,errCode:%v,err:%v", req, CloseTradeDesKeyFormatErrCode, err.Error())
		return closeTradeRsp, CloseTradeDesKeyFormatErrCode, errors.New(CloseTradeDesKeyFormatErrMessage)
	}
	encryptBytes, err := TripleEcbDesEncrypt([]byte(xmlStr), desKeyBytes)
	if err != nil {
		logrus.Errorf(CloseTradeRequestDataEncryptFailedErrMessage+",request:%+v,errCode:%v,err:%v", req, CloseTradeRequestDataEncryptFailedErrCode, err.Error())
		return closeTradeRsp, CloseTradeRequestDataEncryptFailedErrCode, errors.New(CloseTradeRequestDataEncryptFailedErrMessage)
	}
	reqEncrypt := DecimalByteSlice2HexString(encryptBytes)
	reqEncrypt = base64.StdEncoding.EncodeToString([]byte(reqEncrypt))
	closedWithEncrypt := CloseWithEncrypt{
		Version:  Version,
		Merchant: configParamMap["merchant"],
		Encrypt:  reqEncrypt,
	}
	xmlBytes, err = xml.Marshal(closedWithEncrypt)
	xmlStr = strings.TrimRight(xml.Header, "\n") + string(xmlBytes)
	fmt.Println("with 3des xml", xmlStr)
	//closeTradeRsp.EncryptRes = xmlStr

	var closeResult CloseResult
	playLoad := strings.NewReader(xmlStr)
	err = curl.PostXML(configParamMap["close_way"], &closeResult, playLoad)
	if err != nil {
		logrus.Errorf(CloseTradeNetErrMessage+",request:%+v,errCode:%v,err:%v", req, CloseTradeNetErrCode, err.Error())
		return closeTradeRsp, CloseTradeNetErrCode, errors.New(CloseTradeNetErrMessage)
	}
	//fmt.Printf("closedResult:%+v\n", closedResult)
	/*closedResultBytes, err := xml.Marshal(closedResult)
	if err != nil {
		logrus.Errorf(CloseTradeResponseDataEncryptFormatErrMessage+",request:%+v,response:%v,errCode:%v,err:%v", arg, closedResult, CloseTradeResponseDataEncryptFormatErrCode, err.Error())
		return closeTradeRsp, CloseTradeResponseDataEncryptFormatErrCode, errors.New(CloseTradeResponseDataEncryptFormatErrMessage)
	}*/
	//closeTradeRsp.EncryptRsp = string(closedResultBytes)

	//解密数据
	rspEncryptBytes, err := base64.StdEncoding.DecodeString(closeResult.Encrypt)
	if err != nil {
		logrus.Errorf(CloseTradeResponseDataDecryptFailedErrMessage+",request:%+v,response:%v,errCode:%v,err:%v", req, closeResult, CloseTradeResponseDataDecryptFailedErrCode, err.Error())
		return closeTradeRsp, CloseTradeResponseDataDecryptFailedErrCode, errors.New(CloseTradeResponseDataDecryptFailedErrMessage)
	}
	rspEncryptBytes, err = HexString2Bytes(string(rspEncryptBytes))
	rspDecryptBytes, err := TripleEcbDesDecrypt(rspEncryptBytes, desKeyBytes)
	if err != nil {
		logrus.Errorf(CloseTradeResponseDataDecryptFailedErrMessage+",request:%+v,response:%v,errCode:%v,err:%v", req, closeResult, CloseTradeResponseDataDecryptFailedErrCode, err.Error())
		return closeTradeRsp, CloseTradeResponseDataDecryptFailedErrCode, errors.New(CloseTradeResponseDataDecryptFailedErrMessage)
	}
	//fmt.Println("search rsp", string(rspDecrypt))
	//closeTradeRsp.DecryptRsp = string(rspDecryptBytes)

	var closeDecryptRsp CloseDecryptRsp
	err = xml.Unmarshal(rspDecryptBytes, &closeDecryptRsp)
	if err != nil {
		logrus.Errorf(CloseTradeResponseDataDecryptFormatErrMessage+",request:%+v,response:%v,errCode:%v,err:%v", req, closeResult, CloseTradeResponseDataDecryptFormatErrCode, err.Error())
		return closeTradeRsp, CloseTradeResponseDataDecryptFormatErrCode, errors.New(CloseTradeResponseDataDecryptFormatErrMessage)
	}
	closeTradeRsp.OrderId = closeDecryptRsp.TradeNum

	//签名校验
	if !checkCloseTradeSignature(closeDecryptRsp.Sign, string(rspDecryptBytes), configParamMap["public_key"]) {
		logrus.Errorf(CloseTradeResponseDataSignErrMessage+",request:%+v,response:%v,errCode:%v,err:%v", req, closeResult, CloseTradeResponseDataSignErrCode)
		return closeTradeRsp, CloseTradeResponseDataSignErrCode, errors.New(CloseTradeResponseDataSignErrMessage)
	}

	if closeDecryptRsp.Status == ClosedTradeSuccessStatus {
		closeTradeRsp.Status = true
	}

	return closeTradeRsp, 0, nil
}

//验证查询交易结果
func checkCloseTradeSignature(sign, decryptRsp, publicKey string) bool {
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

func (jd *Jd) GetCloseTradeConfigCode() []string {
	return []string{
		"merchant",
		"des_key", "private_key", "public_key", "close_way",
	}
}
