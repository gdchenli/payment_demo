package dbops

import "time"

type PaymentMethod struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`       //支付方式名称
	CreatedAt time.Time `json:"created_at"` //创建时间
	UpdatedAt time.Time `json:"updated_at"` //更新时间
}
