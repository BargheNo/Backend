package monitoringdto

type CustomerPanelStatusListRequest struct {
	PanelID uint
	OwnerID uint
	Offset  int
	Limit   int
}

type CorporationPanelStatusListRequest struct {
	CorporationID uint
	UserID        uint
	PanelID       uint
	Offset        int
	Limit         int
}

type AdminPanelStatusListRequest struct {
	PanelID uint
	Offset  int
	Limit   int
}
