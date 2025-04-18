package chatdto

import "github.com/BargheNo/Backend/internal/domain/enum"

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
	Offset int
	Limit  int
}

type BlockServiceChatRequest struct {
	UserID     uint
	RoomID     uint
	BlockedBy  enum.BlockedBy
	ChatStatus enum.ChatStatus
}
