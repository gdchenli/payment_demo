package defs

import "errors"

type PayMethod interface {
	OrderSubmit(Order) (string, int, error)            //发起支付
	Notify(string, string) (NotifyRsp, int, error)     //异步通知
	Callback(string, string) (CallbackRsp, int, error) //同步通知
	Trade(string, string) (TradeRsp, int, error)       //交易查询
}

const (
	RequiredPayOrderIdErrCode            = 10150
	RequiredPayOrderIdErrMessage         = "请输入发起支付的订单编号"
	RequiredPayTotalFeeErrCode           = 10151
	RequiredPayTotalFeeErrMessage        = "请输入发起支付的订单金额"
	RequiredPayCurrencyErrCode           = 10152
	RequiredPayCurrencyErrMessage        = "请输入发起支付的订单币种"
	RequiredPayMethodErrCode             = 10153
	RequiredPayMethodMessage             = "请选择发起支付的支付方式"
	RequiredPayOrgErrCode                = 10154
	RequiredPayOrgErrMessage             = "请选择发起支付的支付机构"
	RequiredUserIdErrCode                = 10155
	RequiredUserIdErrMessage             = "请输入发起支付的用户Id"
	RequiredSearchTradeOrderIdErrCode    = 10550
	RequiredSearchTradeOrderIdErrMessage = "请输入需要查询交易的订单编号"
	RequiredSearchTradeMethodErrCode     = 10551
	RequiredSearchTradeMethodMessage     = "请选择需要查询交易的支付方式"
	RequiredSearchTradeOrgErrCode        = 10552
	RequiredSearchTradeOrgErrMessage     = "请选择需要查询交易的支付机构"
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

type NotifyRsp struct {
	OrderId string `json:"order_id"` //订单号
	Status  bool   `json:"status"`   //交易状态，true交易成功 false交易失败
	TradeNo string `json:"trade_no"` //支付机构交易流水号
	Message string `json:"message"`  //支付成功响应字符串
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
		return RequiredSearchTradeOrderIdErrCode, errors.New(RequiredSearchTradeOrderIdErrMessage)
	}

	if trade.MethodCode == "" {
		return RequiredSearchTradeMethodErrCode, errors.New(RequiredSearchTradeMethodMessage)
	}

	if trade.OrgCode == "" {
		return RequiredSearchTradeOrgErrCode, errors.New(RequiredSearchTradeOrgErrMessage)
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
		return RequiredCloseTradeOrderIdErrCode, errors.New(RequiredCloseTradeOrderIdErrMessage)
	}

	if closed.OrgCode == "jd" && closed.TotalFee == 0 {
		return RequiredCloseTradeTotalFeeErrCode, errors.New(RequiredCloseTradeTotalFeeErrMessage)
	}

	if closed.OrgCode == "jd" && closed.Currency == "" {
		return RequiredCloseTradeCurrencyErrCode, errors.New(RequiredCloseTradeCurrencyErrMessage)
	}

	if closed.MethodCode == "" {
		return RequiredCloseTradeMethodErrCode, errors.New(RequiredCloseTradeMethodErrMessage)
	}

	if closed.OrgCode == "" {
		return RequiredCloseTradeOrgErrCode, errors.New(RequiredCloseTradeOrgErrMessage)
	}

	return 0, nil
}

type Logistics struct {
	OrderId          string `form:"order_id" json:"order_id"`                   //订单号
	LogisticsNo      string `form:"logistics_no" json:"logistics_no"`           //物流单号
	LogisticsCompany string `form:"logistics_company" json:"logistics_company"` //物流公司名称
	OrgCode          string `form:"org_code" json:"org_code"`                   //支付机构
}

type LogisticsRsp struct {
	Status  bool   `json:"status"`   //上传状态
	OrderId string `json:"order_id"` //订单号
}
