package controllers

import (
	"context"
	"fmt"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/utils"
	"net/http"
	"time"
)

type UserCon struct{}

type UserRequest struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Response string `json:"response"`
	Id string `json:"id"`
}

// https://github.com/Joojo7/user-athentication-golang/blob/master/controllers/userController.go
func (t *UserCon) Login(r *http.Request, args *UserRequest,	reply *UserResponse) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
	var foundUser models.User

	var userCollection = utils.OpenCollection(utils.MongoClient, "user")

	err := userCollection.FindOne(ctx, bson.M{"email": args.Email}).Decode(&foundUser)
	if err != nil {
		fmt.Println("No user found")
		cancel()
		return err
	}
	defer cancel()

	passwordIsValid := utils.VerifyPassword(args.Password, *foundUser.Password)
	if !passwordIsValid {
		cancel()
		return nil
	}
	defer cancel()

	jwtToken, err := utils.GenerateToken(args.Email)
	if err != nil {
		fmt.Println("Generating token error")
	}

	utils.UpdateToken(userCollection, args.Email, jwtToken)
	reply.Response = "ok"
	reply.Id = foundUser.ID.String()
	return nil
}


func (t *UserCon) RegisterUser(r *http.Request, args *UserRequest,	reply *UserResponse) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user models.User
	var userCollection = utils.OpenCollection(utils.MongoClient, "user")

	user.Email = &args.Email
	user.ID = primitive.NewObjectID()
	token, _ := utils.GenerateToken(args.Email)
	password := utils.HashPassword(args.Password)
	user.ValidationToken = &token
	user.Password = &password

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	fmt.Println(user)
	insertId, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		cancel()
		return err
	}
	reply.Id = cast.ToString(insertId)
	defer cancel()
	return nil
}