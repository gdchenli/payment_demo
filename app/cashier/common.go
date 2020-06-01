package cashier

import "payment_demo/internal/method"

func getPayMethod(orgCode string) method.PayMethod {
	switch orgCode {
	case JdOrg:
		return new(method.Jd)
	case AllpayOrg:
		return new(method.Allpay)
	case AlipayOrg:
		return new(method.Alipay)
	case EpaymentsOrg:
		return new(method.Epayments)
	default:
		return nil
	}
}
