package payment

const (
	KRW                  = "KRW"    //韩币种符号
	JPY                  = "JPY"    //日币符号
	SignTypeMD5          = "MD5"    //MD5加密方式
	ChannelWechat        = "WECHAT" //微信支付
	ChannelAlipay        = "ALIPAY" //支付宝支付
	DateTimeFormatLayout = "2006-01-02 15:04:05"
)

const (
	TradeRefund     = "TRADE_REFUND"   //转入退款
	TradeSuccess    = "TRADE_SUCCESS"  //支付成功
	TradeWaitBuyPay = "WAIT_BUYER_PAY" //交易创建,等待用户支付
	TradeNotPay     = "TRADE_NOT_PAY"  //用户未支付
	TradeClosed     = "TRADE_CLOSED"   //交易关闭
	TradeError      = "TRADE_ERROR"    //交易失败
	TradeRevoked    = "TRADE_REVOKED"  //交易撤销
	TradeFinished   = "TRADE_FINISHED" //交易完成
)
