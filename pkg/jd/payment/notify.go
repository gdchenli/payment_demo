package payment

import (
	"encoding/base64"
	"encoding/xml"
	"errors"
	"payment_demo/pkg/jd/util"
	"strings"
)

const (
	NotifySuccessStatus = "2" //支付成功状态
	NotifySuccessCode   = "000000"
)

type Notify struct{}

type NotifyArg struct {
	PublicKey string `json:"public_key"` //公钥
	DesKey    string `json:"des_key"`    //des密钥
}

type NotifyQuery struct {
	XMLName  xml.Name     `xml:"jdpay" json:"-"`
	Version  string       `xml:"version" json:"version"`   //版本号
	Merchant string       `xml:"merchant" json:"merchant"` //商户号
	Result   NotifyResult `xml:"result" json:"result"`     //交易结果
	Encrypt  string       `xml:"encrypt" json:"encrypt"`   //加密信息
}

type NotifyDecrypt struct {
	XMLName   xml.Name      `xml:"jdpay" json:"-"`
	Version   string        `xml:"version" json:"version"`     //版本号
	Merchant  string        `xml:"merchant" json:"merchant"`   //商户号
	Result    NotifyResult  `xml:"result" json:"result"`       //交易结果
	TradeNum  string        `xml:"tradeNum" json:"tradeNum"`   //订单号
	TradeType int           `xml:"tradeType" json:"tradeType"` //交易类型
	Sign      string        `xml:"sign" json:"sign"`           //数据签名
	Amount    int64         `xml:"amount" json:"amount"`       //人民币支付总金额
	OrderId   string        `json:"order_id"`                  //京东交易流水号
	Status    string        `xml:"status" json:"status"`       //交易状态
	PayList   NotifyPayList `xml:"payList" json:"payList"`     //支付方式明细
}

type NotifyResult struct {
	Code string `xml:"code" json:"code"` //交易返回码
	Desc string `xml:"desc" json:"desc"` //返回码信息
}

type NotifyPayList struct {
	Pay []NotifyPay `xml:"pay" json:"pay"`
}

type NotifyPay struct {
	PayType   int    `xml:"payType" json:"payType"`     //支付方式
	Amount    int64  `xml:"amount" json:"amount"`       //交易金额
	Currency  string `xml:"currency" json:"currency"`   //交易币种
	TradeTime string `xml:"tradeTime" json:"tradeTime"` //交易时间
}

type NotifyRsp struct {
	OrderId    string `json:"order_id"`    //订单号
	Status     bool   `json:"status"`      //交易状态，true交易成功 false交易失败
	EncryptRsp string `json:"encrypt_rsp"` //返回的加密数据
	DecryptRsp string `json:"decrypt_rsp"` //返回的解密数据
}

func (notify *Notify) Validate(query string, arg NotifyArg) (notifyRsp NotifyRsp, errCode int, err error) {
	notifyRsp.EncryptRsp = query

	//解析加密的支付机构参数为结构体
	var notifyQuery NotifyQuery
	err = xml.Unmarshal([]byte(query), &notifyQuery)
	if err != nil {
		return notifyRsp, 10201, errors.New("异步通知，加密数据格式错误")
	}

	//解密支付机构参数
	decryptBytes, err := notify.decryptArg(notifyQuery, arg.DesKey)
	if err != nil {
		return notifyRsp, 10202, errors.New("异步通知，解密失败")
	}
	notifyRsp.DecryptRsp = string(decryptBytes)

	//解析解密后的支付机构参数为结构体
	var notifyDecrypt NotifyDecrypt
	err = xml.Unmarshal(decryptBytes, &notifyDecrypt)
	if err != nil {
		return notifyRsp, 10203, errors.New("异步通知，解密后数据格式错误")
	}

	//判断请求结果
	if notifyDecrypt.Result.Code != NotifySuccessCode {
		return notifyRsp, 10104, errors.New("网络请求错误")
	}
	notifyRsp.OrderId = notifyDecrypt.TradeNum

	//校验签名
	if !notify.checkSign(decryptBytes, notifyDecrypt.Sign, arg.PublicKey) {
		return notifyRsp, 10105, errors.New("签名校验失败")
	}

	//交易状态
	if notifyDecrypt.Status != NotifySuccessStatus {
		return notifyRsp, 10106, errors.New("交易状态不正确")
	}
	notifyRsp.Status = true

	return notifyRsp, 0, nil
}

func (notify *Notify) decryptArg(notifyQuery NotifyQuery, desKey string) (decryptBytes []byte, err error) {
	desKeyBytes, err := base64.StdEncoding.DecodeString(desKey)
	if err != nil {
		return nil, err
	}
	encryptBytes, err := base64.StdEncoding.DecodeString(notifyQuery.Encrypt)
	if err != nil {
		return nil, err
	}
	encryptBytes, err = util.HexString2Bytes(string(encryptBytes))
	if err != nil {
		return nil, err
	}
	decryptBytes, err = util.TripleEcbDesDecrypt(encryptBytes, desKeyBytes)
	if err != nil {
		return nil, err
	}

	return decryptBytes, nil
}

func (notify *Notify) checkSign(decryptBytes []byte, sign, publicKey string) bool {
	decrypt := string(decryptBytes)
	clipStartIndex := strings.Index(decrypt, "<sign>")
	clipEndIndex := strings.Index(decrypt, "</sign>")
	xmlStart := decrypt[0:clipStartIndex]
	xmlEnd := decrypt[clipEndIndex+7 : len(decrypt)]
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

	return util.VerifyPKCS1v15([]byte(sha256), signByte, []byte(publicKey))
}
