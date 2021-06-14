package controllers

import (
	"lalela-backend/internal/pkg/middleware"
	models "lalela-backend/internal/pkg/models"
	services "lalela-backend/internal/pkg/services"
	"lalela-backend/internal/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type DashCon struct{}

func (t *DashCon) DashboardData(r *http.Request, args *models.DashboardDataRequest, reply *models.DashboardDataResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	u := &models.User{
		Email: args.Email,
	}
	userDb, err := services.FindByEmail(u)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	dashDb, err := services.GetDashboardData(userDb)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}
	dashDbFinal := services.SortDashboards(dashDb)
	reply.DashboardData = dashDbFinal

	return nil
}
func (t *DashCon) DashboardsGet(r *http.Request, args *models.DashboardDataRequest, reply *models.DashboardsGetResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	dashDb := services.GetDashboards()
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	reply.Dashboards = dashDb

	return nil
}
func (t *DashCon) DashboardGet(r *http.Request, args *models.DashboardSingleDataRequest, reply *models.DashboardSingleDataResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	dashDb := services.GetDashboard(models.Dashboard{Id: args.Id})
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	reply.DashboardData = dashDb

	return nil
}
func (t *DashCon) DashboardAdd(r *http.Request, args *models.DashboardAddRequest, reply *models.DashboardAddResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	sourceUser, err := services.FindByEmailString(args.Email)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	for _, dashboard := range args.Dashboards {

		// Check if Source User is allowed to Edit User Group
		if sourceUser.UserGroupId != 1 {
			reply.Messages = append(reply.Messages, "User Does Not Have Permission to Add")
			return nil
		}

		dashDb, err := services.GetDashboardExists(&dashboard)
		if err != nil {
			log.Print(middleware.NewError(err))
			return nil
		}

		if dashDb == true {
			if x := services.AddDashBoard(&dashboard); x != false {
				reply.Messages = append(reply.Messages, "Dashboard Does not Exist, Added "+dashboard.DashboardUrl)

				gs := services.GetGroups()
				t := SliceContainsGroup(gs, strings.ToLower(dashboard.CategoryName))

				utils.PolicyExistsAdd("role::4", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "1")
			} else {
				reply.Messages = append(reply.Messages, "Empty Dashboard Fields, Not Added "+dashboard.DashboardUrl)
			}
		} else {
			reply.Messages = append(reply.Messages, "Dashboard Exists, Not Added "+dashboard.DashboardUrl)
		}
	}
	return nil
}
func (t *DashCon) DashboardDelete(r *http.Request, args *models.DashboardDeleteRequest, reply *models.DashboardDeleteResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	sourceUser, err := services.FindByEmailString(args.Email)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	for _, dashboard := range args.Dashboards {

		// Check if Source User is allowed to Edit User Group
		if sourceUser.UserGroupId != 1 {
			reply.Messages = append(reply.Messages, "User Does Not Have Permission to Delete")
			return nil
		}

		dashDb, err := services.GetDashboardExists(&dashboard)
		if err != nil {
			log.Print(middleware.NewError(err))
			return nil
		}

		if dashDb == true {
			reply.Messages = append(reply.Messages, "Dashboard Does Not Exist "+dashboard.DashboardUrl)
		} else {
			if dashboard.Id != 0 {
				services.DeleteDashboard(&dashboard)

				reply.Messages = append(reply.Messages, "Dashboard Exists, Deleted "+dashboard.DashboardUrl)

				gs := services.GetGroups()
				t := SliceContainsGroup(gs, strings.ToLower(dashboard.CategoryName))

				utils.PolicyExistsRemove("role::4", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "1")
				utils.PolicyExistsRemove("role::3", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "2")
				utils.PolicyExistsRemove("role::2", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "2")
			} else {
				reply.Messages = append(reply.Messages, "No Dashboard ID Given "+dashboard.DashboardUrl)
			}
		}

	}
	return nil
}
func (t *DashCon) DashboardDeleteSingle(r *http.Request, args *models.DashboardDeleteSingleRequest, reply *models.DashboardDeleteResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Print(middleware.NewError(err))
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

	dashboard := services.GetDashboard(models.Dashboard{Id: args.Dashboard})
	dashDb, err := services.GetDashboardExists(&dashboard)
	if err != nil {
		log.Print(middleware.NewError(err))
		return nil
	}

	if dashDb == true {
		reply.Messages = append(reply.Messages, "Dashboard Does Not Exist "+string(args.Dashboard))
	} else {

		if dashboard.Id != 0 {
			services.DeleteDashboard(&dashboard)

			reply.Messages = append(reply.Messages, "Dashboard Exists, Deleted "+dashboard.DashboardUrl)

			gs := services.GetGroups()
			t := SliceContainsGroup(gs, strings.ToLower(dashboard.CategoryName))

			utils.PolicyExistsRemove("role::4", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "1")
			utils.PolicyExistsRemove("role::3", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "2")
			utils.PolicyExistsRemove("role::2", fmt.Sprint(t), fmt.Sprint(dashboard.Id), "2")
		} else {
			reply.Messages = append(reply.Messages, "No Dashboard ID Given "+dashboard.DashboardUrl)
		}
	}

	return nil
}
func (t *DashCon) DashboardUpdate(r *http.Request, args *models.DashboardDeleteRequest, reply *models.DashboardDeleteResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	sourceUser, err := services.FindByEmailString(args.Email)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	for _, dashboard := range args.Dashboards {

		// Check if Source User is allowed to Edit User Group
		if sourceUser.UserGroupId != 1 {
			reply.Messages = append(reply.Messages, "User Does Not Have Permission to Delete")
			return nil
		}

		dashDb, err := services.GetDashboardExists(&dashboard)
		if err != nil {
			log.Print(middleware.NewError(err))
			return nil
		}

		if dashDb == true {
			reply.Messages = append(reply.Messages, "Dashboard Does Not Exist "+dashboard.DashboardUrl)
		} else {
			if dashboard.Id != 0 {
				services.UpdateDashboard(&dashboard)
				reply.Messages = append(reply.Messages, "Dashboard Exists, Updated "+dashboard.DashboardUrl)
			} else {
				reply.Messages = append(reply.Messages, "No Dashboard ID Given "+dashboard.DashboardUrl)
			}
		}

	}
	return nil
}
func (t *DashCon) DashboardUpdateSingle(r *http.Request, args *models.DashboardUpdateSingleRequest, reply *models.DashboardDeleteResponse) error {

	err := services.JwtVerify(r)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	sourceUser, err := services.FindByEmailString(args.Email)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	// Check if Source User is allowed to Edit User Group
	if sourceUser.UserGroupId != 1 {
		reply.Messages = append(reply.Messages, "User Does Not Have Permission to Delete")
		return nil
	}
	dashOgUrl := services.GetDashboard(models.Dashboard{Id: args.Dashboard.Id})
	dashDb, err := services.GetDashboardExists(&dashOgUrl)
	if err != nil {
		log.Print(middleware.NewError(err))
		return nil
	}

	if dashDb == true {
		reply.Messages = append(reply.Messages, "Dashboard Does Not Exist "+args.Dashboard.DashboardUrl)
	} else {
		if args.Dashboard.Id != 0 {
			services.UpdateDashboard(&args.Dashboard)
			reply.Messages = append(reply.Messages, "Dashboard Exists, Updated "+args.Dashboard.DashboardUrl)
		} else {
			reply.Messages = append(reply.Messages, "No Dashboard ID Given "+args.Dashboard.DashboardUrl)
		}
	}

	return nil
}
func (t *DashCon) DashboardRolesUpdate(r *http.Request, args *models.DashboardRolesRequest, reply *models.DashboardRolesResponse) error {
	err := services.JwtVerify(r)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	sourceUser, err := services.FindByEmailString(args.Email)
	if err != nil {
		log.Printf("Error: %s", err)
		return err
	}

	// Check if Source User is allowed to Edit User Group
	if sourceUser.UserGroupId != 1 {
		reply.Messages = append(reply.Messages, "User Does Not Have Permission to Delete")
		return nil
	}

	dashDb, err := services.GetDashboardExists(&args.Dash)
	if err != nil {
		log.Print(middleware.NewError(err))
		return nil
	}

	if dashDb == true {
		reply.Messages = append(reply.Messages, "Dashboard Does Not Exist "+args.Dash.DashboardUrl)
	} else {
		if args.Dash.Id != 0 {
			services.DashboardPermission(args.Roles, args.GroupId, args.Dash.Id)
			reply.Messages = append(reply.Messages, "Dashboard Exists, Updated "+args.Dash.DashboardUrl)
		} else {
			reply.Messages = append(reply.Messages, "No Dashboard ID Given "+args.Dash.DashboardUrl)
		}
	}

	return nil
}
