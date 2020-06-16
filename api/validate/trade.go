package validate

import (
	"errors"
	"payment_demo/internal/request"
)

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

type SearchTradeReq struct{}

func (t *SearchTradeReq) Validate(searchTradeReq *request.SearchTradeReq) (errCode int, err error) {
	if searchTradeReq.OrderId == "" {
		return RequiredSearchTradeOrderIdErrCode, errors.New(RequiredSearchTradeOrderIdErrMessage)
	}

	if searchTradeReq.MethodCode == "" {
		return RequiredSearchTradeMethodErrCode, errors.New(RequiredSearchTradeMethodMessage)
	}

	if searchTradeReq.OrgCode == "" {
		return RequiredSearchTradeOrgErrCode, errors.New(RequiredSearchTradeOrgErrMessage)
	}

	if searchTradeReq.TotalFee == 0 {
		return RequiredTradeTotalFeeErrCode, errors.New(RequiredTradeTotalFeeErrMessage)
	}

	if searchTradeReq.Currency == "" {
		return RequiredTradeCurrencyErrCode, errors.New(RequiredTradeCurrencyErrMessage)
	}

	return 0, nil
}

type CloseTradeReq struct{}

func (t *CloseTradeReq) Validate(closeTradeReq *request.CloseTradeReq) (errCode int, err error) {
	if closeTradeReq.OrderId == "" {
		return RequiredCloseTradeOrderIdErrCode, errors.New(RequiredCloseTradeOrderIdErrMessage)
	}

	if closeTradeReq.OrgCode == "jd" && closeTradeReq.TotalFee == 0 {
		return RequiredCloseTradeTotalFeeErrCode, errors.New(RequiredCloseTradeTotalFeeErrMessage)
	}

	if closeTradeReq.OrgCode == "jd" && closeTradeReq.Currency == "" {
		return RequiredCloseTradeCurrencyErrCode, errors.New(RequiredCloseTradeCurrencyErrMessage)
	}

	if closeTradeReq.MethodCode == "" {
		return RequiredCloseTradeMethodErrCode, errors.New(RequiredCloseTradeMethodErrMessage)
	}

	if closeTradeReq.OrgCode == "" {
		return RequiredCloseTradeOrgErrCode, errors.New(RequiredCloseTradeOrgErrMessage)
	}

	return 0, nil
}
