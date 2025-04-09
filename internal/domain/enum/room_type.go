package enum

type RoomType uint

const (
	RoomTypeCorporation RoomType = iota + 1
)

func (roomType RoomType) String() string {
	switch roomType {
	case RoomTypeCorporation:
		return "corporation"
	}
	return ""
}

func GetAllRoomTypes() []RoomType {
	return []RoomType{
		RoomTypeCorporation,
	}
}
