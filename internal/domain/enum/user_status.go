package enum

type UserStatus uint

const (
	UserStatusActive UserStatus = iota + 1
	UserStatusBlock
)

func (status UserStatus) String() string {
	switch status {
	case UserStatusActive:
		return "active"
	case UserStatusBlock:
		return "block"
	}
	return ""
}

func GetAllUserStatus() []UserStatus {
	return []UserStatus{
		UserStatusActive,
		UserStatusBlock,
	}
}
