package store

import (
	"context"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lalela-backend/internal/pkg/auth"
	"lalela-backend/internal/pkg/database"
	"lalela-backend/internal/pkg/users"
	"time"
)

func CreateOne(email string, password string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user users.User
	var collection = database.OpenCollection("user")

	user.ID = primitive.NewObjectID()
	user.Email = email


	token, _ := auth.GenerateToken(email)
	user.ValidationToken = token
	pass := auth.HashPassword(password)
	user.Password = pass
	user.OrganizationId = primitive.NilObjectID

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

func FindUserByEmail(email string) (users.User, error) {
	var collection = database.OpenCollection("user")
	var foundUser users.User
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&foundUser)
	if err != nil {
		return users.User{}, err
	}
	return foundUser, err
}
