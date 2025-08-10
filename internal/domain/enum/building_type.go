package enum

type BuildingType uint

const (
	BuildingTypeResidential BuildingType = iota + 1
	BuildingTypeCommercial
	BuildingTypeIndustrial
	BuildingTypeAgricultural
	BuildingTypeEducational
	BuildingTypeGovernmental
)

func (b BuildingType) String() string {
	switch b {
	case BuildingTypeResidential:
		return "مسکونی"
	case BuildingTypeCommercial:
		return "تجاری"
	case BuildingTypeIndustrial:
		return "صنعتی"
	case BuildingTypeAgricultural:
		return "کشاورزی"
	case BuildingTypeEducational:
		return "آموزشی"
	case BuildingTypeGovernmental:
		return "دولتی"
	}
	return "unknown"
}

func GetAllBuildingTypes() []BuildingType {
	return []BuildingType{
		BuildingTypeResidential,
		BuildingTypeCommercial,
		BuildingTypeIndustrial,
		BuildingTypeAgricultural,
		BuildingTypeEducational,
		BuildingTypeGovernmental,
	}
}
