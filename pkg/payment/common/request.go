package common

type OrderArg struct {
	OrderId       string  `form:"order_id" json:"order_id"`               //订单编号
	TotalFee      float64 `form:"total_fee" json:"total_fee"`             //金额
	Currency      string  `form:"currency" json:"currency"`               //币种
	MethodCode    string  `form:"method_code" json:"method_code"`         //支付方式
	OrgCode       string  `form:"org_code" json:"org_code"`               //支付机构
	UserId        string  `form:"user_id" json:"user_id"`                 //用户Id
	UserAgentType int     `form:"user_agent_type" json:"user_agent_type"` //环境
}

type UploadLogisticsArg struct {
	OrderId          string `form:"order_id" json:"order_id"`                   //订单号
	LogisticsNo      string `form:"logistics_no" json:"logistics_no"`           //物流单号
	LogisticsCompany string `form:"logistics_company" json:"logistics_company"` //物流公司名称
	OrgCode          string `form:"org_code" json:"org_code"`                   //支付机构
}

type SearchTradeArg struct {
	OrderId    string  `form:"order_id" json:"order_id"`       //订单号
	MethodCode string  `form:"method_code" json:"method_code"` //支付方式
	OrgCode    string  `form:"org_code" json:"org_code"`       //支付机构
	Currency   string  `form:"currency" json:"currency"`       //币种
	TotalFee   float64 `form:"total_fee" json:"total_fee"`     //金额
}

type CloseTradeArg struct {
	OrderId    string  `form:"order_id" json:"order_id"`       //订单号
	TotalFee   float64 `form:"total_fee" json:"total_fee"`     //订单金额
	Currency   string  `form:"currency" json:"currency"`       //订单币种
	MethodCode string  `form:"method_code" json:"method_code"` //支付方式
	OrgCode    string  `form:"org_code" json:"org_code"`       //支付机构
}
