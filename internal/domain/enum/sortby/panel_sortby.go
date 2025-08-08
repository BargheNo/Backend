package sortby

type PanelSortBy uint

const (
	PanelSortByCreatedAt PanelSortBy = iota + 1
	PanelSortByName
	PanelSortByPower
)

func (s PanelSortBy) Name() string {
	switch s {
	case PanelSortByCreatedAt:
		return "تاریخ ایجاد"
	case PanelSortByName:
		return "نام"
	case PanelSortByPower:
		return "توان"
	}
	return "ناشناس"
}

func (s PanelSortBy) DBColumn() string {
	switch s {
	case PanelSortByCreatedAt:
		return "created_at"
	case PanelSortByName:
		return "name"
	case PanelSortByPower:
		return "power"
	}
	return ""
}

func GetPanelSortableColumns() map[PanelSortBy]bool {
	return map[PanelSortBy]bool{
		PanelSortByCreatedAt: true,
		PanelSortByName:      true,
		PanelSortByPower:     true,
	}
}
