package services

import (
	"cog-analytics-engine-go/internal/pkg/controllers"
	"cog-analytics-engine-go/internal/pkg/models"
	"cog-analytics-engine-go/internal/pkg/utils"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type UserCon struct{}

func (t *UserCon) UsersGet(r *http.Request, args *controllers.UsersGetRequest, reply *controllers.UsersGetResponse) error {

	err := utils.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	sourceUser, err := utils.FindByEmailString(args.Email)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	userData, err := utils.GetUsers(sourceUser)
	reply.Users = userData
	return nil
}
func (t *UserCon) UserGet(r *http.Request, args *controllers.UserGetRequest, reply *controllers.UserGetResponse2) error {

	err := utils.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	userData, err := utils.GetUser(models.User{ID: args.Id})
	reply.User = userData
	reply.Dashboards = utils.GetUserRoles("user::" + fmt.Sprint(args.Id))
	return nil
}
func (t *UserCon) UserIsAdmin(r *http.Request, args *controllers.UserIsAdminRequest, reply *controllers.UserIsAdminResponse) error {

	err := utils.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	userData, err := utils.GetUserRaw(models.User{Email: args.Email})
	userIsAdmin := utils.IsUserAdmin(userData)
	reply.Admin = userIsAdmin
	return nil
}
func (t *UserCon) UserAdd(r *http.Request, args *controllers.UserAddRequest, reply *controllers.UserAddResponse) error {

	err := utils.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	sourceUser, err := utils.FindByEmailString(args.Email)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	for _, user := range args.Users {

		// Check if Source User is allowed to Edit User Group
		if sourceUser.UserGroupId != 1 {
			if !utils.PermissionCanEditUser(fmt.Sprint(sourceUser.ID), fmt.Sprint(user.UserGroupId)) {
				reply.Messages = append(reply.Messages, "User Does Not Have Permission to Add Users")
				return nil
			}
		}

		userDB, err := utils.GetUserExists(&user)
		if err != nil {
			log.Printf("Error: %s", err)
			return nil
		}

		if userDB == true {
			if x := utils.AddUser(&user); x != false {
				reply.Messages = append(reply.Messages, "User Does not Exist, Added "+user.Email)
				tempUser, _ := utils.FindByEmailString(user.Email)
				utils.GroupPolicyExistsAdd("user::"+fmt.Sprint(tempUser.ID), "role::4", fmt.Sprint(tempUser.UserGroupId))
			} else {
				reply.Messages = append(reply.Messages, "Empty User Fields, Not Added "+user.Email)
			}
		} else {
			reply.Messages = append(reply.Messages, "User Exists, Not Added "+user.Email)
		}
	}
	return nil
}
func (t *UserCon) UserUpdate(r *http.Request, args *controllers.UserUpdateRequest, reply *controllers.UserUpdateResponse) error {

	err := utils.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	//sourceUser, err := services.FindByEmailString(args.Email)
	//if err != nil {
	//	log.Printf("Error: %s", err)
	//	return err
	//}

	for _, user := range args.Users {

		//// Check if Source User is allowed to Edit User Group
		//if sourceUser.UserGroupId != 1 {
		//	reply.Messages = append(reply.Messages, "User Does Not Have Permission to Delete")
		//	return nil
		//}

		userDb, err := utils.GetUserExists(&user)
		if err != nil {
			log.Printf("Error: %s", err)
			return nil
		}

		if userDb == true {
			reply.Messages = append(reply.Messages, "User Does Not Exist "+user.Email)
		} else {
			if user.ID != 0 {
				pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
				if err != nil {

				}
				user.Password = string(pass)
				utils.UpdateUser(&user)
				reply.Messages = append(reply.Messages, "User Exists, Updated "+user.Email)
			} else {
				reply.Messages = append(reply.Messages, "No User ID Given "+user.Email)
			}
		}

	}
	return nil
}
func (t *UserCon) UserUpdateSingle(r *http.Request, args *controllers.UserUpdateRequestSingle, reply *controllers.UserUpdateResponse) error {

	err := utils.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	sourceUser, err := utils.FindByEmailString(args.Email)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	// Check if Source User is allowed to Edit User Group
	if sourceUser.UserGroupId != 1 {
		reply.Messages = append(reply.Messages, "User Does Not Have Permission to Delete")
		return nil
	}

	userDb, err := utils.GetUserExists(&args.User)
	if err != nil {
		log.Printf("Error: %s", err)
		return nil
	}

	if userDb == true {
		reply.Messages = append(reply.Messages, "User Does Not Exist "+args.User.Email)
	} else {
		if args.User.ID != 0 {
			user, _ := utils.FindByEmail(&args.User)

			user.FirstName = args.User.FirstName
			user.LastName = args.User.LastName
			utils.UpdateUser(user)
			reply.Messages = append(reply.Messages, "User Exists, Updated "+args.User.Email)
		} else {
			reply.Messages = append(reply.Messages, "No User ID Given "+args.User.Email)
		}
	}

	return nil
}
func (t *UserCon) UserDelete(r *http.Request, args *controllers.UserDeleteRequest, reply *controllers.UserDeleteResponse) error {

	err := utils.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	sourceUser, err := utils.FindByEmailString(args.Email)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	// Check if Source User is allowed to Edit User Group
	if sourceUser.UserGroupId != 1 {
		reply.Messages = append(reply.Messages, "User Does Not Have Permission to Delete")
		return nil
	}

	userDb, err := utils.GetUserExists(&args.User)
	if err != nil {
		log.Printf("Error: %s", err)
		return nil
	}

	if userDb == true {
		reply.Messages = append(reply.Messages, "User Does Not Exist "+args.User.Email)
	} else {
		if args.User.ID != 0 {
			user, _ := utils.FindByEmail(&args.User)
			utils.DeleteUser(&args.User)
			reply.Messages = append(reply.Messages, "User Exists, Deleted "+args.User.Email)
			utils.GroupPolicyExistsRemove("user::"+fmt.Sprint(user.ID), "role::4", fmt.Sprint(user.UserGroupId))
		} else {
			reply.Messages = append(reply.Messages, "No User ID Given "+args.User.Email)
		}
	}

	return nil
}
func (t *UserCon) UserRolesUpdate(r *http.Request, args *controllers.UserRoleUpdateRequest, reply *controllers.UserRoleUpdateResponse) error {

	err := utils.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	sourceUser, err := utils.FindByEmailString(args.Email)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	// Check if Source User is allowed to Edit User Group
	if sourceUser.UserGroupId != 1 {
		reply.Messages = append(reply.Messages, "User Does Not Have Permission to Delete")
		return nil
	}

	userDb, err := utils.GetUserExists(&args.User)
	if err != nil {
		log.Printf("Error: %s", err)

		return nil
	}

	if userDb == true {
		reply.Messages = append(reply.Messages, "User Does Not Exist "+args.User.Email)
	} else {
		if args.User.ID != 0 {
			utils.UserPermission(args.Roles, args.GroupId, args.User.ID)
			reply.Messages = append(reply.Messages, "User Exists, Deleted "+args.User.Email)
			log.Printf("Error: %s", args.Roles)
			//utils.GroupPolicyExistsRemove("user::"+fmt.Sprint(args.User.ID), "role::4", fmt.Sprint(args.User.UserGroupId))
		} else {
			reply.Messages = append(reply.Messages, "No User ID Given "+args.User.Email)
		}
	}

	return nil
}
