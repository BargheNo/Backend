package reportdto

type CreateReportRequest struct {
	ObjectID       uint
	ObjectType     string
	Description    string
	ReportedByID   uint
	ReportedByType string
}
