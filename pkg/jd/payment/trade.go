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
	CustomTradeType = "0" //交易类型
)

const (
	TradeCreate  = "0" //交易创建
	TradePending = "1" //交易处理中
	TradeProcess = "2" //交易成功
	TradeClosed  = "3" //交易关闭
	TradeFailed  = "4" //交易失败
)

type Trade struct{}

type TradeArg struct {
	Merchant   string `json:"merchant"`   //商户号
	TradeNum   string `json:"tradeNum"`   //订单编号
	OTradeNum  string `json:"oTradeNum"`  //原交易流水号
	DesKey     string `json:"signKey"`    //desKey
	PrivateKey string `json:"privateKey"` //私钥
	PublicKey  string `json:"publicKey"`  //公钥
	GateWay    string `json:"gateWay"`    //网关地址
}

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

type TradeRsp struct {
	Status     string `json:"status"`      //交易状态
	OrderId    string `json:"order_id"`    //订单号
	TradeNo    string `json:"trade_no"`    //支付机构交易流水号
	EncryptRsp string `json:"encrypt_rsp"` //返回的加密数据
	DecryptRsp string `json:"decrypt_rsp"` //返回的解密数据
	EncryptRes string `json:"encrypt_res"` //请求的加密数据
	DecryptRes string `json:"decrypt_res"` //请求的未加密数据
}

func (trade *Trade) Search(arg TradeArg) (tradeRsp TradeRsp, errCode int, err error) {
	searchWithoutSignRequest := SearchWithoutSignRequest{
		Version:   Version,
		Merchant:  arg.Merchant,
		TradeNum:  arg.TradeNum,
		OTradeNum: arg.OTradeNum,
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
	sha256 := util.HaSha256(xmlStr)
	signBytes, err := util.SignPKCS1v15([]byte(sha256), []byte(arg.PrivateKey))
	if err != nil {
		return tradeRsp, 10301, errors.New("查询交易流水，签名生成失败")
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
	tradeRsp.DecryptRes = xmlStr

	desKeyBytes, err := base64.StdEncoding.DecodeString(arg.DesKey)
	if err != nil {
		return tradeRsp, 10302, errors.New("查询交易流水，desKey格式错误")
	}
	encryptBytes, err := util.TripleEcbDesEncrypt([]byte(xmlStr), desKeyBytes)
	if err != nil {
		return tradeRsp, 10303, err
	}
	reqEncrypt := util.DecimalByteSlice2HexString(encryptBytes)
	reqEncrypt = base64.StdEncoding.EncodeToString([]byte(reqEncrypt))
	searchWithEncrypt := SearchWithEncrypt{
		Version:  Version,
		Merchant: arg.Merchant,
		Encrypt:  reqEncrypt,
	}
	xmlBytes, err = xml.Marshal(searchWithEncrypt)
	xmlStr = strings.TrimRight(xml.Header, "\n") + string(xmlBytes)
	//fmt.Println("with 3des xml", xmlStr)
	tradeRsp.EncryptRes = xmlStr

	var searchResult SearchResult
	playLoad := strings.NewReader(xmlStr)
	err = curl.PostXML(arg.GateWay, &searchResult, playLoad)
	if err != nil {
		return tradeRsp, 10504, errors.New("查询交易流水,网络错误")
	}
	searchResultBytes, err := xml.Marshal(searchResult)
	if err != nil {
		return tradeRsp, 10505, errors.New("查询交易流水,返回结果信格式错误")
	}
	tradeRsp.EncryptRsp = string(searchResultBytes)

	//解密数据
	rspEncryptBytes, err := base64.StdEncoding.DecodeString(searchResult.Encrypt)
	if err != nil {
		return tradeRsp, 10506, errors.New("查询交易流水,解密数据失败")
	}
	rspEncryptBytes, err = util.HexString2Bytes(string(rspEncryptBytes))
	rspDecryptBytes, err := util.TripleEcbDesDecrypt(rspEncryptBytes, desKeyBytes)
	if err != nil {
		return tradeRsp, 10506, errors.New("查询交易流水,解密数据失败")
	}
	//fmt.Println("search rsp", string(rspDecrypt))
	tradeRsp.DecryptRsp = string(rspDecryptBytes)

	var searchDecryptRsp SearchDecryptRsp
	err = xml.Unmarshal(rspDecryptBytes, &searchDecryptRsp)
	if err != nil {
		return tradeRsp, 10506, errors.New("查询交易流水,解密数据格式错误")
	}
	tradeRsp.OrderId = searchDecryptRsp.TradeNum

	//签名校验
	if !trade.checkSignature(searchDecryptRsp.Sign, tradeRsp.DecryptRsp, arg.PublicKey) {
		return tradeRsp, 10507, errors.New("查询交易流水,签名校验错误")
	}

	tradeRsp.Status = searchDecryptRsp.Status
	tradeRsp.TradeNo = searchDecryptRsp.TradeNum //该接口不会返回交易流水号，使用请求交易时的订单号

	return tradeRsp, 0, err
}

//验证查询交易结果
func (trade *Trade) checkSignature(sign, decryptRsp, publicKey string) bool {
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
