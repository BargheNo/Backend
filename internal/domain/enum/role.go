package enum

type RoleName uint

const (
	SuperAdmin RoleName = iota + 1
	Customer
	Technician
	CorporationManager
	SupportAgent
	ContentManager
	Moderator
	CorporationOwner
)

var rolePermissions = map[RoleName][]PermissionType{
	SuperAdmin: {
		PermissionAll,
	},
	Customer: {},
	Technician: {
		BidViewInstallationRequests, BidViewAll, BidCreate, BidEdit, BidCancel,
		MaintenanceViewAll, MaintenanceAcceptRequest, MaintenanceCreateRecord, MaintenanceUpdateRecord,
	},
	CorporationManager: {
		PanelViewAll, PanelCreate,
		BidViewInstallationRequests, BidViewAll, BidCreate, BidEdit, BidCancel,
		MaintenanceViewAll, MaintenanceAcceptRequest, MaintenanceCreateRecord, MaintenanceUpdateRecord,
		GuaranteeViewAll, GuaranteeCreate, GuaranteeArchiveUnarchive,
		CorpBlogViewAll, CorpBlogCreate, CorpBlogEdit, CorpBlogDelete,
		ProfileViewPrivate, ProfileUpdate,
	},
	SupportAgent: {
		TicketViewAll, TicketRespond, TicketClose, TicketComment,
		ReportViewAll, ReportRespond,
	},
	ContentManager: {
		NewsViewAll, NewsCreate, NewsEdit, NewsDelete,
	},
	Moderator: {
		CorpBlogViewAll, CorpBlogCreate, CorpBlogEdit, CorpBlogDelete,
	},
	CorporationOwner: {
		CorporationPermissionAll,
	},
}

func (role RoleName) Permissions() []PermissionType {
	if permissions, ok := rolePermissions[role]; ok {
		return permissions
	}
	return nil
}

func (role RoleName) String() string {
	switch role {
	case SuperAdmin:
		return "سوپر ادمین"
	case Customer:
		return "مشتری"
	case Technician:
		return "تکنسین"
	case CorporationManager:
		return "مدیر سازمان"
	case SupportAgent:
		return "پشتیبان"
	case ContentManager:
		return "مدیر محتوا"
	case Moderator:
		return "ناظر"
	case CorporationOwner:
		return "صاحب شرکت"
	}
	return "unknown"
}

func GetAllRoleNames() []RoleName {
	return []RoleName{
		SuperAdmin,
		Customer,
		Technician,
		CorporationManager,
		SupportAgent,
		ContentManager,
		Moderator,
		CorporationOwner,
	}
}

func GetAdminRoleNames() []RoleName {
	return []RoleName{
		SuperAdmin,
		Customer,
		SupportAgent,
		ContentManager,
	}
}

func GetCorporationRoleNames() []RoleName {
	return []RoleName{
		Technician,
		CorporationManager,
		Moderator,
		CorporationOwner,
	}
}
