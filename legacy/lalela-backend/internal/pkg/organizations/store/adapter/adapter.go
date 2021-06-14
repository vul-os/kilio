package adapter

import (

	"net/http"
)

type OrganizationsCon struct{}


type CreateOneRequest struct {
	Name string `json:"name"`
}

type CreateOneResponse struct {
	Id string `json:"id"`
}


func (t *OrganizationsCon) CreateOne(r *http.Request, args *CreateOneRequest, reply *CreateOneResponse) error {
	//canDo, err := auth.Authorize(r, "organization", "create")
	//if err != nil {
	//	return err
	//}
	//if !canDo {
	//	claims := auth.ValidateJWTRequest(r)
	//	user, err := userStore.FindUserByEmail(claims.Email)
	//	if err != nil {
	//		return err
	//	}
	//	if !user.OrganizationId.IsZero() {
	//		fmt.Println("User has no Org, let them make one")
	//		return nil
	//	}
	//}
	//
	//_, err = store.CreateOne(args.Name)
	//if err != nil {
	//	return err
	//}
	return nil
}
