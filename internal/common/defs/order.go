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
		return 10150, errors.New("请输入订单编号")
	}

	if order.TotalFee == 0 {
		return 10151, errors.New("请输入订单金额")
	}
	if order.Currency == "" {
		return 10152, errors.New("请输入订单币种")
	}

	if order.MethodCode == "" {
		return 10153, errors.New("请选择支付方式")
	}

	if order.OrgCode == "" {
		return 10154, errors.New("请选择支付机构")
	}

	if order.OrgCode == "Jd" && order.UserId == "" {
		return 10155, errors.New("请输入用户Id")
	}

	return 0, nil
}

type NotifyRsp struct {
	OrderId string `json:"order_id"` //订单号
	Status  bool   `json:"status"`   //交易状态，true交易成功 false交易失败
	TradeNo string `json:"trade_no"` //支付机构交易流水号
}

type CallbackRsp struct {
	OrderId string `json:"order_id"` //订单号
	Status  bool   `json:"status"`   //交易状态，true交易成功 false交易失败
}

type TradeRsp struct {
	Status  string `json:"status"`   //交易状态
	OrderId string `json:"order_id"` //订单号
	TradeNo string `json:"trade_no"` //支付机构交易流水号
}

type Trade struct {
	OrderId    string `form:"order_id" json:"order_id"`       //订单号
	MethodCode string `form:"method_code" json:"method_code"` //支付方式
	OrgCode    string `form:"org_code" json:"org_code"`       //支付机构
}

func (trade *Trade) Validate() (errCode int, err error) {
	if trade.OrderId == "" {
		return 10550, errors.New("请输入订单编号")
	}

	if trade.MethodCode == "" {
		return 10551, errors.New("请选择支付方式")
	}

	if trade.OrgCode == "" {
		return 10552, errors.New("请选择支付机构")
	}

	return 0, nil
}

type ClosedRsp struct {
	Status  bool   `json:"status"`   //交易关闭状态
	OrderId string `json:"order_id"` //订单号
}

type Closed struct {
	OrderId    string  `form:"order_id" json:"order_id"`       //订单号
	TotalFee   float64 `form:"total_fee" json:"total_fee"`     //订单金额
	Currency   string  `form:"currency" json:"currency"`       //订单币种
	MethodCode string  `form:"method_code" json:"method_code"` //支付方式
	OrgCode    string  `form:"org_code" json:"org_code"`       //支付机构
}

func (closed *Closed) Validate() (errCode int, err error) {
	if closed.OrderId == "" {
		return 10450, errors.New("请输入订单编号")
	}

	if closed.TotalFee == 0 {
		return 10451, errors.New("请输入订单金额")
	}

	if closed.Currency == "" {
		return 10452, errors.New("请输入币种")
	}

	if closed.MethodCode == "" {
		return 10453, errors.New("请选择支付方式")
	}

	if closed.OrgCode == "" {
		return 10454, errors.New("请选择支付机构")
	}

	return 0, nil
}
