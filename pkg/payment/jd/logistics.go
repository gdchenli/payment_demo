package jd

import (
	"crypto"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/gdchenli/pay/dialects/jd/util"
	"github.com/gdchenli/pay/pkg/curl"
	"github.com/sirupsen/logrus"
)

const (
	LogisticsUploadSuccessCode = "10000"
)
const (
	LogisticsUploadBuildSignErrCode                    = 10601
	LogisticsUploadBuildSignErrMessage                 = "回传物流信息，签名生成失败"
	LogisticsUploadDesKeyFormatErrCode                 = 10602
	LogisticsUploadDesKeyFormatErrMessage              = "回传物流信息，desKey格式错误"
	LogisticsUploadRequestDataEncryptFailedErrCode     = 10603
	LogisticsUploadRequestDataEncryptFailedErrMessage  = "回传物流信息，请求数据加密失败"
	LogisticsUploadNetErrCode                          = 10604
	LogisticsUploadNetErrMessage                       = "回传物流信息,网络错误"
	LogisticsUploadResponseDataEncryptFormatErrCode    = 10605
	LogisticsUploadResponseDataEncryptFormatErrMessage = "回传物流信息,返回加密数据格式错误"
	LogisticsUploadResponseDataDecryptFailedErrCode    = 10606
	LogisticsUploadResponseDataDecryptFailedErrMessage = "回传物流信息,解密返回数据失败"
	LogisticsUploadResponseDataDecryptFormatErrCode    = 10607
	LogisticsUploadResponseDataDecryptFormatErrMessage = "回传物流信息,解密数据格式错误"
	LogisticsUploadResponseDataSignErrCode             = 10608
	LogisticsUploadResponseDataSignErrMessage          = "回传物流信息,返回数据签名校验错误"
)

type Logistics struct{}

type LogisticsArg struct {
	OrderId          string `json:"order_id"`          //订单编号
	LogisticsNo      string `json:"logistics_no"`      //物流单号
	MerchantName     string `json:"merchant_name"`     //商户名称
	LogisticsCompany string `json:"logistics_company"` //物流公司名称
	Merchant         string `json:"merchant"`          //商户号
	DesKey           string `json:"des_key"`           //desKey
	PrivateKey       string `json:"private_key"`       //私钥
	PublicKey        string `json:"public_key"`        //公钥
	GateWay          string `json:"gate_way"`          //网关地址
}

type LogisticsWithoutSignRequest struct {
	XMLName          xml.Name `xml:"jdpay" json:"-"`
	Version          string   `xml:"version" json:"version"`                   //版本
	Merchant         string   `xml:"merchant" json:"merchant"`                 //商户ID
	TradeNum         string   `xml:"tradeNum" json:"tradeNum"`                 //订单编号
	MerchantName     string   `xml:"merchantName" json:"merchantName"`         //商户名称
	LogisticsNo      string   `xml:"logisticsNo" json:"logisticsNo"`           //物流单号
	LogisticsCompany string   `xml:"logisticsCompany" json:"logisticsCompany"` //物流公司名称
}

type LogisticsWithSignRequest struct {
	XMLName          xml.Name `xml:"jdpay" json:"-"`
	Version          string   `xml:"version" json:"version"`                   //版本
	Merchant         string   `xml:"merchant" json:"merchant"`                 //商户ID
	TradeNum         string   `xml:"tradeNum" json:"tradeNum"`                 //订单编号
	MerchantName     string   `xml:"merchantName" json:"merchantName"`         //商户名称
	LogisticsNo      string   `xml:"logisticsNo" json:"logisticsNo"`           //物流单号
	LogisticsCompany string   `xml:"logisticsCompany" json:"logisticsCompany"` //物流公司名称
	Sign             string   `xml:"sign" json:"sign"`                         //签名
}

type LogisticsWithEncrypt struct {
	XMLName  xml.Name `xml:"jdpay" json:"-"`
	Version  string   `xml:"version" json:"version"`   //版本
	Merchant string   `xml:"merchant" json:"merchant"` //商户号
	Encrypt  string   `xml:"encrypt" json:"encrypt"`   //加密数据
}

type LogisticsResult struct {
	XMLName  xml.Name     `xml:"jdpay" json:"-"`
	Version  string       `xml:"version" json:"version"`   //版本号
	Merchant string       `xml:"merchant" json:"merchant"` //商户号
	Result   LogisticsRsp `xml:"result" json:"result"`     //结果
	Encrypt  string       `xml:"encrypt" json:"encrypt"`   //加密信息
}

type LogisticsResultRsp struct {
	Code string `xml:"code" json:"code"` //交易返回码
	Desc string `xml:"desc" json:"desc"` //返回码信息
}

type LogisticsDecryptRsp struct {
	XMLName  xml.Name           `xml:"jdpay" json:"-"`
	Version  string             `xml:"version" json:"version"`   //版本
	Merchant string             `xml:"merchant" json:"merchant"` //商户号
	Result   LogisticsResultRsp `xml:"result" json:"result"`     //结果
	Sign     string             `xml:"sign" json:"sign"`         //数据签名
}

type LogisticsRsp struct {
	Status     bool   `json:"status"`      //交易状态
	OrderId    string `json:"order_id"`    //订单号
	EncryptRsp string `json:"encrypt_rsp"` //返回的加密数据
	DecryptRsp string `json:"decrypt_rsp"` //返回的解密数据
	EncryptRes string `json:"encrypt_res"` //请求的加密数据
	DecryptRes string `json:"decrypt_res"` //请求的未加密数据
}

func (l *Logistics) Upload(arg LogisticsArg) (logisticsRsp LogisticsRsp, errCode int, err error) {
	logisticsRsp.OrderId = arg.OrderId

	logisticsWithoutSignRequest := LogisticsWithoutSignRequest{
		Version:          Version,
		Merchant:         arg.Merchant,
		TradeNum:         arg.OrderId,
		MerchantName:     arg.MerchantName,
		LogisticsNo:      arg.LogisticsNo,
		LogisticsCompany: arg.LogisticsCompany,
	}
	xmlBytes, err := xml.Marshal(logisticsWithoutSignRequest)
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
	sha256 := util.HaSha256(xmlStr)
	signBytes, err := util.SignPKCS1v15([]byte(sha256), []byte(arg.PrivateKey), crypto.Hash(0))
	if err != nil {
		logrus.Errorf(LogisticsUploadBuildSignErrMessage+",request:%+v,errCode:%v,err:%v", arg, LogisticsUploadBuildSignErrCode, err.Error())
		return logisticsRsp, LogisticsUploadBuildSignErrCode, errors.New(LogisticsUploadBuildSignErrMessage)
	}
	sign := base64.StdEncoding.EncodeToString(signBytes)
	logisticsWithSignRequest := LogisticsWithSignRequest{
		Version:          logisticsWithoutSignRequest.Version,
		Merchant:         logisticsWithoutSignRequest.Merchant,
		TradeNum:         logisticsWithoutSignRequest.TradeNum,
		MerchantName:     logisticsWithoutSignRequest.MerchantName,
		LogisticsNo:      logisticsWithoutSignRequest.LogisticsNo,
		LogisticsCompany: logisticsWithoutSignRequest.LogisticsCompany,
		Sign:             sign,
	}
	xmlBytes, err = xml.Marshal(logisticsWithSignRequest)
	xmlStr = strings.TrimRight(xml.Header, "\n") + string(xmlBytes)
	//fmt.Println("with sign xml", xmlStr)

	desKey, err := base64.StdEncoding.DecodeString(arg.DesKey)
	if err != nil {
		logrus.Errorf(LogisticsUploadDesKeyFormatErrMessage+",request:%+v,errCode:%v,err:%v", arg, LogisticsUploadDesKeyFormatErrCode, err.Error())
		return logisticsRsp, LogisticsUploadDesKeyFormatErrCode, errors.New(LogisticsUploadDesKeyFormatErrMessage)
	}
	encryptBytes, err := util.TripleEcbDesEncrypt([]byte(xmlStr), desKey)
	if err != nil {
		logrus.Errorf(LogisticsUploadRequestDataEncryptFailedErrMessage+",query:%+v,errCode:%v,err:%v", arg, LogisticsUploadRequestDataEncryptFailedErrCode, err.Error())
		return logisticsRsp, LogisticsUploadRequestDataEncryptFailedErrCode, errors.New(LogisticsUploadRequestDataEncryptFailedErrMessage)
	}
	reqEncrypt := util.DecimalByteSlice2HexString(encryptBytes)
	reqEncrypt = base64.StdEncoding.EncodeToString([]byte(reqEncrypt))
	logisticsWithEncrypt := LogisticsWithEncrypt{
		Version:  Version,
		Merchant: arg.Merchant,
		Encrypt:  reqEncrypt,
	}
	xmlBytes, err = xml.Marshal(logisticsWithEncrypt)
	xmlStr = strings.TrimRight(xml.Header, "\n") + string(xmlBytes)
	//fmt.Println("with 3des xml", xmlStr)

	var logisticsResult LogisticsResult
	playLoad := strings.NewReader(xmlStr)
	err = curl.PostXML(arg.GateWay, &logisticsResult, playLoad)
	if err != nil {
		//fmt.Println("err", err)
		logrus.Errorf(LogisticsUploadNetErrMessage+",request:%+v,errCode:%v,err:%v", arg, LogisticsUploadNetErrCode, err.Error())
		return logisticsRsp, LogisticsUploadNetErrCode, errors.New(LogisticsUploadNetErrMessage)
	}

	logisticsResultBytes, err := xml.Marshal(logisticsResult)
	if err != nil {
		logrus.Errorf(LogisticsUploadResponseDataEncryptFormatErrMessage+",request:%+v,response:%+v,errCode:%v,err:%v", arg, logisticsResult, LogisticsUploadResponseDataEncryptFormatErrCode, err.Error())
		return logisticsRsp, LogisticsUploadResponseDataEncryptFormatErrCode, errors.New(LogisticsUploadResponseDataEncryptFormatErrMessage)
	}
	logisticsRsp.EncryptRsp = string(logisticsResultBytes)

	//解密数据
	rspEncryptBytes, err := base64.StdEncoding.DecodeString(logisticsResult.Encrypt)
	if err != nil {
		logrus.Errorf(LogisticsUploadResponseDataDecryptFailedErrMessage+",request:%+v,response:%+v,errCode:%v,err:%v", arg, logisticsResult, LogisticsUploadResponseDataDecryptFailedErrCode, err.Error())
		return logisticsRsp, LogisticsUploadResponseDataDecryptFailedErrCode, errors.New(LogisticsUploadResponseDataDecryptFailedErrMessage)
	}
	rspEncryptBytes, err = util.HexString2Bytes(string(rspEncryptBytes))
	rspDecryptBytes, err := util.TripleEcbDesDecrypt(rspEncryptBytes, desKey)
	if err != nil {
		logrus.Errorf(LogisticsUploadResponseDataDecryptFailedErrMessage+",request:%+v,response:%+v,errCode:%v,err:%v", arg, logisticsResult, LogisticsUploadResponseDataDecryptFailedErrCode, err.Error())
		return logisticsRsp, LogisticsUploadResponseDataDecryptFailedErrCode, errors.New(LogisticsUploadResponseDataDecryptFailedErrMessage)
	}
	//fmt.Println("logistics rsp", string(rspDecrypt))
	logisticsRsp.DecryptRsp = string(rspDecryptBytes)

	var logisticsDecryptRsp LogisticsDecryptRsp
	err = xml.Unmarshal(rspDecryptBytes, &logisticsDecryptRsp)
	if err != nil {
		logrus.Errorf(LogisticsUploadResponseDataDecryptFormatErrMessage+",request:%+v,response:%+v,errCode:%v,err:%v", arg, logisticsResult, LogisticsUploadResponseDataDecryptFormatErrCode, err.Error())
		return logisticsRsp, LogisticsUploadResponseDataDecryptFormatErrCode, errors.New(LogisticsUploadResponseDataDecryptFormatErrMessage)
	}

	//签名校验
	if !l.checkSignature(logisticsDecryptRsp.Sign, logisticsRsp.DecryptRsp, arg.PublicKey) {
		logrus.Errorf(LogisticsUploadResponseDataSignErrMessage+",request:%+v,response:%+v,errCode:%v,err:%v", arg, logisticsRsp, LogisticsUploadResponseDataSignErrCode)
		return logisticsRsp, LogisticsUploadResponseDataSignErrCode, errors.New(LogisticsUploadResponseDataSignErrMessage)
	}

	if logisticsDecryptRsp.Result.Code == LogisticsUploadSuccessCode {
		logisticsRsp.Status = true

	}

	return logisticsRsp, 0, nil
}

//验证结果
func (l *Logistics) checkSignature(sign, decryptRsp, publicKey string) bool {
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
	sha256 := util.HaSha256(originXml)
	verifySign := util.VerifyPKCS1v15([]byte(sha256), signByte, []byte(publicKey), crypto.Hash(0))
	if !verifySign {
		fmt.Println("签名校验不通过")
	}
	return verifySign
}
