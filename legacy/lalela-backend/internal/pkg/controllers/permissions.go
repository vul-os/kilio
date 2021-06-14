package controllers

import (
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/services"
	"lalela-backend/internal/pkg/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type PermissionCon struct{}

func PermissionExists(w http.ResponseWriter, r *http.Request) {

	t, done := CheckJson(w, r, &PermissionsAR{})
	if done {
		return
	}

	// Here you cast the interface to the concrete type
	perms := t.(*PermissionsAR)

	CreatePermissionsGroups()

	// Loop Through Struct Map
	for _, permission := range perms.Permissions {

		dashDb, err := services.GetPermissionExists(&permission)
		if err != nil {
			resp, err := utils.GetError(404, "en")
			if err != nil {
				json.NewEncoder(w).Encode(err)
				return
			}
			json.NewEncoder(w).Encode(resp)
			return
		}

		var resp = map[string]interface{}{}

		if dashDb == true {
			if x := services.AddPermission(&permission); x != false {
				resp["code"] = "990"
				resp["status"] = "Dashboard Permission Doesn't Exist; Added"
			} else {
				resp["code"] = "991"
				resp["status"] = "Empty Fields; Not Added"
			}
		} else {
			resp["code"] = "992"
			resp["status"] = "Dashboard Permission Exist; Not Added"
		}

		json.NewEncoder(w).Encode(resp)
	}
}

// Contains tells whether a contains x.
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

func CreatePermissionsGroups() {

	gs := services.GetGroups()
	ds := services.GetDashboards()
	us := services.GetUsersRaw()

	for _, dashboard := range ds {

		t := SliceContainsGroup(gs, strings.ToLower(dashboard.CategoryName))

		if t != 0 {
			utils.PolicyExistsAdd("role::4", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "1")
			utils.PolicyExistsAdd("role::2", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "1")
			utils.PolicyExistsAdd("role::2", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "2")
			utils.PolicyExistsAdd("role::2", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "3")
			utils.PolicyExistsAdd("role::3", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "1")
			utils.PolicyExistsAdd("role::3", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "2")
		}
	}

	for _, user := range us {
		if user.UserGroupId == 1 {
			for _, group := range gs {
				utils.GroupPolicyExistsAdd("user::"+fmt.Sprint(user.ID), "role::2", fmt.Sprint(group.ID))
			}
		} else {
			utils.GroupPolicyExistsAdd("user::"+fmt.Sprint(user.ID), "role::4", fmt.Sprint(user.UserGroupId))
		}
	}
}

func (t *PermissionCon) GetUserPermissions(r *http.Request, args *models.GetUserRolesRequest, reply *models.GetUserRolesResponse) error {
	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	reply.Roles = services.GetUserRoles("user::" + fmt.Sprint(args.UserId))

	return nil

}

func (t *PermissionCon) GetDashboardPermissions(r *http.Request, args *models.GetDashboardPermissionsRequest, reply *models.GetDashboardPermissionsResponse) error {
	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	reply.Permissions = services.GetDashboardPermissions(fmt.Sprint(args.DashId))

	return nil

}

func (t *PermissionCon) GetGroupPermissions(r *http.Request, args *models.GetGroupPermissionsRequest, reply *models.GetGroupPermissionsResponse) error {
	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	services.GetPermissionsGroup("4", "role::4")

	return nil

}

// Todo
func (t *PermissionCon) UpdateUserRole(r *http.Request, args *models.GetGroupPermissionsRequest, reply *models.GetGroupPermissionsResponse) error {
	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	services.GetPermissionsGroup("4", "role::4")

	return nil

}

// Todo
func (t *PermissionCon) UpdateGroupPermissions(r *http.Request, args *models.GetGroupPermissionsRequest, reply *models.GetGroupPermissionsResponse) error {
	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	services.GetPermissionsGroup("4", "role::4")

	return nil

}
