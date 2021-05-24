package utils

import (
	"lalela-backend/internal/pkg/middleware"
	"lalela-backend/internal/pkg/models"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
)

func InitSeed(d *gorm.DB) {
	AutoMigrations(d)
	//AutoSeedRoles()
	//AutoSeedActions()
	//AutoSeedSiteConfig()

	//// Uncomment for Examples Seeding
	//AutoSeedExamples()
}

func AutoMigrations(d *gorm.DB) {
	//d.AutoMigrate(&models.Action{})
	//d.AutoMigrate(&models.Dashboard{})
	//d.AutoMigrate(&models.Group{})
	//d.AutoMigrate(&models.Role{})
	//d.AutoMigrate(&models.User{})
	//d.AutoMigrate(&models.Event{})
	//d.AutoMigrate(&models.DateProcessedItem{})
	//d.AutoMigrate(&models.DateProcessedTable{})
	//d.AutoMigrate(&models.SiteConfig{})

	d.AutoMigrate(&models.KanbanMember{})
	d.AutoMigrate(&models.KanbanList{})
	d.AutoMigrate(&models.KanbanCard{})
	d.AutoMigrate(&models.KanbanEvents{})
	d.AutoMigrate(&models.KanbanAssignLookUp{})
}

func AutoSeedSiteConfig() {

	db := GetDB()

	if err := db.Where("rule = ?", "isBlocked").First(&models.SiteConfig{}).Error; gorm.IsRecordNotFoundError(err) {
		r := models.SiteConfig{
			Rule:  "isBlocked",
			Value: false,
		}
		db.Create(&r)
	}
}

func AutoSeedRoles() {

	db := GetDB()

	var roles = [6]string{
		"Root",
		"Administrator",
		"Sub Administrator",
		"User",
		"Management",
		"Other",
	}

	for _, role := range roles {
		if err := db.Where("name = ?", role).First(&models.Role{}).Error; gorm.IsRecordNotFoundError(err) {
			r := models.Role{Name: role}
			db.Create(&r)
		}
	}
}

func AutoSeedActions() {

	db := GetDB()

	var actions = [4]string{
		"Read",
		"Write",
		"Delete",
		"All"}

	for _, action := range actions {
		if err := db.Where("name = ?", action).First(&models.Action{}).Error; gorm.IsRecordNotFoundError(err) {
			r := models.Action{Name: action}
			db.Create(&r)
		}
	}
}

// Example Data

func AutoSeedExamples() {
	//AutoSeedExampleUsers()
	////AutoSeedExampleDashboards()
	//AutoSeedExampleGroups()
	//AutoSeedExamplePermissions()
}

func AutoSeedExampleUsers() {

	db := GetDB()

	var Users = [2]models.User{
		{
			FirstName:    "noreply",
			LastName:     "FF",
			Email:        "noreply@cognizance.vision",
			Password:     "abc",
			UserGroupId:  1,
			RoleID:       0,
			ClientID:     0,
			CategoryName: "FF",
			IsTrail:      false,
		},
		{
			FirstName:    "Administartor",
			LastName:     "FF",
			Email:        "charl@cognizance.vision",
			Password:     "abc",
			UserGroupId:  1,
			RoleID:       0,
			ClientID:     0,
			CategoryName: "FF",
			IsTrail:      false,
		},
		//{
		//	FirstName:    "Client 1",
		//	LastName:     "User",
		//	Email:        "user@client1",
		//	Password:     "abc",
		//	UserGroupId:  2,
		//	RoleID:       0,
		//	ClientID:     0,
		//	CategoryName: "Company A",
		//	IsTrail:      false,
		//},
		//{
		//	FirstName:    "Client 1",
		//	LastName:     "Sub Admin",
		//	Email:        "subadmin@client1",
		//	Password:     "abc",
		//	UserGroupId:  2,
		//	RoleID:       0,
		//	ClientID:     0,
		//	CategoryName: "Company A",
		//	IsTrail:      false,
		//},
		//{
		//	FirstName:    "Client 2",
		//	LastName:     "User",
		//	Email:        "user@client2",
		//	Password:     "abc",
		//	UserGroupId:  3,
		//	RoleID:       0,
		//	ClientID:     0,
		//	CategoryName: "Company B",
		//	IsTrail:      false,
		//},
		//{
		//	FirstName:    "Client 2",
		//	LastName:     "Trail User",
		//	Email:        "trailuser@client2",
		//	Password:     "abc",
		//	UserGroupId:  3,
		//	RoleID:       0,
		//	ClientID:     0,
		//	CategoryName: "Company B",
		//	IsTrail:      true,
		//},
	}

	for _, user := range Users {
		if err := db.Where("email = ?", user.Email).First(&models.User{}).Error; gorm.IsRecordNotFoundError(err) {
			pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Print(middleware.NewError(err))
			}
			user.Password = string(pass)
			db.Create(&user)
		}
	}
}

func AutoSeedExampleDashboards() {

	db := GetDB()

	var Dashboards = [4]models.Dashboard{
		{
			DashboardUrl:  "https://embed.chartio.com/d/future-fragment-1/production_pns_distribution_percentage_report/",
			DashboardName: "Distribution Report",
			ClientID:      2,
			DashboardId:   452112,
			OrgId:         51657,
			CategoryName:  "Company A",
			DisplayOrder:  1,
		},
		{
			DashboardUrl:  "www.acompany.co.za",
			DashboardName: "Company A2",
			ClientID:      2,
			DashboardId:   0,
			OrgId:         0,
			CategoryName:  "Company A",
			DisplayOrder:  2,
		},
		{
			DashboardUrl:  "www.acompany.co.za",
			DashboardName: "Company B",
			ClientID:      3,
			DashboardId:   0,
			OrgId:         0,
			CategoryName:  "Company B",
			DisplayOrder:  2,
		},
		{
			DashboardUrl:  "www.acompany.co.za",
			DashboardName: "Company B2",
			ClientID:      3,
			DashboardId:   0,
			OrgId:         0,
			CategoryName:  "Company B",
			DisplayOrder:  1,
		},
	}

	for _, dashboard := range Dashboards {
		if err := db.Where("dashboard_url = ?", dashboard.DashboardUrl).First(&models.Dashboard{}).Error; gorm.IsRecordNotFoundError(err) {
			db.Create(&dashboard)
		}
	}
}

func AutoSeedExampleGroups() {

	db := GetDB()

	var Groups = [3]models.Group{
		{
			Name: "FF",
		},
		{
			Name: "Company A",
		},
		{
			Name: "Company B",
		},
	}

	for _, group := range Groups {
		if err := db.Where("name = ?", group.Name).First(&models.Group{}).Error; gorm.IsRecordNotFoundError(err) {
			db.Create(&group)
		}
	}
}

func AutoSeedExamplePermissions() {
	gs := GetGroups()
	ds := GetDashboards()
	us := GetUsersRaw()

	for _, dashboard := range ds {

		t := SliceContainsGroup(gs, strings.ToLower(dashboard.CategoryName))

		if t != 0 {
			PolicyExistsAdd("role::4", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "1")
		}
	}

	for _, user := range us {
		if user.UserGroupId == 1 {
			for _, group := range gs {
				GroupPolicyExistsAdd("user::"+fmt.Sprint(user.ID), "role::2", fmt.Sprint(group.ID))
			}
		} else {
			GroupPolicyExistsAdd("user::"+fmt.Sprint(user.ID), "role::4", fmt.Sprint(user.UserGroupId))
		}
	}
}

// General Functions
func GetGroups() []models.Group {
	db := GetDB()
	var groups []models.Group
	db.Find(&groups)
	return groups
}

func GetDashboards() []models.Dashboard {
	db := GetDB()
	var dashboards []models.Dashboard
	db.Find(&dashboards)
	return dashboards
}

func GetUsersRaw() []models.User {
	db := GetDB()
	var users []models.User
	db.Find(&users)
	return users
}

func SliceContainsGroup(a []models.Group, x string) uint {
	var r uint
	for _, n := range a {
		m := strings.ToLower(n.Name)
		if x == m {
			return n.ID
		} else {
			r = 0
		}
	}
	return r
}
