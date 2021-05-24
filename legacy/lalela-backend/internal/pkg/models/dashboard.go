package models


type DashboardData struct {
	OrgName    string
	Dashboards []Dashboard
}

//Dashboard declaration
type Dashboard struct {
	Id            uint   `json:"id",gorm:"AUTO_INCREMENT"`
	DashboardUrl  string `json:"dashboardUrl",gorm:"type:varchar(255)"`
	DashboardName string `json:"dashboardName"`
	ClientID      int    `json:"clientID"`
	DashboardId   int    `json:"dashboardId"`
	OrgId         int    `json:"orgId"`
	CategoryName  string `json:"categoryName"`
	DisplayOrder  int    `json:"displayOrder"`
}

//Dashboard With Roles declaration
type DashboardWithRoles struct {
	Id            uint     `json:"id",gorm:"AUTO_INCREMENT"`
	DashboardUrl  string   `json:"dashboardUrl",gorm:"type:varchar(255)"`
	DashboardName string   `json:"dashboardName"`
	CategoryName  string   `json:"categoryName"`
	Roles         []string `json:"roles"`
}

type DashboardPermission struct {
	ID          uint   `gorm:"AUTO_INCREMENT"`
	DashboardID string `json:"dashboardID"`
	UserID      string `json:"userID"`
}


//// DashCon

type DashboardDataRequest struct {
	Email string `json:"email"`
}

type DashboardSingleDataRequest struct {
	Email string `json:"email"`
	Id    uint   `json:"id"`
}

type DashboardSingleDataResponse struct {
	DashboardData Dashboard `json:"dashboardData"`
}

//type DashboardDataResponse struct {
//	DashboardData []DashboardData `json:"dashboardData"`
//}

type DashboardDataHolder struct {
	OrganizationName string      `json:"organizationName"`
	Dashboards       []Dashboard `json:"dashboards"`
}

type DashboardDataResponse struct {
	DashboardData []DashboardDataHolder `json:"dashboardData"`
}

type DashboardsGetResponse struct {
	Dashboards []Dashboard `json:"dashboardData"`
}

type DashboardAddRequest struct {
	Email      string `json:"email"`
	Dashboards []Dashboard
}

type DashboardAddResponse struct {
	Messages []string `json:"message"`
}

type DashboardDeleteRequest struct {
	Email      string `json:"email"`
	Dashboards []Dashboard
}

type DashboardDeleteSingleRequest struct {
	Email     string `json:"email"`
	Dashboard uint
}

type DashboardUpdateSingleRequest struct {
	Email     string `json:"email"`
	Dashboard Dashboard
}

type DashboardDeleteResponse struct {
	Messages []string `json:"message"`
}

type DashboardUpdateRequest struct {
	Email      string `json:"Email"`
	Dashboards []Dashboard
}

type DashboardUpdateResponse struct {
	Messages []string `json:"message"`
}

type DashboardRolesRequest struct {
	Email   string    `json:"Email"`
	Dash    Dashboard `json:"Dash"`
	Roles   []string
	GroupId int
}

type DashboardRolesResponse struct {
	Messages []string `json:"message"`
}