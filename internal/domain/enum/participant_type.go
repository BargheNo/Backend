package enum

type ParticipantType uint

const (
	ParticipantTypeUser ParticipantType = iota + 1
	ParticipantTypeAdmin
)

func (participantType ParticipantType) String() string {
	switch participantType {
	case ParticipantTypeUser:
		return "user"
	case ParticipantTypeAdmin:
		return "admin"
	}
	return ""
}

func GetAllParticipantTypes() []ParticipantType {
	return []ParticipantType{
		ParticipantTypeUser,
		ParticipantTypeAdmin,
	}
}
