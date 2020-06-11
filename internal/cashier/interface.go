package cashier

import "payment_demo/internal/common/defs"

type PayMethod interface {
	AmpSubmit(defs.Order) (string, int, error)                         //发起支付
	OrderSubmit(defs.Order) (string, int, error)                       //发起支付
	OrderQrCode(defs.Order) (string, int, error)                       //发起支付
	Notify(string, string) (defs.NotifyRsp, int, error)                //异步通知
	Callback(string, string) (defs.CallbackRsp, int, error)            //同步通知
	Trade(string, string, string, float64) (defs.TradeRsp, int, error) //交易查询
	Closed(defs.Closed) (defs.ClosedRsp, int, error)                   //关闭交易
}
