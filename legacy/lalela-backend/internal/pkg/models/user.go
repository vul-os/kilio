package models

import "time"

//User struct declaration
type User struct {
	ID               uint      `gorm:"AUTO_INCREMENT"`
	FirstName        string    `json:"firstName"`
	LastName         string    `json:"lastName"`
	Email            string    `gorm:"type:varchar(100);unique_index"`
	Password         string    `json:"password"`
	UserGroupId      int       `json:"userGroupId"`
	RoleID           int       `json:"roleID"`
	ClientID         int       `json:"clientID"`
	CategoryName     string    `json:"categoryName"`
	Avatar           string    `json:"avatar"`
	IsTrail          bool      `json:"isTrail"`
	TrailExpire      time.Time `json:"-"`
	ResetToken       string    `json:"-"`
	ResetTokenExpiry time.Time `json:"-"`
	CreatedAt        time.Time `json:"-"`
	UpdatedAt        time.Time `json:"-"`
}

//UserGetReport struct declaration
type UserGetReport struct {
	ID          uint      `gorm:"AUTO_INCREMENT"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `gorm:"type:varchar(100);unique_index"`
	Password    string    `json:"password"`
	UserGroupId string    `json:"userGroupId"`
	Avatar      string    `json:"avatar"`
	IsTrail     bool      `json:"isTrail"`
	TrailExpire time.Time `json:"trailExpire"`
	LastLogin   string    `json:"lastLogin"`
}

//User struct declaration
type UserWithRoles struct {
	ID        uint     `gorm:"AUTO_INCREMENT"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Email     string   `gorm:"type:varchar(100);unique_index"`
	Avatar    string   `json:"avatar"`
	Roles     []string `json:"roles"`
}

//// UserCon

type UsersGetRequest struct {
	Email string `json:"email"`
}

type UsersGetResponse struct {
	Users []User
}

type UserGetRequest struct {
	Email string `json:"email"`
	Id    uint   `json:"id"`
}

type UserGetResponse struct {
	User User
}

type UserIsAdminRequest struct {
	Email string `json:"email"`
}

type UserIsAdminResponse struct {
	Admin bool
}

type UserGetResponse2 struct {
	User       UserGetReport
	Dashboards map[string][]string
}

type UserAddRequest struct {
	Email string `json:"email"`
	Users []User
}

type UserAddResponse struct {
	Messages []string `json:"message"`
}

type UserUpdateRequest struct {
	Email string `json:"email"`
	Users []User
}

type UserUpdateRequestSingle struct {
	Email string `json:"email"`
	User  User
}

type UserUpdateResponse struct {
	Messages []string `json:"message"`
}

type UserDeleteRequest struct {
	Email string `json:"email"`
	User  User
}

type UserDeleteResponse struct {
	Messages []string `json:"message"`
}

type UserRoleUpdateRequest struct {
	Email   string `json:"email"`
	User    User
	Roles   []string
	GroupId int
}

type UserRoleUpdateResponse struct {
	Messages []string `json:"message"`
}