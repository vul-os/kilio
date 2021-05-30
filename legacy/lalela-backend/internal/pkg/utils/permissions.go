package utils

import (
	"cog-analytics-engine-go/internal/pkg/controllers"
	"cog-analytics-engine-go/internal/pkg/models"
	"fmt"
	"strconv"
)

type CasbinRule struct {
	PType string
	V0    string
	V1    string
	V2    string
	V3    string
	V4    string
	V5    string
}

func GetPermissionExists(permission *controllers.DashboardPermission) (bool, error) {
	db := GetPostgreDB()
	var dashboardPermissions controllers.DashboardPermission
	if len(permission.DashboardID) != 0 && len(permission.UserID) != 0 {
		if err := db.Where(&controllers.DashboardPermission{DashboardID: permission.DashboardID, UserID: permission.UserID}).Find(&dashboardPermissions).Error; err != nil {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return true, nil
	}
}

func AddPermission(permission *controllers.DashboardPermission) bool {
	db := GetPostgreDB()
	if permission.DashboardID == "" || permission.UserID == "" {
		return false
	} else {
		db.Create(&permission)
		return true
	}
	return false
}

func AppendUniqueSlice(slice []string, i string) []string {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}
	return append(slice, i)
}

// todo :: Check if Used
func GetGroups() []models.Group {
	db := GetPostgreDB()
	var groups []models.Group
	db.Find(&groups)
	return groups
}

func GetGroups2() map[int]string {
	db := GetPostgreDB()
	var groups []models.Group
	groupsT := make(map[int]string)
	db.Find(&groups)
	for _, r := range groups {
		groupsT[int(r.ID)] = r.Name
	}
	return groupsT
}

func GetRoles() map[string]string {
	db := GetPostgreDB()
	var roles []models.Role
	rolesT := make(map[string]string)
	db.Find(&roles)

	for _, r := range roles {
		t := "role::" + fmt.Sprint(r.ID)
		rolesT[t] = r.Name
	}
	return rolesT
}

func GetDashboardsGroup() map[uint]string {
	db := GetPostgreDB()
	var dashboards []controllers.Dashboard
	dashboardsT := make(map[uint]string)
	db.Find(&dashboards)

	for _, r := range dashboards {
		dashboardsT[r.Id] = r.CategoryName + ": " + r.DashboardName
	}
	return dashboardsT
}

func GetPermissionsGroup(g string, r string) {
	db := GetPostgreDB().Table("casbin_rule")

	var permissions []CasbinRule

	// Find User Groupings
	db.Where("V0 = ? AND V1 = ?", r, g).Find(&permissions)

	// Find All Permissions From Group

	return
}

// Grab All Permissions for Users
func GetPermissionsUser(u string) {
	db := GetPostgreDB().Table("casbin_rule")

	var permissions []CasbinRule
	var dashboards []string

	// Find User Groupings
	db.Where("V0 = ?", u).Find(&permissions)

	for _, perm := range permissions {
		db.Where("V0 = ? AND V1 = ? AND V3 = '1'", perm.V1, perm.V2).Find(&permissions)
		for _, perms := range permissions {
			dashboards = append(dashboards, perms.V2+" "+perms.V1)
		}
	}

	return
}

// Grab All Permissions For Dashboard
func GetDashboardPermissions(d string) map[string][]string {

	// DB
	db := GetPostgreDB().Table("casbin_rule")

	// Declares
	var permissions []CasbinRule
	dashT := make(map[string][]string)
	roles := GetRoles()
	groups := GetGroups2()

	// Find User Groupings
	db.Where("p_type = 'p' AND V2 = ?", d).Find(&permissions)

	// Find All Permissions From Group
	for _, p := range permissions {
		r := roles[p.V0]
		g, _ := strconv.Atoi(p.V1)
		dashT[r] = append(dashT[r], groups[g])
	}

	return dashT
}

// Grab all Permissions for Users
func GetUserRoles(u string) map[string][]string {

	//Init Db
	db := GetPostgreDB().Table("casbin_rule")

	// Declare
	var permissions []CasbinRule
	roles := GetRoles()
	groups := GetGroups2()
	rolesT := make(map[string][]string)

	// Find User Groupings
	db.Where("V0 = ?", u).Find(&permissions)

	// Loop Through permissions, creating names slices
	for _, perm := range permissions {
		r := roles[perm.V1]
		g, _ := strconv.Atoi(perm.V2)
		rolesT[r] = append(rolesT[r], groups[g])
	}

	return rolesT
}

// todo
func GetUserRoles2(u string, g int) []string {

	//Init Db
	db := GetPostgreDB().Table("casbin_rule")

	// Declare
	var permissions []CasbinRule
	var userRoles []string
	roles := GetRoles()

	// Find User Groupings
	db.Where("V0 = ? AND V2 = ?", u, g).Find(&permissions)
	for _, role := range permissions {
		userRoles = append(userRoles, roles[role.V1])
	}

	return userRoles
}

func GetDashRoles(u int, g int) []string {

	//Init Db
	db := GetPostgreDB().Table("casbin_rule")

	// Declare
	var permissions []CasbinRule
	var dashRoles []string
	roles := GetRoles()

	// Find Dash Groupings
	db.Where("V1 = ? AND V2 = ?", u, g).Find(&permissions)
	for _, role := range permissions {
		dashRoles = AppendIfMissing(dashRoles, roles[role.V0])
	}

	return dashRoles
}

func GetAllPermissionsGroup() []CasbinRule {
	db := GetPostgreDB().Table("casbin_rule")

	var permissions []CasbinRule

	// Find User Groupings
	db.Where("p_type = ?", "g").Find(&permissions)

	// Find All Permissions From Group

	return permissions
}

func GetAllPermissionsPolicies() []CasbinRule {
	db := GetPostgreDB().Table("casbin_rule")

	var permissions []CasbinRule

	// Find User Groupings
	db.Where("p_type = ?", "p").Find(&permissions)

	// Find All Permissions From Group

	return permissions
}

func AppendIfMissing(slice []string, i string) []string {
	for _, ele := range slice {
		if len(i) <= 0 {
			return slice
		}
		if ele == i {
			return slice
		}

	}
	return append(slice, i)
}
