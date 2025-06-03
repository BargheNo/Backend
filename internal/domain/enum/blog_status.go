package enum

type PostStatus uint

const (
	PostStatusDraft PostStatus = iota + 1
	PostStatusPublished
)

func (status PostStatus) String() string {
	switch status {
	case PostStatusDraft:
		return "پیش نویس"
	case PostStatusPublished:
		return "منتشر شده"
	}
	return ""
}

func GetAllPostStatus() []PostStatus {
	return []PostStatus{
		PostStatusDraft,
		PostStatusPublished,
	}
}
