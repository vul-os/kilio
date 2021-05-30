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

