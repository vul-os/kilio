package services

import (
	"lalela-backend/internal/pkg/middleware"
	"lalela-backend/internal/pkg/models"
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	jwtauth "github.com/go-chi/jwtauth"
)

var tokenAuth *jwtauth.JWTAuth

//Exception struct
type Exception models.Exception

func GetToken(user *models.User) (string, error) {
	expiresAt := time.Now().Add(time.Minute * 60).Unix()

	tk := models.Token{
		UserID: user.ID,
		Name:   user.FirstName + " " + user.LastName,
		Email:  user.Email,
		RoleID: string(user.RoleID),
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
		Type: "Login",
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	//secretJwt := os.Getenv("secretJwt")
	secretJwt := viper.Get("secretJwt").(string)

	tokenString, error := token.SignedString([]byte(secretJwt))
	if error != nil {
		log.Print(middleware.NewError(error))
		return "", error
	}

	return tokenString, nil
}

// JwtVerify Middleware function
func JwtVerify(r *http.Request) error {

	var header = r.Header.Get("x-access-token") //Grab the token from the header
	header = strings.TrimSpace(header)

	if header == "" {
		return errors.New("Missing auth ")
	}

	tk := &models.Token{}
	secretJwt := viper.Get("secretJwt").(string)

	_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretJwt), nil
	})

	if err != nil {
		log.Print(middleware.NewError(err))
		return errors.New("Token Fault ")
	}

	return nil
}

func JwtGetUser(r *http.Request) (error, string) {

	var header = r.Header.Get("x-access-token") //Grab the token from the header
	header = strings.TrimSpace(header)

	if header == "" {
		return errors.New("Missing auth "), ""
	}

	tk := &models.Token{}
	secretJwt := viper.Get("secretJwt").(string)

	_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretJwt), nil
	})

	fmt.Print(tk)
	if err != nil {
		log.Print(middleware.NewError(err))
		return errors.New("Token Fault "), ""
	}

	return nil, tk.Email
}

