package models


//// GroupCon

type GroupsGetRequest struct {
	Email string `json:"email"`
}

type GroupsGetResponse struct {
	Groups []Group
}

type GroupGetRequest struct {
	Email string `json:"email"`
	Id    uint   `json:"id"`
}

type GroupGetResponse struct {
	Users             []UserWithRoles
	Dashboards        []DashboardWithRoles
	GroupPermissionsG []CasbinRule
	GroupPermissionsP []CasbinRule
}

type GroupDeleteResponse struct {
	Messages []string `json:"message"`
}

type GroupsGetRawRequest struct {
	Email string `json:"email"`
}

type GroupsGetRawResponse struct {
	Groups []Group
}

type GroupAddRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type GroupAddResponse struct {
	Messages []string `json:"message"`
}

