package enum

type BucketType uint

const (
	VATTaxpayerCertificate BucketType = iota + 1
	OfficialNewspaperAD
	ProfilePic
)

func (bt BucketType) String() string {
	switch bt {
	case VATTaxpayerCertificate:
		return "vatTaxpayerCertificate"
	case OfficialNewspaperAD:
		return "officialNewspaperAD"
	case ProfilePic:
		return "profilePic"
	}
	return ""
}

func GetAllBucketTypes() []BucketType {
	return []BucketType{
		VATTaxpayerCertificate,
		OfficialNewspaperAD,
		ProfilePic,
	}
}
