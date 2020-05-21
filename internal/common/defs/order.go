package defs

type Order struct {
	OrderId         string  `json:"order_id"`          //订单编号
	TotalFee        float64 `json:"total_fee"`         //金额
	Currency        string  `json:"currency"`          //币种
	PaymentMethodId int     `json:"payment_method_id"` //支付方式
}
