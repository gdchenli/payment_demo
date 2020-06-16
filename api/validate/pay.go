package validate

import (
	"errors"
	"payment_demo/internal/request"
)

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

type Order struct{}

func (o *Order) Validate(order *request.Order) (errCode int, err error) {
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
