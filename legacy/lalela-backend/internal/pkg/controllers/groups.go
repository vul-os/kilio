package controllers

import (
	"lalela-backend/internal/pkg/middleware"
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/services"
	"log"
	"net/http"
)

type GroupsCon struct{}

func (t *GroupsCon) GroupsGet(r *http.Request, args *models.GroupsGetRequest, reply *models.GroupsGetResponse) error {

	//var Groups map[string]string

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	//sourceUser, err := services.FindByEmailString(args.Email)
	//if err != nil {
	//	log.Printf("Error: %s", err)
	//	return err
	//}

	groupsData, err := services.GetAllGroups()

	//Groups, err = services.GetGroupsMembers(groupsData)
	reply.Groups = groupsData
	return nil
}

func (t *GroupsCon) GroupGet(r *http.Request, args *models.GroupGetRequest, reply *models.GroupGetResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	// Todo Add Security

	//sourceUser, err := services.FindByEmailString(args.Email)
	//if err != nil {
	//	log.Printf("Error: %s", err)
	//	return err
	//}

	userData, err := services.GetGroupUsers(args.Id)
	dashData, err := services.GetGroupDashboards(args.Id)
	groupPermissionsDataG, err := services.GetGroupPermissionsG(args.Id)
	groupPermissionsDataP, err := services.GetGroupPermissionsP(args.Id)

	reply.Users = userData
	reply.Dashboards = dashData
	reply.GroupPermissionsG = groupPermissionsDataG
	reply.GroupPermissionsP = groupPermissionsDataP

	return nil
}

func (t *GroupsCon) GroupAdd(r *http.Request, args *models.GroupAddRequest, reply *models.GroupAddResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	// Todo Add Security

	groupData, err := services.AddGroup(args.Name)

	reply.Messages = append(reply.Messages, "Group Added "+groupData.Name)

	return nil
}

func (t *GroupsCon) GroupsGetRaw(r *http.Request, args *models.GroupsGetRawRequest, reply *models.GroupsGetRawResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	//sourceUser, err := services.FindByEmailString(args.Email)
	//if err != nil {
	//	log.Printf("Error: %s", err)
	//	return err
	//}

	groupsData, err := services.GetAllGroups()

	reply.Groups = groupsData
	return nil
}

func (t *GroupsCon) GroupDelete(r *http.Request, args *models.GroupGetRequest, reply *models.GroupDeleteResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	sourceUser, err := services.FindByEmailString(args.Email)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	if sourceUser.UserGroupId != 1 {
		reply.Messages = append(reply.Messages, "User Does Not Have Permission to Delete")
		return nil
	}

	groupId := models.Group{ID: args.Id}
	groupDb, err := services.GetGroupExists(&groupId)
	if err != nil {
		log.Print(middleware.NewError(err))
		return nil
	}

	if groupDb == true {
		reply.Messages = append(reply.Messages, "Dashboard Does Not Exist "+string(args.Id))
	} else {
		// TODO: Permission Removal
		if groupId.ID != 0 {
			//services.DeleteDashboard(&dashboard)
			services.DeleteGroup(&groupId)
			reply.Messages = append(reply.Messages, "Dashboard Exists, Deleted "+groupId.Name)

		} else {
			reply.Messages = append(reply.Messages, "No Dashboard ID Given "+groupId.Name)
		}
	}

	return nil
}
