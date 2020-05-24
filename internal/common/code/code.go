package code

const (
	AmountFormatErrCode               = 10001
	AmountFormatErrMessage            = "金额转换异常"
	PrivateKeyNotExitsErrCode         = 10002
	PrivateKeyNotExitsErrMessage      = "私钥文件不存在"
	PrivateKeyContentErrCode          = 10002
	PrivateKeyContentErrMessage       = "私钥内容错误"
	PublicKeyNotExitsErrCode          = 10003
	PublicKeyNotExitsErrMessage       = "公钥文件不存在"
	PublicKeyContentErrCode           = 10004
	PublicKeyContentErrMessage        = "公钥内容错误"
	GateWayNotExitsErrCode            = 10005
	GateWayNotExitsErrMessage         = "网关不存在"
	MerchantNotExitsErrCode           = 10006
	MerchantNotExitsErrMessage        = "商户Id不存在"
	DesKeyNotExitsErrCode             = 10007
	DesKeyNotExitsErrMessage          = "des key不存在"
	ExpireTimeNotExitsErrCode         = 10008
	ExpireTimeNotExitsErrMessage      = "订单交易过期时间不存在"
	TransCurrencyNotExitsErrCode      = 10009
	TransCurrencyNotExitsErrMessage   = "订单结算币种不存在"
	NotifyUrlNotExitsErrCode          = 10010
	NotifyUrlNotExitsErrMessage       = "异步通知地址不存在"
	CallbackUrlNotExitsErrCode        = 10011
	CallbackUrlNotExitsErrMessage     = "同步通知地址不存在"
	NotSupportPaymentMethodErrCode    = 10012
	NotSupportPaymentMethodErrMessage = "不支持该支付方式"
)
