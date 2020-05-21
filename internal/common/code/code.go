package code

const (
	AmountFormatErrCode          = 10101
	AmountFormatErrMessage       = "金额转换异常"
	PrivateKeyNotExitsErrCode    = 10102
	PrivateKeyNotExitsErrMessage = "私钥文件不存在"
	PrivateKeyContentErrCode     = 10102
	PrivateKeyContentErrMessage  = "私钥内容错误"

	PaymentGoodsInfoFormatErrCode    = 10150
	PaymentGoodsInfoFormatErrMessage = "发起支付，商品数据转换失败"
	PaymentKjInfoFormatErrCode       = 10151
	PaymentKjInfoFormatErrMessage    = "发起支付，跨境数据转换失败"
	PaymentSignErrCode               = 10152
	PaymentSignErrMessage            = "发起支付，签名计算错误"
)
