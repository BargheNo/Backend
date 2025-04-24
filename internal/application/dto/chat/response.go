package chatdto

import (
	"time"

	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type ChatRoomDetailsResponse struct {
	RoomID                uint                                         `json:"roomID"`
	CustomerCredential    userdto.CredentialResponse                   `json:"customer"`
	CorporationCredential corporationdto.CorporationCredentialResponse `json:"corporation"`
	Status                string                                       `json:"status"`
	BlockedBy             string                                       `json:"blockedBy"`
}

type RoomMessagesResponse struct {
	ID        uint                       `json:"id"`
	Sender    userdto.CredentialResponse `json:"sender"`
	Content   string                     `json:"content"`
	TimeStamp time.Time                  `json:"timeStamp"`
}
