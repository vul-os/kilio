package controllers

import (
	"lalela-backend/internal/pkg/middleware"
	models "lalela-backend/internal/pkg/models"
	services "lalela-backend/internal/pkg/services"
	"lalela-backend/internal/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthCon struct{}

func (t *AuthCon) Login(r *http.Request, args *models.AuthLoginRequest, reply *models.AuthLoginResponse) error {
	fmt.Println(args.Email, args.Password)

	//blocked := services.IsLoginBlocked()
	//if blocked != nil {
	//	log.Print(middleware.NewError(blocked))
	//	return blocked
	//}

	u := &models.User{
		Email: args.Email,
	}
	userDb, err := services.FindByEmail(u)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	errPw := services.CheckPassword(userDb, args.Password)
	if errPw != nil {
		log.Print(middleware.NewError(errPw))
		return errPw
	}

	token, err := services.GetToken(userDb)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}
	services.SetAvatar(userDb, r)
	reply.Jwt = token
	reply.User = *userDb
	EventLog(fmt.Sprint(userDb.ID), "User Logged In", "Auth", "Login")

	return nil
}

func (t *AuthCon) LoginViaJWT(r *http.Request, args *models.AuthLoginJWTRequest, reply *models.AuthLoginJWTResponse) error {

	_, user := services.JwtGetUser(r)

	u := &models.User{
		Email: user,
	}
	userDb, err := services.FindByEmail(u)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	reply.User = *userDb

	return nil
}

func (t *AuthCon) ResetPassword(r *http.Request, args *models.AuthResetRequest, reply *models.GeneralMessageResponse) error {

	u := &models.User{
		Email: args.Email,
	}
	userDb, err := services.FindByEmail(u)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	resetToken, timein := services.GetTokensPassword()

	userDb.ResetToken = resetToken
	userDb.ResetTokenExpiry = timein
	services.UpdateUser(userDb)

	frontUrl := GetEnvVar("frontUrl")

	contentMsg := utils.EmailTemplateData{
		To:         userDb.Email,
		Name:       userDb.Email,
		Subject:    "Reset Password",
		Text:       "Please click this button to reset your password",
		ButtonText: "Reset",
		Link:       frontUrl + "/change/" + resetToken,
		MainText:   "We hope it all works out!",
	}
	utils.Send(contentMsg)

	reply.Message = "Success"
	EventLog(fmt.Sprint(userDb.ID), "User Password Reset", "Auth", "Reset")

	return nil
}

func (t *AuthCon) UpdatePassword(r *http.Request, args *models.AuthUpdatePasswordRequest, reply *models.GeneralMessageResponse) error {
	fmt.Println(args.ResetToken)

	userDb, err := services.FindByToken(args.ResetToken)

	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	diff := time.Now().Sub(userDb.ResetTokenExpiry)
	if diff.Seconds() > 0 {
		reply.Message = "Reset Token Expired"
		log.Print(middleware.NewError(err))
		return err
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(middleware.NewError(err))
		return err
	}

	userDb.Password = string(pass)
	userDb.ResetToken = ""

	services.UpdateUser(userDb)

	contentMsg := utils.EmailTemplateData{
		To:      userDb.Email,
		Name:    userDb.Email,
		Subject: "Password Updated",
		Text:    "Thank you for updating your password!",
	}
	utils.Send(contentMsg)

	reply.Message = "Success"
	EventLog(fmt.Sprint(userDb.ID), "User Password Changed", "Auth", "Login")

	return nil
}

func (t *AuthCon) UpdatePasswordAdmin(r *http.Request, args *models.AuthUpdatePasswordAdminRequest, reply *models.AuthUpdatePasswordAdminResponse) error {

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

	userDb, err := services.GetUserExists(&args.User)
	if err != nil {
		log.Printf("Error: %s", err)
		return nil
	}

	if userDb == true {
		reply.Messages = append(reply.Messages, "User Does Not Exist "+args.User.Email)
	} else {
		if args.User.ID != 0 {
			userOrg, _ := services.GetUserRaw(args.User)
			pass, _ := bcrypt.GenerateFromPassword([]byte(args.User.Password), bcrypt.DefaultCost)
			userOrg.Password = string(pass)
			services.UpdateUser(&userOrg)
			reply.Messages = append(reply.Messages, "User Exists, Updated "+args.User.Email)
		} else {
			reply.Messages = append(reply.Messages, "No User ID Given "+args.User.Email)
		}
	}

	return nil

}
