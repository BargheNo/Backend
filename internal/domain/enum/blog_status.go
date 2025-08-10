package enum

type PostStatus uint

const (
	PostStatusDraft PostStatus = iota + 1
	PostStatusPublished
	PostStatusAll
)

func (status PostStatus) String() string {
	switch status {
	case PostStatusDraft:
		return "پیش نویس"
	case PostStatusPublished:
		return "منتشر شده"
	case PostStatusAll:
		return "همه"
	}
	return ""
}

func GetAllPostStatus() []PostStatus {
	return []PostStatus{
		PostStatusDraft,
		PostStatusPublished,
		PostStatusAll,
	}
}
