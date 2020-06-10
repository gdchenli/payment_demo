package payment

const (
	CharsetUTF8          = "UTF-8"      //UTF-8编码
	Web                  = "web"        //PC端
	MobileWeb            = "mobile_web" //移动端
	Amp                  = "amp"        //小程序
	KRW                  = "KRW"        //韩币种符号
	JPY                  = "JPY"        //日币符号
	CNY                  = "CNY"        //人民币符号
	SignTypeMD5          = "MD5"        //MD5加密方式
	DateTimeFormatLayout = "2006-01-02 15:04:05"
)

const (
	TradeSuccess    = "TRADE_SUCCESS" //支付成功
	TradeFinished   = "TRADE_FINISHED"
	TradeWaitBuyPay = "WAIT_BUYER_PAY"
	TradeClosed     = "TRADE_CLOSED"
)
