package auth

import (
	//"fmt"
	//"lalela-backend/internal/pkg/database"
	//"lalela-backend/internal/pkg/users/store"
	"net/http"
)

type AuthenticatorCon struct{}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	JWT string
}

func (t *AuthenticatorCon) Login(r *http.Request, args *LoginRequest, reply *LoginResponse) error {
	jwtToken, err := Login(args.Email, args.Password)
	if err != nil {
		return err
	}
	reply.JWT = jwtToken
	return nil
}

func Login(email string, password string) (string, error) {
	//var collection = database.OpenCollection("user")
	//
	//foundUser, err := store.FindUserByEmail(email)
	//if err != nil {
	//	fmt.Println("User Not Found")
	//	return "", err
	//}
	//
	//passwordIsValid := VerifyPassword(password, foundUser.Password)
	//if !passwordIsValid {
	//	return "", nil
	//}
	//
	//jwtToken, err := GenerateToken(email)
	//if err != nil {
	//	fmt.Println("Generating token error")
	//	return "", err
	//}
	//
	//UpdateToken(collection, email, jwtToken)
	return "jwtToken", nil
}
