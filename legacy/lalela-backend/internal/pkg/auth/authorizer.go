package auth

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/mongodb-adapter/v3"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
	"lalela-backend/internal/pkg/users/store"
	"net/http"
)

var Enforcer *casbin.Enforcer

func InitCasbin(URI string, dbName string, modelConfigPath string) {
	// Initialize a MongoDB adapter with NewAdapterWithClientOption:
	// The adapter will use custom mongo client options.
	// custom database name.
	// default collection name 'casbin_rule'.
	mongoClientOption := mongooptions.Client().ApplyURI(URI)
	a, err := mongodbadapter.NewAdapterWithClientOption(mongoClientOption, dbName)
	// Or you can use NewAdapterWithCollectionName for custom collection name.
	if err != nil {
		panic(err)
	}

	e, err := casbin.NewEnforcer(modelConfigPath, a)
	if err != nil {
		panic(err)
	}
	// Load the policy from DB.
	e.LoadPolicy()
	Enforcer = e
}

func Authorize(r *http.Request, object string, action string) (bool, error) {
	claims := ValidateJWTRequest(r)
	user, err := store.FindUserByEmail(claims.Email)
	if err != nil {
		return false, err
	}
	return Enforcer.Enforce(user.ID.Hex(), user.OrganizationId.Hex(), object, action)
}

