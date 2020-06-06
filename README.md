# payment demo

### 支付方式
|支付机构|支付方式|版本|
|---|---|---|
|京东|京东支付|2.0|
|allpay|支付宝支付、银联支付|5.0|
|epayments|微信支付||
|支付宝|支付宝支付||

### 运行项目

拷贝配置文件
```
cp cmd/cashier/config.toml.example cmd/cashier/config.toml
```
配置支付参数 

运行项目
```
cd cmd/cashier/ && go build . && ./cashier
```

### 目录结构
```
.
├── README.md
├── app
│   └── cashier                 http响应处理文件夹
│       ├── callback.go         支付同步回调
│       ├── common.go           
│       ├── logistics.go        jd支付上传物流信息
│       ├── notify.go           支付异步回调
│       ├── pay.go              发起支付
│       └── trade.go            交易查询
├── cmd
│   └── cashier
│       ├── config.toml.exmple  配置示例文件
│       ├── logs                日志目录
│       └── main.go             入口
├── go.mod
├── go.sum
├── internal
│   ├── cashier
│   │   ├── alipay.go           支付宝直连
│   │   ├── allpay.go           allpay
│   │   ├── epayments.go        epayments
│   │   ├── interface.go
│   │   └── jd.go               京东支付
│   └── common
├── pkg                         工具文件夹
└── vendor                      第三方依赖包文件夹
```