package payment

const (
	Version              = "V2.0" //京东支付版本
	DateTimeFormatLayout = "2006-01-02 15:04:05"
)

const (
	TradeCreate  = "0" //交易创建
	TradePending = "1" //交易处理中
	TradeProcess = "2" //交易成功
	TradeClosed  = "3" //交易关闭
	TradeFailed  = "4" //交易失败
)
