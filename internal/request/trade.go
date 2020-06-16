package request

type SearchTradeReq struct {
	OrderId    string  `form:"order_id" json:"order_id"`       //订单号
	MethodCode string  `form:"method_code" json:"method_code"` //支付方式
	OrgCode    string  `form:"org_code" json:"org_code"`       //支付机构
	Currency   string  `form:"currency" json:"currency"`       //币种
	TotalFee   float64 `form:"total_fee" json:"total_fee"`     //金额
}

type CloseTradeReq struct {
	OrderId    string  `form:"order_id" json:"order_id"`       //订单号
	TotalFee   float64 `form:"total_fee" json:"total_fee"`     //订单金额
	Currency   string  `form:"currency" json:"currency"`       //订单币种
	MethodCode string  `form:"method_code" json:"method_code"` //支付方式
	OrgCode    string  `form:"org_code" json:"org_code"`       //支付机构
}
