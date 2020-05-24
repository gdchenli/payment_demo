package payment

type Notify struct{}

type NotifyRsp struct {
	OrderNum  string `json:"orderNum"`
	TransId   string `json:"transID"`
	RespCode  string `json:"RespCode"`
	RespMsg   string `json:"RespMsg"`
	Signature string `json:"signature"`
}
