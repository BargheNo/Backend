package reportdto

type CreateReportRequest struct {
	ObjectID       uint
	ObjectType     string
	Description    string
	ReportedByID   uint
	ReportedByType string
}

type ReportListRequest struct {
	OwnerID uint
	Status  uint
	Offset  int
	Limit   int
	SortBy  uint
	Asc     bool
}

type ResolveReportRequest struct {
	ReportID uint
	UserID   uint
}

type ReportNotificationData struct {
	ReportID uint
}
