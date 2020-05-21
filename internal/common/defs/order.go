package defs

import "errors"

type Order struct {
	OrderId    string  `form:"order_id" json:"order_id"`       //订单编号
	TotalFee   float64 `form:"total_fee" json:"total_fee"`     //金额
	Currency   string  `form:"currency" json:"currency"`       //币种
	MethodCode string  `form:"method_code" json:"method_code"` //支付方式
	OrgCode    string  `form:"org_code" json:"org_code"`       //支付机构
	UserId     string  `form:"user_id" json:"user_id"`         //用户Id
}

func (order *Order) Validate() (errCode int, err error) {
	if order.OrderId == "" {
		return 10101, errors.New("请输入订单编号")
	}

	if order.TotalFee == 0 {
		return 10102, errors.New("请输入订单金额")
	}
	if order.Currency == "" {
		return 10103, errors.New("请输入订单币种")
	}

	if order.MethodCode == "" {
		return 10104, errors.New("请选择支付方式")
	}

	if order.OrgCode == "" {
		return 10105, errors.New("请选择支付机构")
	}

	if order.OrgCode == "Jd" && order.UserId == "" {
		return 10106, errors.New("请输入用户Id")
	}

	return 0, nil
}

type NotifyRsp struct {
	OrderId string `json:"order_id"` //订单号
	Status  bool   `json:"status"`   //交易状态，true交易成功 false交易失败
}
