package services

import (
	"net/http"
)

// todo: fucking hacky, not proud of this at all
//func checkRootAdmin(email string) bool {
//	if email == dbSeed.RootAdminEmail {
//		return true
//	}
//	return false
//}



func Authorize(r *http.Request, orgId string, object string, action string) (bool, error) {
	claims := ValidateJWTRequest(r)
	user, err := FindUserByEmail(claims.Email)
	if err != nil {
		return false, err
	}
	return Enforcer.Enforce(user.ID, orgId, object, action)
}
