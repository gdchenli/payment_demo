package defs

type PaymentConfig struct {
	MethodId int    `json:"method_id"` //支付方式
	Code     string `json:"code"`      //配置标识
	Name     string `json:"name"`      //配置名称
	Value    string `json:"value"`     //配置值
}
