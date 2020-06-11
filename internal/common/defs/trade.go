package defs

type SearchRsp struct {
	Status  string  `json:"status"`   //交易状态
	OrderId string  `json:"order_id"` //订单号
	TradeNo string  `json:"trade_no"` //支付机构交易流水号
	Rate    float64 `json:"rate"`     //汇率
	RmbFee  float64 `json:"rmb_fee"`  //人民币金额
}

type CloseRsp struct {
	Status  bool   `json:"status"`   //交易关闭状态
	OrderId string `json:"order_id"` //订单号
}

/*type SearchReq struct {
	OrderId    string  `form:"order_id" json:"order_id"`       //订单号
	MethodCode string  `form:"method_code" json:"method_code"` //支付方式
	OrgCode    string  `form:"org_code" json:"org_code"`       //支付机构
	Currency   string  `form:"currency" json:"currency"`       //币种
	TotalFee   float64 `form:"total_fee" json:"total_fee"`     //金额
}*/

type CloseReq struct {
	OrderId    string  `form:"order_id" json:"order_id"`       //订单号
	TotalFee   float64 `form:"total_fee" json:"total_fee"`     //订单金额
	Currency   string  `form:"currency" json:"currency"`       //订单币种
	MethodCode string  `form:"method_code" json:"method_code"` //支付方式
	OrgCode    string  `form:"org_code" json:"org_code"`       //支付机构
}
