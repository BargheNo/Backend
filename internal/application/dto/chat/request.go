package chatdto

type CreateOrGetUserRoomRequest struct {
	CorporationID uint
	UserID        uint
}

type GetCorporationRoomRequest struct {
	CorporationID uint
	ApplicantID   uint
	UserPhone     string
}

type GetCorporationRoomsRequest struct {
	CorporationID uint
	ApplicantID   uint
}

type GetRoomMessageRequest struct {
	RoomID uint
	UserID uint
}
