package enum

type ChatStatus uint

const (
	ChatStatusActive ChatStatus = iota + 1
	ChatStatusBlocked
)

func (status ChatStatus) String() string {
	switch status {
	case ChatStatusActive:
		return "active"
	case ChatStatusBlocked:
		return "blocked"
	}
	return "active"
}

func GetAllChatStatuses() []ChatStatus {
	return []ChatStatus{
		ChatStatusActive,
		ChatStatusBlocked,
	}
}
