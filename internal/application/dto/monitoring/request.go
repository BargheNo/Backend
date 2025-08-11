package monitoringdto

type CustomerPanelStatusListRequest struct {
	PanelID uint
	OwnerID uint
	Offset  int
	Limit   int
}
