package enum

type ReviewAction uint

const (
	ReviewActionApproved ReviewAction = iota + 1
	ReviewActionRejected
	ReviewActionSuspended
	ReviewActionAll
)

func (action ReviewAction) String() string {
	switch action {
	case ReviewActionApproved:
		return "تایید شده"
	case ReviewActionRejected:
		return "رد شده"
	case ReviewActionSuspended:
		return "معلق"
	case ReviewActionAll:
		return "همه"
	}
	return ""
}

func GetAllReviewActions() []ReviewAction {
	return []ReviewAction{
		ReviewActionApproved,
		ReviewActionRejected,
		ReviewActionSuspended,
		ReviewActionAll,
	}
}
