package services

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lalela-backend/internal/pkg/models"
	"time"
)

func FindUserByEmail(email string) (models.User, error) {
	var collection = OpenCollection("user")
	var foundUser models.User
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&foundUser)
	if err != nil {
		return models.User{}, err
	}
	return foundUser, err
}

func getUserOrganizations(email string) ([]models.Organizations, error) {
	foundUser, err := FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	orgs, err := Enforcer.GetDomainsForUser(foundUser.ID.Hex())
	if err != nil {
		return nil, err
	}

	var collection = OpenCollection("organizations")
	var organizations []models.Organizations

	for _, orgId := range orgs {
		var organization models.Organizations
		err := collection.FindOne(context.Background(), bson.M{"_id": orgId}).Decode(&organization)
		if err != nil {
			return []models.Organizations{}, err
		}
		organizations = append(organizations, organization)
	}
	return organizations, nil
}


func GetUsers() {

}

func CreateUser(email string, password string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	var collection = OpenCollection("user")

	user.ID = primitive.NewObjectID()
	user.Email = email


	token, _ := GenerateToken(email)
	user.ValidationToken = token
	pass := HashPassword(password)
	user.Password = pass

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	insertId, err := collection.InsertOne(ctx, user)
	if err != nil {
		cancel()
		return "", err
	}
	defer cancel()

	return cast.ToString(insertId.InsertedID), err
}


func LoginUserCredentials(email string, password string) (string, error) {
	var collection = OpenCollection("user")

	foundUser, err := FindUserByEmail(email)
	if err != nil {
		fmt.Println("User Not Found")
		return "", err
	}

	passwordIsValid := VerifyPassword(password, foundUser.Password)
	if !passwordIsValid {
		return "", nil
	}

	jwtToken, err := GenerateToken(email)
	if err != nil {
		fmt.Println("Generating token error")
		return "", err
	}

	UpdateToken(collection, email, jwtToken)
	return foundUser.ID.Hex(), nil
}