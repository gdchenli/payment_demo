package response

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
