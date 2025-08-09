package enum

type UserStatus uint

const (
	UserStatusActive UserStatus = iota + 1
	UserStatusBlock
	UserStatusAll
)

func (status UserStatus) String() string {
	switch status {
	case UserStatusActive:
		return "فعال"
	case UserStatusBlock:
		return "لیست سیاه"
	case UserStatusAll:
		return "همه"
	}
	return ""
}

func GetAllUserStatus() []UserStatus {
	return []UserStatus{
		UserStatusActive,
		UserStatusBlock,
		UserStatusAll,
	}
}
