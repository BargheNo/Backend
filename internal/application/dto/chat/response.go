package chatdto

import (
	corporationdto "github.com/BargheNo/Backend/internal/application/dto/corporation"
	userdto "github.com/BargheNo/Backend/internal/application/dto/user"
)

type ChatRoomDetailsResponse struct {
	RoomID                uint                                      `json:"roomID"`
	UserCredential        userdto.CredentialResponse                `json:"user"`
	CorporationCredential corporationdto.CorporationDetailsResponse `json:"corporation"`
}

type RoomMessagesResponse struct {
	Sender  userdto.CredentialResponse `json:"sender"`
	Content string                     `json:"content"`
}
