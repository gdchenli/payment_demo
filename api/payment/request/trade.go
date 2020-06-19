package request

import "errors"

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
