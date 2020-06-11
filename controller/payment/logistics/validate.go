package logistics

type Req struct {
	OrderId          string `form:"order_id" json:"order_id"`                   //订单号
	LogisticsNo      string `form:"logistics_no" json:"logistics_no"`           //物流单号
	LogisticsCompany string `form:"logistics_company" json:"logistics_company"` //物流公司名称
	OrgCode          string `form:"org_code" json:"org_code"`                   //支付机构
}
