package chatdto

type CreateOrGetRoomRequest struct {
	CorporationID uint
	UserID        uint
}

type GetRoomMessageRequest struct {
	RoomID uint
	UserID uint
}
