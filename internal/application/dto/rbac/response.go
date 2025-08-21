package rbacdto

type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	IsCorpStaff bool                 `json:"isCorpStaff"`
	Permissions []PermissionResponse `json:"permissions"`
}

type PermissionResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IsCorpStaff bool   `json:"isCorpStaff"`
	Category    string `json:"category"`
}
