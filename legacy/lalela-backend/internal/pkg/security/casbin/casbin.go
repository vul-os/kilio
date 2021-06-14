package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/rs/zerolog/log"
	"lalela-backend/internal/pkg/mongo"
)

type Casbin struct {
	Enforcer *casbin.Enforcer
}

func NewCasbinEnforcer(
	configPath string,
	database *mongo.Database,
) *Casbin {
	a, err := NewAdapter(database)
	if err != nil {
		log.Fatal().Err(err).Msg("error setting up casbin adapter")
	}
	e, err := casbin.NewEnforcer(configPath, a)
	if err != nil {
		log.Fatal().Err(err).Msg("error setting up enforcer")
	}
	// Load the policy from DB.
	_ = e.LoadPolicy()
	return &Casbin{
		Enforcer: e,
	}
}

//func Authorize(r *http.Request, object string, action string) (bool, error) {
//	claims := ValidateJWTRequest(r)
//	user, err := store.FindUserByEmail(claims.Email)
//	if err != nil {
//		return false, err
//	}
//	return Enforcer.Enforce(user.ID.Hex(), user.OrganizationId.Hex(), object, action)
//}