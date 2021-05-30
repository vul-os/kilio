package utils

import (
	"cog-analytics-engine-go/internal/pkg/controllers"
	"cog-analytics-engine-go/internal/pkg/models"
	"encoding/json"
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
type Exception controllers.Exception


type PermissionsAR struct {
	Permissions []controllers.DashboardPermission `json:"permissions"`
}

func CheckJson(w http.ResponseWriter, r *http.Request, i interface{}) (interface{}, bool) {
	// Define Struct Map
	err := json.NewDecoder(r.Body).Decode(&i)

	if err != nil {
		resp, err := GetError(400, "en")
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return i, true
		}
		json.NewEncoder(w).Encode(resp)
		return i, true
	}
	return i, false
}

func GetToken(user *models.User) (string, error) {
	expiresAt := time.Now().Add(time.Minute * 60).Unix()

	tk := controllers.Token{
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
		log.Print(NewError(error))
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

	tk := &controllers.Token{}
	secretJwt := viper.Get("secretJwt").(string)

	_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretJwt), nil
	})

	if err != nil {
		log.Print(NewError(err))
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

	tk := &controllers.Token{}
	secretJwt := viper.Get("secretJwt").(string)

	_, err := jwt.ParseWithClaims(header, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretJwt), nil
	})

	fmt.Print(tk)
	if err != nil {
		log.Print(NewError(err))
		return errors.New("Token Fault "), ""
	}

	return nil, tk.Email
}

//func IsLoginBlocked() error {
//	db := utils.GetDB()
//	var siteConfig models.SiteConfig
//	db.Where("rule = ?","isBlocked").Find(&siteConfig)
//	if siteConfig.Value == true {
//		return errors.New("isBlocked True")
//	}
//	return nil
//}
