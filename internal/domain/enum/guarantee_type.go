package enum

type GuaranteeType uint

const (
	GuaranteeTypeWarranty GuaranteeType = iota + 1
	GuaranteeTypeInsurance
	GuaranteeTypeServiceContract
)

func (g GuaranteeType) String() string {
	switch g {
	case GuaranteeTypeWarranty:
		return "گارانتی تعویض"
	case GuaranteeTypeInsurance:
		return "بیمه"
	case GuaranteeTypeServiceContract:
		return "قرارداد خدمات پس از فروش"
	}
	return "unknown"
}

func GetAllGuaranteeTypes() []GuaranteeType {
	return []GuaranteeType{
		GuaranteeTypeWarranty,
		GuaranteeTypeInsurance,
		GuaranteeTypeServiceContract,
	}
}
