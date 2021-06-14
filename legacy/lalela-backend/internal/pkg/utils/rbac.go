package utils

import (
	"lalela-backend/internal/pkg/middleware"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"log"
)

// Init RBCA Enforcer
func InitRBCA() *casbin.Enforcer {

	a, _ := gormadapter.NewAdapter("postgres", dbURIString, true)

	e, err := casbin.NewEnforcer("./internal/config/model.conf", a)
	if err != nil {
		log.Print(middleware.NewError(err))
	}

	e.LoadPolicy()

	return e
}

// Can User View Dash
func PermissionCanViewDash(e *casbin.Enforcer, UserID string, GroupID string, DashID string) bool {

	// Check Permissions
	var d bool

	// todo : Change Root Account Detection
	if UserID == "2" {
		d = true
	} else {
		d, _ = e.Enforce("user::"+UserID, GroupID, DashID, "1")
	}

	//fmt.Println("user::"+UserID,GroupID, DashID, "1")
	return d
}

// Can User Edit Dash
func PermissionCanEditDash(UserID string, GroupID string, DashID string) bool {

	// Check Permissions
	var d bool

	e := InitRBCA()

	// todo : Change Root Account Detection
	if UserID == "2" {
		d = true
	} else {
		d, _ = e.Enforce("user::"+UserID, GroupID, DashID, "2")
	}

	//fmt.Println("user::"+UserID,GroupID, DashID, "1")
	return d
}

// Can User Delete Dash
func PermissionCanDeleteDash(e *casbin.Enforcer, UserID string, GroupID string, DashID string) bool {

	// Check Permissions
	var d bool

	// todo : Change Root Account Detection
	if UserID == "2" {
		d = true
	} else {
		d, _ = e.Enforce("user::"+UserID, GroupID, DashID, "3")
	}

	//fmt.Println("user::"+UserID,GroupID, DashID, "1")
	return d
}

// Can User View Group Members
func PermissionCanViewUser(e *casbin.Enforcer, UserID string, GroupID string) bool {

	// Check Permissions
	var d bool

	// todo : Change Root Account Detection
	if UserID == "2" {
		d = true
	} else {
		d, _ = e.Enforce("user::"+UserID, GroupID, "group::"+GroupID, "1")
	}

	return d
}

// Can User Edit Group Members
func PermissionCanEditUser(UserID string, GroupID string) bool {
	e := InitRBCA()
	// Check Permissions
	var d bool

	// todo : Change Root Account Detection
	if UserID == "2" {
		d = true
	} else {
		d, _ = e.Enforce("user::"+UserID, GroupID, "group::"+GroupID, "2")
	}

	return d
}

// Can User Delete Group Members
func PermissionCanDeleteUser(e *casbin.Enforcer, UserID string, GroupID string) bool {

	// Check Permissions
	var d bool

	// todo : Change Root Account Detection
	if UserID == "2" {
		d = true
	} else {
		d, _ = e.Enforce("user::"+UserID, GroupID, "group::"+GroupID, "3")
	}

	return d
}

// Check if Policy exits
func PolicyExists(sub string, dom string, obj string, act string) bool {

	// ex. utils.GroupPolicyExists(e,"role::2","2", "2", "2")
	// sub can be user::id (ex. user::1), role::id (ex. role::1)
	// dom can only be group_id
	// obj can be dashboard_id (ex. 2), group::id (ex. group::1)
	// act can only be actions_id
	e := InitRBCA()
	d, _ := e.Enforce(sub, dom, obj, act)
	return d
}

// Check if Policy Exists, Add if not
func PolicyExistsAdd(sub string, dom string, obj string, act string) bool {
	e := InitRBCA()
	if PolicyExists(sub, dom, obj, act) {
		return true
	} else {
		e.AddPolicy(sub, dom, obj, act)
		return false
	}
}

// Check if Policy Exists, Add if not
func PolicyExistsRemove(sub string, dom string, obj string, act string) bool {
	e := InitRBCA()
	if PolicyExists(sub, dom, obj, act) {
		e.RemovePolicy(sub, dom, obj, act)
		return true
	} else {
		return false
	}
}

// Check if Grouping Policy Exists
func GroupPolicyExists(sub string, obj string, dom string) bool {

	// ex. utils.GroupPolicyExists(e,"user::19","role::2", "2")
	// sub can only be user::id (ex. user::1)
	// obj can only be role::id (ex. role::2)
	// dom can only be group_id
	e := InitRBCA()
	d := e.HasGroupingPolicy(sub, obj, dom)
	return d

}

// Check if Grouping Policy Exists, Add if not
func GroupPolicyExistsAdd(sub string, obj string, dom string) bool {
	e := InitRBCA()
	if GroupPolicyExists(sub, obj, dom) {
		return true
	} else {
		e.AddGroupingPolicy(sub, obj, dom)
		return false
	}
}

// Check if Grouping Policy Exists, Remove if does
func GroupPolicyExistsRemove(sub string, obj string, dom string) bool {
	e := InitRBCA()
	if GroupPolicyExists(sub, obj, dom) {
		e.RemoveGroupingPolicy(sub, obj, dom)
		return true
	} else {
		return false
	}
}

//// KANBAN

// Check if User can assign to card
func PermissionsKanbanUserCanAssign() {

}

// Check if User is allowed to move card
func PermissionsKanbanUserCanMoveCard() {

}

// Check if Useer has access to card, if not Add
func PermissionsKanbanAddUserToCard() {

}

// Check if Useer has access to card, if so Remove
func PermissionsKanbanRemoveUserFromCard() {

}
