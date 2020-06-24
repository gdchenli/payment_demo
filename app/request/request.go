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

type OrderArg struct {
	OrderId       string  `form:"order_id" json:"order_id"`               //订单编号
	TotalFee      float64 `form:"total_fee" json:"total_fee"`             //金额
	Currency      string  `form:"currency" json:"currency"`               //币种
	MethodCode    string  `form:"method_code" json:"method_code"`         //支付方式
	OrgCode       string  `form:"org_code" json:"org_code"`               //支付机构
	UserId        string  `form:"user_id" json:"user_id"`                 //用户Id
	UserAgentType int     `form:"user_agent_type" json:"user_agent_type"` //环境
}

func (arg *OrderArg) Validate() (errCode int, err error) {
	if arg.OrderId == "" {
		return RequiredPayOrderIdErrCode, errors.New(RequiredPayOrderIdErrMessage)
	}

	if arg.TotalFee == 0 {
		return RequiredPayTotalFeeErrCode, errors.New(RequiredPayTotalFeeErrMessage)
	}

	if arg.Currency == "" {
		return RequiredPayCurrencyErrCode, errors.New(RequiredPayCurrencyErrMessage)
	}

	if arg.MethodCode == "" {
		return RequiredPayMethodErrCode, errors.New(RequiredPayMethodMessage)
	}

	if arg.OrgCode == "" {
		return RequiredPayOrgErrCode, errors.New(RequiredPayOrgErrMessage)
	}

	if arg.OrgCode == "jd" && arg.UserId == "" {
		return RequiredUserIdErrCode, errors.New(RequiredUserIdErrMessage)
	}

	return 0, nil
}

type UploadLogisticsArg struct {
	OrderId          string `form:"order_id" json:"order_id"`                   //订单号
	LogisticsNo      string `form:"logistics_no" json:"logistics_no"`           //物流单号
	LogisticsCompany string `form:"logistics_company" json:"logistics_company"` //物流公司名称
	OrgCode          string `form:"org_code" json:"org_code"`                   //支付机构
}

const (
	RequiredSearchTradeOrderIdErrCode    = 10550
	RequiredSearchTradeOrderIdErrMessage = "请输入需要查询交易的订单编号"
	RequiredSearchTradeMethodErrCode     = 10551
	RequiredSearchTradeMethodMessage     = "请选择需要查询交易的支付方式"
	RequiredSearchTradeOrgErrCode        = 10552
	RequiredSearchTradeOrgErrMessage     = "请选择需要查询交易的支付机构"
	RequiredTradeTotalFeeErrCode         = 10553
	RequiredTradeTotalFeeErrMessage      = "请输入查询交易的订单金额"
	RequiredTradeCurrencyErrCode         = 10554
	RequiredTradeCurrencyErrMessage      = "请输入查询交易的订单币种"
)

const (
	RequiredCloseTradeOrderIdErrCode     = 10450
	RequiredCloseTradeOrderIdErrMessage  = "请输入需要关闭交易的订单编号"
	RequiredCloseTradeTotalFeeErrCode    = 10451
	RequiredCloseTradeTotalFeeErrMessage = "请输入需要关闭交易的订单金额"
	RequiredCloseTradeCurrencyErrCode    = 10453
	RequiredCloseTradeCurrencyErrMessage = "请输入需要关闭交易的订单币种"
	RequiredCloseTradeMethodErrCode      = 10453
	RequiredCloseTradeMethodErrMessage   = "请选择需要关闭交易的支付方式"
	RequiredCloseTradeOrgErrCode         = 10454
	RequiredCloseTradeOrgErrMessage      = "请选择需要关闭交易的支付机构"
)

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

func (arg *SearchTradeArg) Validate() (errCode int, err error) {
	if arg.OrderId == "" {
		return RequiredSearchTradeOrderIdErrCode, errors.New(RequiredSearchTradeOrderIdErrMessage)
	}

	if arg.MethodCode == "" {
		return RequiredSearchTradeMethodErrCode, errors.New(RequiredSearchTradeMethodMessage)
	}

	if arg.OrgCode == "" {
		return RequiredSearchTradeOrgErrCode, errors.New(RequiredSearchTradeOrgErrMessage)
	}

	if arg.TotalFee == 0 {
		return RequiredTradeTotalFeeErrCode, errors.New(RequiredTradeTotalFeeErrMessage)
	}

	if arg.Currency == "" {
		return RequiredTradeCurrencyErrCode, errors.New(RequiredTradeCurrencyErrMessage)
	}

	return 0, nil
}

func (arg *CloseTradeArg) Validate() (errCode int, err error) {
	if arg.OrderId == "" {
		return RequiredCloseTradeOrderIdErrCode, errors.New(RequiredCloseTradeOrderIdErrMessage)
	}

	if arg.OrgCode == "jd" && arg.TotalFee == 0 {
		return RequiredCloseTradeTotalFeeErrCode, errors.New(RequiredCloseTradeTotalFeeErrMessage)
	}

	if arg.OrgCode == "jd" && arg.Currency == "" {
		return RequiredCloseTradeCurrencyErrCode, errors.New(RequiredCloseTradeCurrencyErrMessage)
	}

	if arg.MethodCode == "" {
		return RequiredCloseTradeMethodErrCode, errors.New(RequiredCloseTradeMethodErrMessage)
	}

	if arg.OrgCode == "" {
		return RequiredCloseTradeOrgErrCode, errors.New(RequiredCloseTradeOrgErrMessage)
	}

	return 0, nil
}
