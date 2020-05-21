package dbops

import "time"

type Order struct {
	Id              int       `json:"id"`
	OrderId         string    `json:"order_id"`          //订单编号
	TotalFee        float64   `json:"total_fee"`         //金额
	Currency        string    `json:"currency"`          //币种
	PaymentMethodId int       `json:"payment_method_id"` //支付方式
	Status          string    `json:"status"`            //订单状态
	PaidAt          time.Time `json:"paid_at"`           //付款时间
	CreatedAt       time.Time `json:"created_at"`        //创建时间
	UpdatedAt       time.Time `json:"updated_at"`        //更新时间
}
