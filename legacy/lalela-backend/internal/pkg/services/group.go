package services

import (
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/utils"
	"fmt"
	"strconv"
)

func GetAllGroups() ([]models.Group, error) {

	db := utils.GetDB()

	var groups []models.Group

	db.Find(&groups)

	return groups, nil
}

func GetGroupsMembers(groups []models.Group) (map[string]string, error) {

	db := utils.GetDB().Table("casbin_rule")
	groupsRaw := make(map[string]string)
	for _, groups := range groups {
		var count int64
		var groupDB []CasbinRule

		db.Where("p_type = 'g' AND V2 = ?", fmt.Sprint(groups.ID)).Find(&groupDB).Count(&count)
		s := strconv.FormatInt(count, 10)
		groupsRaw[groups.Name] = s
	}

	return groupsRaw, nil
}

func GetGroupUsers(id uint) ([]models.UserWithRoles, error) {

	db := utils.GetDB()

	var users []models.User
	var usersFinal []models.UserWithRoles

	db.Where("user_group_id = ?", id).Find(&users)

	for _, user := range users {
		userFinal := models.UserWithRoles{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Avatar:    user.Avatar,
			Roles:     GetUserRoles2("user::"+fmt.Sprint(user.ID), int(id)),
		}
		usersFinal = append(usersFinal, userFinal)
	}

	return usersFinal, nil
}

func GetGroupDashboards(id uint) ([]models.DashboardWithRoles, error) {

	db := utils.GetDB()

	var group models.Group

	group.ID = id

	db.Find(&group)

	var dashboards []models.Dashboard
	var dashsFinal []models.DashboardWithRoles

	db.Where("category_name = ?", group.Name).Find(&dashboards)

	for _, dash := range dashboards {
		dashFinal := models.DashboardWithRoles{
			Id:            dash.Id,
			DashboardUrl:  dash.DashboardUrl,
			DashboardName: dash.DashboardName,
			CategoryName:  dash.CategoryName,
			Roles:         GetDashRoles(int(group.ID), int(dash.Id)),
		}
		dashsFinal = append(dashsFinal, dashFinal)
	}

	return dashsFinal, nil
}

func GetGroupPermissionsG(id uint) ([]models.CasbinRule, error) {

	db := utils.GetDB().Table("casbin_rule")

	var permissions []models.CasbinRule
	//var permissionsFinal []string

	db.Where("p_type = 'g' AND V2 = ?", id).Find(&permissions)

	//for _, permission := range permissions {
	//
	//}

	return permissions, nil
}

func GetGroupPermissionsP(id uint) ([]models.CasbinRule, error) {

	db := utils.GetDB().Table("casbin_rule")

	var permissions []models.CasbinRule

	db.Where("p_type = 'p' AND V1 = ?", id).Find(&permissions)

	return permissions, nil
}

func AddGroup(name string) (models.Group, error) {

	db := utils.GetDB()
	group := models.Group{
		Name: name,
	}
	db.Create(&group)

	return group, nil
}

func GetGroupCount() int64 {
	var count int64
	db := utils.GetDB()
	var groups []models.Group
	db.Find(&groups).Count(&count)
	return count
}

func GetGroupExists(group *models.Group) (bool, error) {
	db := utils.GetDB()
	if group.ID != 0 {
		if err := db.Where(&models.Group{ID: group.ID}).Find(&group).Error; err != nil {
			return true, nil
		} else {
			return false, nil
		}
	} else {
		return true, nil
	}
}

func DeleteGroup(group *models.Group) {
	db := utils.GetDB().Table("groups")
	db.Delete(&group)
}
