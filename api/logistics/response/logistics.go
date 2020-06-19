package response

type UploadLogisticsRsp struct {
	Status  bool   `json:"status"`   //上传状态
	OrderId string `json:"order_id"` //订单号
}
