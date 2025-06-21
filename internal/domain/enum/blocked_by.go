package enum

type BlockedBy uint

const (
	BlockedByUser BlockedBy = iota + 1
	BlockedByCorporation
)

func (by BlockedBy) String() string {
	switch by {
	case BlockedByUser:
		return "کاربر"
	case BlockedByCorporation:
		return "شرکت"
	}
	return "user"
}

func GetAllBlockedBy() []BlockedBy {
	return []BlockedBy{
		BlockedByUser,
		BlockedByCorporation,
	}
}
