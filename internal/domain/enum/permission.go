package enum

type PermissionType uint
type PermissionCategory uint

const (
	PermissionAll PermissionType = iota + 1
	PermissionGeneral
	UserViewAll
	UserViewCustomers
	UserApproveCorporation
	UserRejectCorporation
	UserManageRolesPermissions
	TicketViewAll
	TicketCreate
	TicketViewOwn
	TicketRespond
	ReportViewAll
	ReportManage
	ReportCreate
	ReportViewOwn
	BlogCreate
	BlogEdit
	BlogDelete
	BlogView
	BlogComment
	NewsCreate
	NewsEdit
	NewsDelete
	NewsView
	NewsComment
	PanelViewAll
	PanelAssignToCustomer
	PanelRemove
	PanelViewAssigned
	PanelRequestPurchase
	PanelViewOwned
	PanelViewUsageData
	ChatViewAll
	ChatSendMessage
	SystemSettingsUpdate
	RepairViewAll
	RepairAssignExpert
	RepairAcceptRequest
	RepairMarkComplete
	RepairRequest
	RepairViewAssigned
	AnalyticsViewAll
	ProfileUpdate
	CorporationViewAll
	CorporationManage
	CorporationViewStats
	CorporationUpdateSettings
)

const (
	CategoryGeneral PermissionCategory = iota + 1
	CategoryUser
	CategoryTicket
	CategoryReport
	CategoryBlog
	CategoryNews
	CategoryPanel
	CategoryChat
	CategorySystem
	CategoryRepair
	CategoryAnalytics
	CategoryProfile
	CategoryCorporation
)

var permissionNames = map[PermissionType]string{
	PermissionAll:              "general.all",
	PermissionGeneral:          "general.general",
	UserViewAll:                "user.view_all",
	UserViewCustomers:          "user.view_customers",
	UserApproveCorporation:     "user.approve_Corporation",
	UserRejectCorporation:      "user.reject_Corporation",
	UserManageRolesPermissions: "user.manage_roles_permissions",
	TicketViewAll:              "ticket.view_all",
	TicketCreate:               "ticket.create",
	TicketViewOwn:              "ticket.view_own",
	TicketRespond:              "ticket.respond",
	ReportViewAll:              "report.view_all",
	ReportManage:               "report.manage",
	ReportCreate:               "report.create",
	ReportViewOwn:              "report.view_own",
	BlogCreate:                 "blog.create",
	BlogEdit:                   "blog.edit",
	BlogDelete:                 "blog.delete",
	BlogView:                   "blog.view",
	BlogComment:                "blog.comment",
	NewsCreate:                 "news.create",
	NewsEdit:                   "news.edit",
	NewsDelete:                 "news.delete",
	NewsView:                   "news.view",
	NewsComment:                "news.comment",
	PanelViewAll:               "panel.view_all",
	PanelAssignToCustomer:      "panel.assign_to_customer",
	PanelRemove:                "panel.remove",
	PanelViewAssigned:          "panel.view_assigned",
	PanelRequestPurchase:       "panel.request_purchase",
	PanelViewOwned:             "panel.view_owned",
	PanelViewUsageData:         "panel.view_usage_data",
	ChatViewAll:                "chat.view_all",
	ChatSendMessage:            "chat.send_message",
	SystemSettingsUpdate:       "systemSettings.update",
	RepairViewAll:              "repair.view_all",
	RepairAssignExpert:         "repair.assign_expert",
	RepairAcceptRequest:        "repair.accept_request",
	RepairMarkComplete:         "repair.mark_complete",
	RepairRequest:              "repair.request",
	RepairViewAssigned:         "repair.view_assigned",
	AnalyticsViewAll:           "analytics.view_all",
	ProfileUpdate:              "profile.update",
	CorporationViewAll:         "Corporation.view_all",
	CorporationManage:          "Corporation.manage",
	CorporationViewStats:       "Corporation.view_stats",
	CorporationUpdateSettings:  "Corporation.update_settings",
}

var permissionDescriptions = map[PermissionType]string{
	PermissionAll:              "دسترسی کامل به سیستم",
	PermissionGeneral:          "دسترسی عمومی",
	UserViewAll:                "مشاهده تمامی کاربران سیستم",
	UserViewCustomers:          "مشاهده تمامی مشتریان شرکت",
	UserApproveCorporation:     "تایید فروشندگان و تامین کنندگان",
	UserRejectCorporation:      "رد درخواست فروشندگان و تامین کنندگان",
	UserManageRolesPermissions: "مدیریت نقش و دسترسی کاربران",
	TicketViewAll:              "مدیریت تمامی تیکت های سیستم",
	TicketCreate:               "ایجاد تیکت جدید",
	TicketViewOwn:              "مشاهده تیکت های خود",
	TicketRespond:              " پاسخ دادن به تیکت ها",
	ReportViewAll:              "مشاهده تمامی گزارش های سیستم",
	ReportManage:               "مدیریت گزارش های سیستم",
	ReportCreate:               "ایجاد گزارش جدید",
	ReportViewOwn:              "مشاهده گزارش های خود",
	BlogCreate:                 "ایجاد مطلب جدید در بلاگ",
	BlogEdit:                   "ویرایش مطالب بلاگ",
	BlogDelete:                 "حذف مطالب بلاگ",
	BlogView:                   "مشاهده مطالب بلاگ",
	BlogComment:                "ثبت نظر برای مطالب بلاگ",
	NewsCreate:                 "ایجاد خبر جدید",
	NewsEdit:                   "ویرایش اخبار",
	NewsDelete:                 "حذف اخبار",
	NewsView:                   "مشاهده اخبار",
	NewsComment:                "ثبت نظر برای اخبار",
	PanelViewAll:               "مشاهده تمامی پنل های سیستم",
	PanelAssignToCustomer:      "اختصاص پنل به مشتری",
	PanelRemove:                "حذف پنل از سیستم",
	PanelViewAssigned:          "مشاهده پنل های اختصاص داده شده",
	PanelRequestPurchase:       "درخواست خرید پنل",
	PanelViewOwned:             "مشاهده پنل های متعلق به خود",
	PanelViewUsageData:         "مشاهده داده های مصرف پنل",
	ChatViewAll:                "مشاهده تمامی مکالمات چت",
	ChatSendMessage:            "ارسال پیام در چت",
	SystemSettingsUpdate:       "به روزرسانی تنظیمات سیستم",
	RepairViewAll:              "مشاهده تمامی درخواست های تعمیر",
	RepairAssignExpert:         "اختصاص متخصص به درخواست تعمیر",
	RepairAcceptRequest:        "پذیرش درخواست تعمیر",
	RepairMarkComplete:         "علامت گذاری تعمیر به عنوان تکمیل شده",
	RepairRequest:              "ثبت درخواست تعمیر",
	RepairViewAssigned:         "مشاهده درخواست های تعمیر اختصاص داده شده",
	AnalyticsViewAll:           "مشاهده تمامی تحلیل ها و آمار سیستم",
	ProfileUpdate:              "به روزرسانی پروفایل کاربری",
	CorporationViewAll:         "مشاهده تمامی فروشگاه ها",
	CorporationManage:          "مدیریت کامل فروشگاه",
	CorporationViewStats:       "مشاهده گزارشات فروشگاه",
	CorporationUpdateSettings:  "به روزرسانی تنظیمات فروشگاه",
}

var permissionCategories = map[PermissionType]PermissionCategory{
	PermissionAll:     CategoryGeneral,
	PermissionGeneral: CategoryGeneral,

	UserViewAll:                CategoryUser,
	UserViewCustomers:          CategoryUser,
	UserApproveCorporation:     CategoryUser,
	UserRejectCorporation:      CategoryUser,
	UserManageRolesPermissions: CategoryUser,

	TicketViewAll: CategoryTicket,
	TicketCreate:  CategoryTicket,
	TicketViewOwn: CategoryTicket,
	TicketRespond: CategoryTicket,

	ReportViewAll: CategoryReport,
	ReportManage:  CategoryReport,
	ReportCreate:  CategoryReport,
	ReportViewOwn: CategoryReport,

	BlogCreate:  CategoryBlog,
	BlogEdit:    CategoryBlog,
	BlogDelete:  CategoryBlog,
	BlogView:    CategoryBlog,
	BlogComment: CategoryBlog,

	NewsCreate:  CategoryNews,
	NewsEdit:    CategoryNews,
	NewsDelete:  CategoryNews,
	NewsView:    CategoryNews,
	NewsComment: CategoryNews,

	PanelViewAll:          CategoryPanel,
	PanelAssignToCustomer: CategoryPanel,
	PanelRemove:           CategoryPanel,
	PanelViewAssigned:     CategoryPanel,
	PanelRequestPurchase:  CategoryPanel,
	PanelViewOwned:        CategoryPanel,
	PanelViewUsageData:    CategoryPanel,

	ChatViewAll:     CategoryChat,
	ChatSendMessage: CategoryChat,

	SystemSettingsUpdate: CategorySystem,

	RepairViewAll:       CategoryRepair,
	RepairAssignExpert:  CategoryRepair,
	RepairAcceptRequest: CategoryRepair,
	RepairMarkComplete:  CategoryRepair,
	RepairRequest:       CategoryRepair,
	RepairViewAssigned:  CategoryRepair,

	AnalyticsViewAll: CategoryAnalytics,

	ProfileUpdate: CategoryProfile,

	CorporationViewAll:        CategoryCorporation,
	CorporationManage:         CategoryCorporation,
	CorporationViewStats:      CategoryCorporation,
	CorporationUpdateSettings: CategoryCorporation,
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
		return "general"
	case CategoryUser:
		return "user"
	case CategoryTicket:
		return "ticket"
	case CategoryReport:
		return "report"
	case CategoryBlog:
		return "blog"
	case CategoryNews:
		return "news"
	case CategoryPanel:
		return "panel"
	case CategoryChat:
		return "chat"
	case CategorySystem:
		return "system"
	case CategoryRepair:
		return "repair"
	case CategoryAnalytics:
		return "analytics"
	case CategoryProfile:
		return "profile"
	case CategoryCorporation:
		return "corporation"
	}
	return "unknown"
}

func GetAllPermissionTypes() []PermissionType {
	return []PermissionType{
		PermissionAll, PermissionGeneral,
		UserViewAll, UserViewCustomers, UserApproveCorporation, UserRejectCorporation, UserManageRolesPermissions,
		TicketViewAll, TicketCreate, TicketViewOwn, TicketRespond,
		ReportViewAll, ReportManage, ReportCreate, ReportViewOwn,
		BlogCreate, BlogEdit, BlogDelete, BlogView, BlogComment,
		NewsCreate, NewsEdit, NewsDelete, NewsView, NewsComment,
		PanelViewAll, PanelAssignToCustomer, PanelRemove, PanelViewAssigned, PanelRequestPurchase, PanelViewOwned, PanelViewUsageData,
		ChatViewAll, ChatSendMessage,
		SystemSettingsUpdate,
		RepairViewAll, RepairAssignExpert, RepairAcceptRequest, RepairMarkComplete, RepairRequest, RepairViewAssigned,
		AnalyticsViewAll,
		ProfileUpdate,
		CorporationViewAll, CorporationManage, CorporationViewStats, CorporationUpdateSettings,
	}
}
