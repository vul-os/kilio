package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lalela-backend/internal/pkg/auth"
	"lalela-backend/internal/pkg/auth/utils"
	"lalela-backend/internal/pkg/database"
	organizations "lalela-backend/internal/pkg/organizations"
	"lalela-backend/internal/pkg/users"

	"time"
)

// todo: this should not be part of the backend?


var RootAdminEmail = "imranparuk@live.com"
var orgSP1AdminEmail = "notimran@live.com"
var orgSP1NormalEmail = "skaapie@live.com"

var RootAdminId string

var passwordForAll = "lalela"

var organizationz = [][]string {
	{"Spar Branch 1", ""},
	{"Spar Branch 2", ""},
}

func AdminPolicies() {
	actions := []string{"get", "create"}
	objects := []string{"form", "forms", "user", "users", "submission", "submissions"}
	for _, obj := range objects {
		for _, act := range actions {
			_, err := auth.Enforcer.AddPolicy("admin", "global", obj, act)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	_, _ = auth.Enforcer.AddGroupingPolicy(RootAdminId, "admin", "global")
}


func SeedDB() {
	addOrgs()
	addSuperAdminUser()
	addOrgAdminUser()
	addOrgUser()
	AdminPolicies()
}

func addOrgs() {
	for _, org := range organizationz {
		org[1] = addOrg(org[0])
	}
}

func addOrg(orgName string) string {
	var model organizations.Organizations
	var collection = database.OpenCollection("organizations")

	model.ID = primitive.NewObjectID()
	model.Name = orgName
	model.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	model.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	fmt.Println(model)
	result, err := collection.InsertOne(context.Background(), model)
	if err != nil {
		fmt.Println(err)
	}
	return result.InsertedID.(primitive.ObjectID).Hex()
}

func addSuperAdminUser() {
	var user users.User
	var collection = database.OpenCollection("user")

	user.ID = primitive.NewObjectID()
	user.Email = RootAdminEmail

	token, _ := auth.GenerateToken(RootAdminEmail)
	user.ValidationToken = token
	password := utils.HashPassword(passwordForAll)
	user.Password = password

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println(err)
	}
	userId := result.InsertedID.(primitive.ObjectID).Hex()
	_ = fmt.Sprintf(`UserId: %s`, userId)
	RootAdminId = userId

	for _, org := range organizationz {
		//alice, admin, domain1
		added, err := auth.Enforcer.AddGroupingPolicy(userId, "admin", org[1])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(added)
	}
}


func addOrgAdminUser() {
	var user users.User
	var collection = database.OpenCollection("user")

	user.ID = primitive.NewObjectID()
	user.Email = orgSP1AdminEmail

	token, _ := auth.GenerateToken(orgSP1AdminEmail)
	user.ValidationToken = token
	password := utils.HashPassword(passwordForAll)
	user.Password = password

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println(err)
	}
	userId := result.InsertedID.(primitive.ObjectID).Hex()
	_ = fmt.Sprintf(`UserId: %s`, userId)

	//alice, admin, domain1
	added, err := auth.Enforcer.AddGroupingPolicy(userId, "admin", organizationz[0][1])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(added)
}

func addOrgUser()  {
	var user users.User
	var collection = database.OpenCollection("user")

	user.ID = primitive.NewObjectID()
	user.Email = orgSP1NormalEmail

	token, _ := auth.GenerateToken(orgSP1NormalEmail)
	user.ValidationToken = token
	password := utils.HashPassword(passwordForAll)
	user.Password = password

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	result, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println(err)
	}
	userId := result.InsertedID.(primitive.ObjectID).Hex()
	_ = fmt.Sprintf(`UserId: %s`, userId)

	//alice, admin, domain1
	_, err = auth.Enforcer.AddGroupingPolicy(userId, "user", organizationz[0][1])
	if err != nil {
		fmt.Println(err)
	}
}

func addForm() {

}

func addSubmission() {

}