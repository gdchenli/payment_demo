package consts

//支付方式code
const (
	JdMethod       = "jd_payment"
	AlipayMethod   = "alipay_payment"    //支付宝
	UnionpayMethod = "vtpayment_payment" //银联
	WechatMethod   = "wechat"            //微信支付
	JdOrg          = "jd"
	AllpayOrg      = "allpay"
	AlipayOrg      = "alipay"
	EpaymentsOrg   = "epayments"
)

//支付环境
const (
	WebUserAgentType               = 1 //pc端
	MobileUserAgentType            = 2 //移动端
	AndroidAppUserAgentType        = 3 //安卓app
	IOSAppUserAgentType            = 4 //苹果app
	WmpUserAgentType               = 5 //微信浏览器
	WechatMiniProgramUserAgentType = 6 //微信小程序
	AlipayMiniProgramUserAgentType = 7 //支付宝小程序
)
