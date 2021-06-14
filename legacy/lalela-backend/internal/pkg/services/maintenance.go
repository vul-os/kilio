package services

import (
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/utils"
	"fmt"
	"regexp"
	"strconv"
)

func CleanRbac() ([]string, error) {
	var log []string
	permissionsGroup := GetAllPermissionsGroup()
	permissionsPolicies := GetAllPermissionsPolicies()

	logUser := CleanRbacUsers(permissionsGroup)
	logDash := CleanRbacDashboards(permissionsPolicies)
	log = append(log, logUser...)
	log = append(log, logDash...)

	return log, nil
}

func CleanRbacUsers(permissions []CasbinRule) []string {
	var log []string
	var users []models.User
	// Init DB
	db := utils.GetDB()
	db.Find(&users)

	var userArray []int

	for _, r := range users {
		userArray = append(userArray, int(r.ID))
	}

	for _, permission := range permissions {
		re := regexp.MustCompile("[0-9]+")
		userId := re.FindAllString(permission.V0, -1)
		userIdFinal, _ := strconv.Atoi(userId[0])
		if contains(userArray, userIdFinal) {
			fmt.Printf("Contain: ", userIdFinal)
		} else {
			utils.GroupPolicyExistsRemove("user::"+fmt.Sprint(userIdFinal), permission.V1, fmt.Sprint(permission.V2))
			log = append(log, "user::"+fmt.Sprint(userIdFinal)+" Removed")
		}
	}

	return log
}
func CleanRbacDashboards(permissions []CasbinRule) []string {
	var log []string
	var dashboards []models.Dashboard

	// Init DB
	db := utils.GetDB()
	db.Find(&dashboards)

	var dashboardArray []int

	for _, r := range dashboards {
		dashboardArray = append(dashboardArray, int(r.Id))
	}

	for _, permission := range permissions {
		re := regexp.MustCompile("[0-9]+")
		dashId := re.FindAllString(permission.V0, -1)
		dashIdFinal, _ := strconv.Atoi(dashId[0])
		if contains(dashboardArray, dashIdFinal) {
			fmt.Printf("Contain: ", dashIdFinal)
		} else {
			utils.PolicyExistsRemove(permission.V0, fmt.Sprint(permission.V1), fmt.Sprint(dashIdFinal), permission.V3)
			log = append(log, "dash::"+fmt.Sprint(dashIdFinal)+" Removed")
		}
	}

	return log
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
