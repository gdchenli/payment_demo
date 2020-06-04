package cashier

import "payment_demo/internal/cashier"

func getPayMethod(orgCode string) cashier.PayMethod {
	switch orgCode {
	case JdOrg:
		return new(cashier.Jd)
	case AllpayOrg:
		return new(cashier.Allpay)
	case AlipayOrg:
		return new(cashier.Alipay)
	case EpaymentsOrg:
		return new(cashier.Epayments)
	default:
		return nil
	}
}
