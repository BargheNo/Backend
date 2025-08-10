package enum

type PermissionType uint
type PermissionCategory uint

const (
	PermissionAll PermissionType = iota + 1
	PermissionGeneral

	// Admin Role Permissions
	// User Management
	UserViewAll
	UserBanUnban
	UserChangeRole
	UserViewRoles
	UserManageRolePermissions
	UserRemoveRole
	UserCreateRole

	// Corporation Management
	CorporationViewAll
	CorporationApproveDecline

	// Panel Installation Request Management
	InstallationRequestViewAll
	InstallationRequestEdit
	InstallationRequestRemove

	// Ticket Management
	TicketViewAll
	TicketRespond
	TicketClose
	TicketComment

	// Report Management
	ReportViewAll
	ReportRespond

	// Blog Management (Admin)
	AdminBlogViewAll
	AdminBlogCreate
	AdminBlogEdit
	AdminBlogDelete

	// News Management
	NewsViewAll
	NewsCreate
	NewsEdit
	NewsDelete

	// Corporation Role Permissions
	// Panel Management
	PanelViewAll
	PanelCreate

	// Bid Management
	BidViewInstallationRequests
	BidViewAll
	BidCreate
	BidEdit
	BidCancel

	// Maintenance Management
	MaintenanceViewAll
	MaintenanceAcceptRequest
	MaintenanceCreateRecord
	MaintenanceUpdateRecord

	// Guarantee Management
	GuaranteeViewAll
	GuaranteeCreate
	GuaranteeArchiveUnarchive

	// Blog Management (Corporation)
	CorpBlogViewAll
	CorpBlogCreate
	CorpBlogEdit
	CorpBlogDelete

	// Profile Management
	ProfileViewPrivate
	ProfileUpdate
)

const (
	CategoryGeneral PermissionCategory = iota + 1
	CategoryUser
	CategoryCorporation
	CategoryInstallationRequest
	CategoryTicket
	CategoryReport
	CategoryBlog
	CategoryNews
	CategoryPanel
	CategoryBid
	CategoryMaintenance
	CategoryGuarantee
	CategoryProfile
)

var permissionNames = map[PermissionType]string{
	PermissionAll:     "general.all",
	PermissionGeneral: "general.general",

	// User Management
	UserViewAll:               "user.viewAll",
	UserBanUnban:              "user.banUnban",
	UserChangeRole:            "user.changeRole",
	UserViewRoles:             "user.viewRoles",
	UserManageRolePermissions: "user.manageRolePermissions",
	UserRemoveRole:            "user.removeRole",
	UserCreateRole:            "user.createRole",

	// Corporation Management
	CorporationViewAll:        "corporation.viewAll",
	CorporationApproveDecline: "corporation.approveDecline",

	// Installation Request Management
	InstallationRequestViewAll: "installationRequest.viewAll",
	InstallationRequestEdit:    "installationRequest.edit",
	InstallationRequestRemove:  "installationRequest.remove",

	// Ticket Management
	TicketViewAll: "ticket.viewAll",
	TicketRespond: "ticket.respond",
	TicketClose:   "ticket.close",
	TicketComment: "ticket.comment",

	// Report Management
	ReportViewAll: "report.viewAll",
	ReportRespond: "report.respond",

	// Blog Management (Admin)
	AdminBlogViewAll: "adminBlog.viewAll",
	AdminBlogCreate:  "adminBlog.create",
	AdminBlogEdit:    "adminBlog.edit",
	AdminBlogDelete:  "adminBlog.delete",

	// News Management
	NewsViewAll: "news.viewAll",
	NewsCreate:  "news.create",
	NewsEdit:    "news.edit",
	NewsDelete:  "news.delete",

	// Panel Management
	PanelViewAll: "panel.viewAll",
	PanelCreate:  "panel.create",

	// Bid Management
	BidViewInstallationRequests: "bid.viewInstallationRequests",
	BidViewAll:                  "bid.viewAll",
	BidCreate:                   "bid.create",
	BidEdit:                     "bid.edit",
	BidCancel:                   "bid.cancel",

	// Maintenance Management
	MaintenanceViewAll:       "maintenance.viewAll",
	MaintenanceAcceptRequest: "maintenance.acceptRequest",
	MaintenanceCreateRecord:  "maintenance.createRecord",
	MaintenanceUpdateRecord:  "maintenance.updateRecord",

	// Guarantee Management
	GuaranteeViewAll:          "guarantee.viewAll",
	GuaranteeCreate:           "guarantee.create",
	GuaranteeArchiveUnarchive: "guarantee.archiveUnarchive",

	// Blog Management (Corporation)
	CorpBlogViewAll: "corpBlog.viewAll",
	CorpBlogCreate:  "corpBlog.create",
	CorpBlogEdit:    "corpBlog.edit",
	CorpBlogDelete:  "corpBlog.delete",

	// Profile Management
	ProfileViewPrivate: "profile.viewPrivate",
	ProfileUpdate:      "profile.update",
}

var permissionDescriptions = map[PermissionType]string{
	PermissionAll:     "دسترسی کامل به سیستم",
	PermissionGeneral: "دسترسی عمومی",

	// User Management
	UserViewAll:               "مشاهده لیست کاربران وب‌سایت",
	UserBanUnban:              "مسدود/رفع مسدودیت کاربران",
	UserChangeRole:            "تغییر نقش کاربر",
	UserViewRoles:             "مشاهده لیست نقش‌ها و مجوزها",
	UserManageRolePermissions: "تغییر مجوزهای نقش",
	UserRemoveRole:            "حذف نقش",
	UserCreateRole:            "ایجاد نقش جدید",

	// Corporation Management
	CorporationViewAll:        "مشاهده لیست شرکت‌ها",
	CorporationApproveDecline: "تایید/رد درخواست شرکت",

	// Installation Request Management
	InstallationRequestViewAll: "مشاهده لیست درخواست‌های نصب",
	InstallationRequestEdit:    "ویرایش درخواست نصب",
	InstallationRequestRemove:  "حذف درخواست نصب",

	// Ticket Management
	TicketViewAll: "مشاهده لیست تیکت‌ها",
	TicketRespond: "پاسخ به تیکت",
	TicketClose:   "بستن تیکت",
	TicketComment: "افزودن نظر به تیکت",

	// Report Management
	ReportViewAll: "مشاهده لیست گزارش‌ها",
	ReportRespond: "پاسخ به گزارش",

	// Blog Management (Admin)
	AdminBlogViewAll: "مشاهده لیست بلاگ‌ها",
	AdminBlogCreate:  "ایجاد بلاگ جدید",
	AdminBlogEdit:    "ویرایش بلاگ",
	AdminBlogDelete:  "حذف بلاگ",

	// News Management
	NewsViewAll: "مشاهده لیست اخبار",
	NewsCreate:  "ایجاد خبر جدید",
	NewsEdit:    "ویرایش خبر",
	NewsDelete:  "حذف خبر",

	// Panel Management
	PanelViewAll: "مشاهده لیست پنل‌ها",
	PanelCreate:  "ایجاد پنل جدید",

	// Bid Management
	BidViewInstallationRequests: "مشاهده لیست درخواست‌های نصب",
	BidViewAll:                  "مشاهده لیست پیشنهادات",
	BidCreate:                   "ایجاد پیشنهاد جدید",
	BidEdit:                     "ویرایش پیشنهاد",
	BidCancel:                   "لغو پیشنهاد",

	// Maintenance Management
	MaintenanceViewAll:       "مشاهده لیست درخواست‌های تعمیر",
	MaintenanceAcceptRequest: "پذیرش/لغو درخواست تعمیر",
	MaintenanceCreateRecord:  "ایجاد سابقه تعمیر",
	MaintenanceUpdateRecord:  "به‌روزرسانی سابقه تعمیر",

	// Guarantee Management
	GuaranteeViewAll:          "مشاهده لیست ضمانت‌ها",
	GuaranteeCreate:           "ایجاد ضمانت جدید",
	GuaranteeArchiveUnarchive: "آرشیو/رفع آرشیو ضمانت",

	// Blog Management (Corporation)
	CorpBlogViewAll: "مشاهده لیست بلاگ‌ها",
	CorpBlogCreate:  "ایجاد بلاگ جدید",
	CorpBlogEdit:    "ویرایش بلاگ",
	CorpBlogDelete:  "حذف بلاگ",

	// Profile Management
	ProfileViewPrivate: "مشاهده اطلاعات خصوصی پروفایل",
	ProfileUpdate:      "به‌روزرسانی پروفایل",
}

var permissionCategories = map[PermissionType]PermissionCategory{
	PermissionAll:     CategoryGeneral,
	PermissionGeneral: CategoryGeneral,

	// User Management
	UserViewAll:               CategoryUser,
	UserBanUnban:              CategoryUser,
	UserChangeRole:            CategoryUser,
	UserViewRoles:             CategoryUser,
	UserManageRolePermissions: CategoryUser,
	UserRemoveRole:            CategoryUser,
	UserCreateRole:            CategoryUser,

	// Corporation Management
	CorporationViewAll:        CategoryCorporation,
	CorporationApproveDecline: CategoryCorporation,

	// Installation Request Management
	InstallationRequestViewAll: CategoryInstallationRequest,
	InstallationRequestEdit:    CategoryInstallationRequest,
	InstallationRequestRemove:  CategoryInstallationRequest,

	// Ticket Management
	TicketViewAll: CategoryTicket,
	TicketRespond: CategoryTicket,
	TicketClose:   CategoryTicket,
	TicketComment: CategoryTicket,

	// Report Management
	ReportViewAll: CategoryReport,
	ReportRespond: CategoryReport,

	// Blog Management (Admin)
	AdminBlogViewAll: CategoryBlog,
	AdminBlogCreate:  CategoryBlog,
	AdminBlogEdit:    CategoryBlog,
	AdminBlogDelete:  CategoryBlog,

	// News Management
	NewsViewAll: CategoryNews,
	NewsCreate:  CategoryNews,
	NewsEdit:    CategoryNews,
	NewsDelete:  CategoryNews,

	// Panel Management
	PanelViewAll: CategoryPanel,
	PanelCreate:  CategoryPanel,

	// Bid Management
	BidViewInstallationRequests: CategoryBid,
	BidViewAll:                  CategoryBid,
	BidCreate:                   CategoryBid,
	BidEdit:                     CategoryBid,
	BidCancel:                   CategoryBid,

	// Maintenance Management
	MaintenanceViewAll:       CategoryMaintenance,
	MaintenanceAcceptRequest: CategoryMaintenance,
	MaintenanceCreateRecord:  CategoryMaintenance,
	MaintenanceUpdateRecord:  CategoryMaintenance,

	// Guarantee Management
	GuaranteeViewAll:          CategoryGuarantee,
	GuaranteeCreate:           CategoryGuarantee,
	GuaranteeArchiveUnarchive: CategoryGuarantee,

	// Blog Management (Corporation)
	CorpBlogViewAll: CategoryBlog,
	CorpBlogCreate:  CategoryBlog,
	CorpBlogEdit:    CategoryBlog,
	CorpBlogDelete:  CategoryBlog,

	// Profile Management
	ProfileViewPrivate: CategoryProfile,
	ProfileUpdate:      CategoryProfile,
}

func (perm PermissionType) String() string {
	if description, ok := permissionNames[perm]; ok {
		return description
	}
	return ""
}

func (perm PermissionType) Description() string {
	if description, ok := permissionDescriptions[perm]; ok {
		return description
	}
	return ""
}

func (perm PermissionType) Category() PermissionCategory {
	if category, ok := permissionCategories[perm]; ok {
		return category
	}
	return CategoryGeneral
}

func (category PermissionCategory) String() string {
	switch category {
	case CategoryGeneral:
		return "عمومی"
	case CategoryUser:
		return "مدیریت کاربران"
	case CategoryCorporation:
		return "مدیریت شرکت‌ها"
	case CategoryInstallationRequest:
		return "مدیریت درخواست‌های نصب"
	case CategoryTicket:
		return "مدیریت تیکت"
	case CategoryReport:
		return "مدیریت گزارشات"
	case CategoryBlog:
		return "مدیریت بلاگ"
	case CategoryNews:
		return "مدیریت اخبار"
	case CategoryPanel:
		return "مدیریت پنل‌ها"
	case CategoryBid:
		return "مدیریت پیشنهادات"
	case CategoryMaintenance:
		return "مدیریت تعمیرات"
	case CategoryGuarantee:
		return "مدیریت ضمانت‌ها"
	case CategoryProfile:
		return "مدیریت پروفایل"
	}
	return "unknown"
}

func GetAllPermissionTypes() []PermissionType {
	return []PermissionType{
		PermissionAll, PermissionGeneral,

		// User Management
		UserViewAll, UserBanUnban, UserChangeRole, UserViewRoles,
		UserManageRolePermissions, UserRemoveRole, UserCreateRole,

		// Corporation Management
		CorporationViewAll, CorporationApproveDecline,

		// Installation Request Management
		InstallationRequestViewAll, InstallationRequestEdit, InstallationRequestRemove,

		// Ticket Management
		TicketViewAll, TicketRespond, TicketClose, TicketComment,

		// Report Management
		ReportViewAll, ReportRespond,

		// Blog Management (Admin)
		AdminBlogViewAll, AdminBlogCreate, AdminBlogEdit, AdminBlogDelete,

		// News Management
		NewsViewAll, NewsCreate, NewsEdit, NewsDelete,

		// Panel Management
		PanelViewAll, PanelCreate,

		// Bid Management
		BidViewInstallationRequests, BidViewAll, BidCreate, BidEdit, BidCancel,

		// Maintenance Management
		MaintenanceViewAll, MaintenanceAcceptRequest, MaintenanceCreateRecord, MaintenanceUpdateRecord,

		// Guarantee Management
		GuaranteeViewAll, GuaranteeCreate, GuaranteeArchiveUnarchive,

		// Blog Management (Corporation)
		CorpBlogViewAll, CorpBlogCreate, CorpBlogEdit, CorpBlogDelete,

		// Profile Management
		ProfileViewPrivate, ProfileUpdate,
	}
}

func GetCorporationPermissionTypes() []PermissionType {
	return []PermissionType{
		// Panel Management
		PanelViewAll, PanelCreate,

		// Bid Management
		BidViewInstallationRequests, BidViewAll, BidCreate, BidEdit, BidCancel,

		// Maintenance Management
		MaintenanceViewAll, MaintenanceAcceptRequest, MaintenanceCreateRecord, MaintenanceUpdateRecord,

		// Guarantee Management
		GuaranteeViewAll, GuaranteeCreate, GuaranteeArchiveUnarchive,

		// Blog Management (Corporation)
		CorpBlogViewAll, CorpBlogCreate, CorpBlogEdit, CorpBlogDelete,

		// Profile Management
		ProfileViewPrivate, ProfileUpdate,
	}
}
