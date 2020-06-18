# payment demo

### 支付方式
|支付机构|支付方式|版本|支付机构文档|
|---|---|---|---|
|京东|京东支付|2.0|[文档](http://payapi.jd.com/docList.html?methodName=0#)|
|allpay|支付宝支付、银联支付|5.0|[文档](https://git.allpayx.com/OpenAPI/common/src/master/AllPay_Integration_Specification_CH.md)|
|epayments|微信支付||[文档](https://www.kiwifast.com/doc/)|
|支付宝|支付宝支付||[文档](https://global.alipay.com/docs/ac/global/create_forex_trade)|


### 运行项目

拷贝配置文件
```
cp cmd/config.toml.example cmd/config.toml
```
配置支付参数 

运行项目
```
cd cmd/ && go build . && ./cmd
```

### 目录结构
```
.
├── README.md
├── api
│   ├── controller          #http处理逻辑
│   │   └── payment.go
│   ├── response            #响应数据结构
│   │   └── payment.go
│   └── validate            #请求数据校验
│       ├── logistics.go
│       ├── pay.go
│       └── trade.go
├── cmd
│   ├── config.toml.exmple  #示例配置文件
│   ├── logs                #日志记录
│   └── main.go             #入口文件
├── go.mod
├── go.sum
├── internal
│   ├── common
│   │   ├── code            #公用错误码
│   │   └── request         #请求数据
│   │       ├── logistics.go
│   │       ├── pay.go
│   │       └── trade.go
│   └── service
│       └── payment         #支付流程逻辑
└── pkg
    ├── config              
    ├── curl
    ├── ginprometheus
    ├── grace
    ├── log
    ├── payment             #支付方式对接处理逻辑
    │   ├── alipay          #支付宝ISV对接处理逻辑
    │   ├── allpay          #支付机构allpay对接处理逻辑
    │   ├── consts          #公用常量
    │   ├── epayments       #支付机构epayments对接处理逻辑
    │   └── jd              #京东支付ISV对接处理逻辑
    └── recovery
```