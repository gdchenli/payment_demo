# payment demo

### 支付方式
|支付机构|支付方式|版本|
|---|---|---|
|京东|京东支付|2.0|
|allpay|支付宝支付、银联支付|5.0|
|epayments|微信支付||
|支付宝|支付宝支付||
|微信|微信支付||


### 运行项目

拷贝配置文件
```
cp cmd/payment/config.toml.example cmd/payment/config.toml
```
配置支付参数 

运行项目
```
cd cmd/payment/ && go build . && ./payment
```

### 目录结构