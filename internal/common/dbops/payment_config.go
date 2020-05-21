package dbops

import "time"

type PaymentConfig struct {
	Id        int       `json:"id"`
	MethodId  int       `json:"method_id"`
	Code      string    `json:"code"`       //配置标识
	Name      string    `json:"name"`       //配置名称
	Value     string    `json:"value"`      //配置值
	CreatedAt time.Time `json:"created_at"` //创建时间
	UpdatedAt time.Time `json:"updated_at"` //更新时间
}
