package common

type UploadLogisticsRsp struct {
	Status  bool   `json:"status"`   //上传状态
	OrderId string `json:"order_id"` //订单号
}

type SearchTradeRsp struct {
	Status  string  `json:"status"`   //交易状态
	OrderId string  `json:"order_id"` //订单号
	TradeNo string  `json:"trade_no"` //支付机构交易流水号
	Rate    float64 `json:"rate"`     //汇率
	RmbFee  float64 `json:"rmb_fee"`  //人民币金额
	PaidAt  string  `json:"paid_at"`
}

type CloseTradeRsp struct {
	Status  bool   `json:"status"`   //交易关闭状态
	OrderId string `json:"order_id"` //订单号
}

type NotifyRsp struct {
	OrderId string  `json:"order_id"` //订单号
	Status  bool    `json:"status"`   //交易状态，true交易成功 false交易失败
	TradeNo string  `json:"trade_no"` //支付机构交易流水号
	Message string  `json:"message"`  //支付成功响应字符串
	RmbFee  float64 `json:"rmb_fee"`  //人民币金额
	PaidAt  string  `json:"paid_at"`
	Rate    float64 `json:"rate"`
}

type VerifyRsp struct {
	OrderId string `json:"order_id"` //订单号
	Status  bool   `json:"status"`   //交易状态，true交易成功 false交易失败
}

type AppPayRsp struct {
}

type WmpPayRsp struct {
}

type AmpPayRsp struct {
}
