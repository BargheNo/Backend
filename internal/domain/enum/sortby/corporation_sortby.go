package sortby

type CorporationSortBy uint

const (
	CorporationSortByCreatedAt CorporationSortBy = iota + 1
	CorporationSortByName
)

func (s CorporationSortBy) Name() string {
	switch s {
	case CorporationSortByCreatedAt:
		return "تاریخ ایجاد"
	case CorporationSortByName:
		return "نام"
	}
	return "ناشناس"
}

func (s CorporationSortBy) DBColumn() string {
	switch s {
	case CorporationSortByCreatedAt:
		return "created_at"
	case CorporationSortByName:
		return "name"
	}
	return ""
}

func GetCorporationSortableColumns() map[CorporationSortBy]bool {
	return map[CorporationSortBy]bool{
		CorporationSortByCreatedAt: true,
		CorporationSortByName:      true,
	}
}
