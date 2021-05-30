package services

import (
	"cog-analytics-engine-go/internal/pkg/controllers"
	models "cog-analytics-engine-go/internal/pkg/models"
	"cog-analytics-engine-go/internal/pkg/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthCon struct{}

func (t *AuthCon) Login(r *http.Request, args *controllers.AuthLoginRequest, reply *controllers.AuthLoginResponse) error {
	fmt.Println(args.Email, args.Password)

	//blocked := services.IsLoginBlocked()
	//if blocked != nil {
	//	log.Print(middleware.NewError(blocked))
	//	return blocked
	//}

	u := &models.User{
		Email: args.Email,
	}
	userDb, err := utils.FindByEmail(u)
	if err != nil {
		log.Print(utils.NewError(err))
		return err
	}

	errPw := utils.CheckPassword(userDb, args.Password)
	if errPw != nil {
		log.Print(utils.NewError(errPw))
		return errPw
	}

	token, err := utils.GetToken(userDb)
	if err != nil {
		log.Print(utils.NewError(err))
		return err
	}
	utils.SetAvatar(userDb, r)
	reply.Jwt = token
	reply.User = *userDb
	utils.EventLog(fmt.Sprint(userDb.ID), "User Logged In", "Auth", "Login")

	return nil
}

func (t *AuthCon) LoginViaJWT(r *http.Request, args *controllers.AuthLoginJWTRequest, reply *controllers.AuthLoginJWTResponse) error {

	_, user := utils.JwtGetUser(r)

	u := &models.User{
		Email: user,
	}
	userDb, err := utils.FindByEmail(u)
	if err != nil {
		log.Print(utils.NewError(err))
		return err
	}

	reply.User = *userDb

	return nil
}

func (t *AuthCon) ResetPassword(r *http.Request, args *controllers.AuthResetRequest, reply *controllers.GeneralMessageResponse) error {

	u := &models.User{
		Email: args.Email,
	}
	userDb, err := utils.FindByEmail(u)
	if err != nil {
		log.Print(utils.NewError(err))
		return err
	}

	resetToken, timein := utils.GetTokensPassword()

	userDb.ResetToken = resetToken
	userDb.ResetTokenExpiry = timein
	utils.UpdateUser(userDb)

	frontUrl := utils.GetEnvVar("frontUrl")

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
	utils.EventLog(fmt.Sprint(userDb.ID), "User Password Reset", "Auth", "Reset")

	return nil
}

func (t *AuthCon) UpdatePassword(r *http.Request, args *controllers.AuthUpdatePasswordRequest, reply *controllers.GeneralMessageResponse) error {
	fmt.Println(args.ResetToken)

	userDb, err := utils.FindByToken(args.ResetToken)

	if err != nil {
		log.Print(utils.NewError(err))
		return err
	}

	diff := time.Now().Sub(userDb.ResetTokenExpiry)
	if diff.Seconds() > 0 {
		reply.Message = "Reset Token Expired"
		log.Print(utils.NewError(err))
		return err
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(args.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Print(utils.NewError(err))
		return err
	}

	userDb.Password = string(pass)
	userDb.ResetToken = ""

	utils.UpdateUser(userDb)

	contentMsg := utils.EmailTemplateData{
		To:      userDb.Email,
		Name:    userDb.Email,
		Subject: "Password Updated",
		Text:    "Thank you for updating your password!",
	}
	utils.Send(contentMsg)

	reply.Message = "Success"
	utils.EventLog(fmt.Sprint(userDb.ID), "User Password Changed", "Auth", "Login")

	return nil
}

func (t *AuthCon) UpdatePasswordAdmin(r *http.Request, args *controllers.AuthUpdatePasswordAdminRequest, reply *controllers.AuthUpdatePasswordAdminResponse) error {

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
			userOrg, _ := utils.GetUserRaw(args.User)
			pass, _ := bcrypt.GenerateFromPassword([]byte(args.User.Password), bcrypt.DefaultCost)
			userOrg.Password = string(pass)
			utils.UpdateUser(&userOrg)
			reply.Messages = append(reply.Messages, "User Exists, Updated "+args.User.Email)
		} else {
			reply.Messages = append(reply.Messages, "No User ID Given "+args.User.Email)
		}
	}

	return nil

}
