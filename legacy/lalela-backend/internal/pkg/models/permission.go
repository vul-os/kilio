package models


//// PermissionCon

type GetUserRolesRequest struct {
	Email  string `json:"email"`
	UserId uint   `json:"userid"`
}

type GetUserRolesResponse struct {
	Roles map[string][]string `json:"roles"`
}

type GetDashboardPermissionsRequest struct {
	Email  string `json:"email"`
	DashId uint   `json:"dashid"`
}

type GetDashboardPermissionsResponse struct {
	Permissions map[string][]string `json:"permissions"`
}

type GetGroupPermissionsRequest struct {
	Email   string `json:"email"`
	GroupId uint   `json:"groupId"`
}

type GetGroupPermissionsResponse struct {
	Permissions map[string][]string `json:"permissions"`
}

type GetUserPermissionsRequest struct {
	Email  string `json:"email"`
	UserId uint   `json:"userid"`
}

type GetUserPermissionsResponse struct {
	Permissions []string `json:"permissions"`
}

