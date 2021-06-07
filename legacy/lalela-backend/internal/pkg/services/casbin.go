package services

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/mongodb-adapter/v3"
	mongooptions "go.mongodb.org/mongo-driver/mongo/options"
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