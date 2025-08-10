package enum

type PaymentMethod uint

const (
	PaymentMethodCash PaymentMethod = iota + 1
	PaymentMethodInstallment
)

func (p PaymentMethod) String() string {
	switch p {
	case PaymentMethodCash:
		return "تفدی"
	case PaymentMethodInstallment:
		return "اقساط"
	}
	return "unknown"
}

func GetAllPaymentMethods() []PaymentMethod {
	return []PaymentMethod{
		PaymentMethodCash,
		PaymentMethodInstallment,
	}
}
