package main

import (
	"flag"
	"github.com/rs/zerolog/log"
	"lalela-backend/internal/pkg/logs"
	"lalela-backend/internal/pkg/mongo"
	mongoOrgsStore "lalela-backend/internal/pkg/organizations/store/mongo"
	orgsStore "lalela-backend/internal/pkg/organizations/store"

)

var configFileName = flag.String("config-file-name", "config", "specify config file")

func main() {
	flag.Parse()
	logs.Setup()

	//// get config
	config, err := GetConfig(configFileName)
	if err != nil {
		log.Fatal().Err(err).Msg("getting config from file")
	}

	// create validator
	//RequestValidator := requestValidator.New()

	// create new mongo db connection
	mongoDb, err := mongo.New(
		config.MongoDBHosts,
		config.MongoDBUsername,
		config.MongoDBPassword,
		config.MongoDBConnectionString,
		config.MongoDBName,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("creating new mongo db client")
	}
	defer func() {
		if err := mongoDb.CloseConnection(); err != nil {
			log.Error().Err(err).Msg("closing mongo db client connection")
		}
	}()

	MongoOrgStore := mongoOrgsStore.New(
		mongoDb,
	)


	_, _ = MongoOrgStore.CreateOne(orgsStore.CreateOneRequest{
		Name: "Spar Yellow Wood Park",
	})

	_, _ = MongoOrgStore.CreateOne(orgsStore.CreateOneRequest{
		Name: "Spar Woodlands East",
	})


	//MongoUserStore := mongoUsersStore.New(
	//	mongoDb,
	//)

}
