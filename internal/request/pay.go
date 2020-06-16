package request

import "errors"

const (
	RequiredPayOrderIdErrCode     = 10150
	RequiredPayOrderIdErrMessage  = "请输入发起支付的订单编号"
	RequiredPayTotalFeeErrCode    = 10151
	RequiredPayTotalFeeErrMessage = "请输入发起支付的订单金额"
	RequiredPayCurrencyErrCode    = 10152
	RequiredPayCurrencyErrMessage = "请输入发起支付的订单币种"
	RequiredPayMethodErrCode      = 10153
	RequiredPayMethodMessage      = "请选择发起支付的支付方式"
	RequiredPayOrgErrCode         = 10154
	RequiredPayOrgErrMessage      = "请选择发起支付的支付机构"
	RequiredUserIdErrCode         = 10155
	RequiredUserIdErrMessage      = "请输入发起支付的用户Id"
)

type Order struct {
	OrderId       string  `form:"order_id" json:"order_id"`               //订单编号
	TotalFee      float64 `form:"total_fee" json:"total_fee"`             //金额
	Currency      string  `form:"currency" json:"currency"`               //币种
	MethodCode    string  `form:"method_code" json:"method_code"`         //支付方式
	OrgCode       string  `form:"org_code" json:"org_code"`               //支付机构
	UserId        string  `form:"user_id" json:"user_id"`                 //用户Id
	UserAgentType int     `form:"user_agent_type" json:"user_agent_type"` //环境
}

func (order *Order) Validate() (errCode int, err error) {
	if order.OrderId == "" {
		return RequiredPayOrderIdErrCode, errors.New(RequiredPayOrderIdErrMessage)
	}

	if order.TotalFee == 0 {
		return RequiredPayTotalFeeErrCode, errors.New(RequiredPayTotalFeeErrMessage)
	}

	if order.Currency == "" {
		return RequiredPayCurrencyErrCode, errors.New(RequiredPayCurrencyErrMessage)
	}

	if order.MethodCode == "" {
		return RequiredPayMethodErrCode, errors.New(RequiredPayMethodMessage)
	}

	if order.OrgCode == "" {
		return RequiredPayOrgErrCode, errors.New(RequiredPayOrgErrMessage)
	}

	if order.OrgCode == "jd" && order.UserId == "" {
		return RequiredUserIdErrCode, errors.New(RequiredUserIdErrMessage)
	}

	return 0, nil
}
