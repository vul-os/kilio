package dbSeed

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/services"
	"time"
)

// todo: this should not be part of the backend?


var RootAdminEmail = "imranparuk@live.com"
var orgSP1AdminEmail = "notimran@live.com"
var orgSP1NormalEmail = "skaapie@live.com"

var RootAdminId string

var passwordForAll = "lalela"

var organizations = [][]string {
	{"Spar Branch 1", ""},
	{"Spar Branch 2", ""},
}

func AdminPolicies() {
	actions := []string{"get", "create"}
	objects := []string{"form", "forms", "user", "users", "submission", "submissions"}
	for _, obj := range objects {
		for _, act := range actions {
			_, err := services.Enforcer.AddPolicy("admin", "global", obj, act)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	_, _ = services.Enforcer.AddGroupingPolicy(RootAdminId, "admin", "global")
}


func SeedDB() {
	addOrgs()
	addSuperAdminUser()
	addOrgAdminUser()
	addOrgUser()
	AdminPolicies()
}

func addOrgs() {
	for _, org := range organizations {
		org[1] = addOrg(org[0])
	}
}

func addOrg(orgName string) string {
	var model models.Organizations
	var collection = services.OpenCollection("organizations")

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
	var user models.User
	var collection = services.OpenCollection("user")

	user.ID = primitive.NewObjectID()
	user.Email = RootAdminEmail

	token, _ := services.GenerateToken(RootAdminEmail)
	user.ValidationToken = token
	password := services.HashPassword(passwordForAll)
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

	for _, org := range organizations {
		//alice, admin, domain1
		added, err := services.Enforcer.AddGroupingPolicy(userId, "admin", org[1])
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(added)
	}
}


func addOrgAdminUser() {
	var user models.User
	var collection = services.OpenCollection("user")

	user.ID = primitive.NewObjectID()
	user.Email = orgSP1AdminEmail

	token, _ := services.GenerateToken(orgSP1AdminEmail)
	user.ValidationToken = token
	password := services.HashPassword(passwordForAll)
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
	added, err := services.Enforcer.AddGroupingPolicy(userId, "admin", organizations[0][1])
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(added)
}

func addOrgUser()  {
	var user models.User
	var collection = services.OpenCollection("user")

	user.ID = primitive.NewObjectID()
	user.Email = orgSP1NormalEmail

	token, _ := services.GenerateToken(orgSP1NormalEmail)
	user.ValidationToken = token
	password := services.HashPassword(passwordForAll)
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
	_, err = services.Enforcer.AddGroupingPolicy(userId, "user", organizations[0][1])
	if err != nil {
		fmt.Println(err)
	}
}

func addForm() {

}

func addSubmission() {

}