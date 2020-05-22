package payment

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"payment_demo/pkg/jd/util"
	"payment_demo/tools/curl"
	"regexp"
	"strings"
)

const (
	JdClosedTradeSuccessStatus = 1
)

type Closed struct{}

type ClosedArg struct {
	Merchant   string `json:"merchant"`   //商户ID
	TradeNum   string `json:"tradeNum"`   //订单编号
	OTradeNum  string `json:"oTradeNum"`  //原交易流水号
	Amount     int64  `json:"amount"`     //交易金额
	Currency   string `json:"currency"`   //交易币种
	DesKey     string `json:"signKey"`    //desKey
	PrivateKey string `json:"privateKey"` //私钥
	PublicKey  string `json:"publicKey"`  //公钥
	GateWay    string `json:"gate_way"`   //网关地址
}

type ClosedWithoutSignRequest struct {
	XMLName   xml.Name `xml:"jdpay" json:"-"`
	Version   string   `xml:"version" json:"version"`     //版本
	Merchant  string   `xml:"merchant" json:"merchant"`   //商户号
	TradeNum  string   `xml:"tradeNum" json:"tradeNum"`   //订单编号
	OTradeNum string   `xml:"oTradeNum" json:"oTradeNum"` //原交易流水号
	Amount    int64    `xml:"amount" json:"amount"`       //金额
	Currency  string   `xml:"currency" json:"currency"`   //币种
}

type ClosedWithSignRequest struct {
	XMLName   xml.Name `xml:"jdpay" json:"-"`
	Version   string   `xml:"version" json:"version"`     //版本
	Merchant  string   `xml:"merchant" json:"merchant"`   //商户号
	TradeNum  string   `xml:"tradeNum" json:"tradeNum"`   //订单编号
	OTradeNum string   `xml:"oTradeNum" json:"oTradeNum"` //原交易流水号
	Amount    int64    `xml:"amount" json:"amount"`       //金额
	Currency  string   `xml:"currency" json:"currency"`   //币种
	Sign      string   `xml:"sign" json:"sign"`           //签名
}

type ClosedWithEncrypt struct {
	XMLName  xml.Name `xml:"jdpay" json:"-"`
	Version  string   `xml:"version" json:"version"`   //版本
	Merchant string   `xml:"merchant" json:"merchant"` //商户号
	Encrypt  string   `xml:"encrypt" json:"encrypt"`   //加密数据
}

type ClosedResult struct {
	XMLName  xml.Name             `xml:"jdpay" json:"-"`
	Version  string               `xml:"version" json:"version"`   //版本号
	Merchant string               `xml:"merchant" json:"merchant"` //商户号
	Result   ClosedResultResponse `xml:"result" json:"result"`     //交易结果
	Encrypt  string               `xml:"encrypt" json:"encrypt"`   //加密信息
}

type ClosedResultResponse struct {
	Code string `xml:"code" json:"code"` //交易返回码
	Desc string `xml:"desc" json:"desc"` //返回码信息
}

type ClosedDecryptRsp struct {
	XMLName   xml.Name        `xml:"jdpay" json:"-"`
	Merchant  string          `xml:"merchant" json:"merchant"`   //商户号
	TradeNum  string          `xml:"tradeNum" json:"tradeNum"`   //订单编号
	TradeType string          `xml:"tradeType" json:"tradeType"` //交易类型
	Result    ClosedResultRsp `xml:"result" json:"result"`       //交易结果
	Sign      string          `xml:"sign" json:"sign"`           //数据签名
	OTradeNum string          `xml:"oTradeNum" json:"oTradeNum"` //原交易流水号
	Amount    int64           `xml:"amount" json:"amount"`       //人民币支付总金额
	Currency  string          `xml:"currency" json:"currency"`   //交易币种
	TradeTime string          `xml:"tradeTime" json:"tradeTime"` //交易时间
	Status    int             `xml:"status" json:"status"`       //交易状态
}

type ClosedResultRsp struct {
	Code string `xml:"code" json:"code"` //交易返回码
	Desc string `xml:"desc" json:"desc"` //返回码信息
}

type ClosedRsp struct {
	Status     bool   `json:"status"`      //交易关闭状态
	OrderId    string `json:"order_id"`    //订单号
	EncryptRsp string `json:"encrypt_rsp"` //返回的加密数据
	DecryptRsp string `json:"decrypt_rsp"` //返回的解密数据
	EncryptRes string `json:"encrypt_res"` //请求的加密数据
	DecryptRes string `json:"decrypt_res"` //请求的未加密数据
}

func (closed *Closed) Trade(arg ClosedArg) (closedRsp ClosedRsp, errCode int, err error) {
	closedWithoutSignRequest := ClosedWithoutSignRequest{
		Version:   Version,
		Merchant:  arg.Merchant,
		TradeNum:  arg.TradeNum,
		OTradeNum: arg.OTradeNum,
		Amount:    arg.Amount,
		Currency:  arg.Currency,
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
	//fmt.Println("without sign xml", xmlStr)

	//生成签名
	sha256 := util.HaSha256(xmlStr)
	signBytes, err := util.SignPKCS1v15([]byte(sha256), []byte(arg.PrivateKey))
	if err != nil {
		return closedRsp, 10401, errors.New("关闭交易，签名生成失败")
	}
	sign := base64.StdEncoding.EncodeToString(signBytes)
	closedWithSignRequest := ClosedWithSignRequest{
		Version:   closedWithoutSignRequest.Version,
		Merchant:  closedWithoutSignRequest.Merchant,
		TradeNum:  closedWithoutSignRequest.TradeNum,
		OTradeNum: closedWithoutSignRequest.OTradeNum,
		Amount:    arg.Amount,
		Currency:  arg.Currency,
		Sign:      sign,
	}
	xmlBytes, err = xml.Marshal(closedWithSignRequest)
	xmlStr = strings.TrimRight(xml.Header, "\n") + string(xmlBytes)
	closedRsp.DecryptRes = xmlStr

	desKeyBytes, err := base64.StdEncoding.DecodeString(arg.DesKey)
	if err != nil {
		return closedRsp, 10402, errors.New("关闭交易，请求数据加密失败")
	}
	encryptBytes, err := util.TripleEcbDesEncrypt([]byte(xmlStr), desKeyBytes)
	if err != nil {
		return closedRsp, 10403, errors.New("关闭交易，请求数据加密失败")
	}
	reqEncrypt := util.DecimalByteSlice2HexString(encryptBytes)
	reqEncrypt = base64.StdEncoding.EncodeToString([]byte(reqEncrypt))
	closedWithEncrypt := ClosedWithEncrypt{
		Version:  Version,
		Merchant: arg.Merchant,
		Encrypt:  reqEncrypt,
	}
	xmlBytes, err = xml.Marshal(closedWithEncrypt)
	xmlStr = strings.TrimRight(xml.Header, "\n") + string(xmlBytes)
	//fmt.Println("with 3des xml", xmlStr)
	closedRsp.EncryptRes = xmlStr

	var closedResult ClosedResult
	playLoad := strings.NewReader(xmlStr)
	err = curl.PostXML(arg.GateWay, &closedResult, playLoad)
	if err != nil {
		return closedRsp, 10404, errors.New("关闭交易，网络错误")
	}
	//fmt.Printf("closedResult:%+v\n", closedResult)
	closedResultBytes, err := xml.Marshal(closedResult)
	if err != nil {
		return closedRsp, 10404, errors.New("关闭交易，返回数据格式错误")
	}
	closedRsp.EncryptRsp = string(closedResultBytes)

	//解密数据
	rspEncryptBytes, err := base64.StdEncoding.DecodeString(closedResult.Encrypt)
	if err != nil {
		return closedRsp, 10405, errors.New("关闭交易，返回数据格解密失败")
	}
	rspEncryptBytes, err = util.HexString2Bytes(string(rspEncryptBytes))
	rspDecryptBytes, err := util.TripleEcbDesDecrypt(rspEncryptBytes, desKeyBytes)
	if err != nil {
		return closedRsp, 10405, errors.New("关闭交易，返回数据格解密失败")
	}
	//fmt.Println("search rsp", string(rspDecrypt))
	closedRsp.DecryptRsp = string(rspDecryptBytes)

	var closedDecryptRsp ClosedDecryptRsp
	err = xml.Unmarshal(rspDecryptBytes, &closedDecryptRsp)
	if err != nil {
		return closedRsp, 10407, errors.New("关闭交易，解密数据格式错误")
	}
	closedRsp.OrderId = closedDecryptRsp.TradeNum

	//签名校验
	if !closed.checkSignature(closedDecryptRsp.Sign, closedRsp.DecryptRsp, arg.PublicKey) {
		return closedRsp, 10408, errors.New("关闭交易,签名校验错误")
	}

	if closedDecryptRsp.Status == JdClosedTradeSuccessStatus {
		closedRsp.Status = true
	}

	return closedRsp, 0, nil
}

//验证查询交易结果
func (closed *Closed) checkSignature(sign, decryptRsp, publicKey string) bool {
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
	verifySign := util.VerifyPKCS1v15([]byte(sha256), signByte, []byte(publicKey))
	if !verifySign {
		fmt.Println("签名校验不通过")
	}
	return verifySign
}
