package controllers

import (

	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/utils"
	"net/http"
	"time"
)

type UserCon struct{}

type UserRequest struct {
	Email string `json:"name"`
	Password string `json:"password"`
}

type UserResponse struct {
	Response string `json:"response"`
	Id string `json:"id"`
}

func (t *UserCon) Login(r *http.Request, args *UserRequest,	reply *UserResponse) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100 * time.Second)
	var foundUser models.User

	var userCollection = utils.OpenCollection(utils.Client, "user")

	err := userCollection.FindOne(ctx, bson.M{"email": args.Email}).Decode(&foundUser)
	if err != nil {
		fmt.Println("No user found")
	}
	defer cancel()

	passwordIsValid := utils.VerifyPassword(args.Password, *foundUser.Password)
	if !passwordIsValid {
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
