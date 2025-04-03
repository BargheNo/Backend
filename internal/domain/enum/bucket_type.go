package enum

type BucketType uint

const (
	VATTaxpayerCertificate BucketType = iota + 1
	OfficialNewspaperAD
)

func (bt BucketType) String() string {
	switch bt {
	case VATTaxpayerCertificate:
		return "vatTaxpayerCertificate"
	case OfficialNewspaperAD:
		return "officialNewspaperAD"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		VATTaxpayerCertificate,
		OfficialNewspaperAD,
	}
}
