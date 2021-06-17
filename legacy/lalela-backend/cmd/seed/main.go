package main

import (
	"flag"
	"fmt"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
	"lalela-backend/internal/pkg/logs"
	"lalela-backend/internal/pkg/mongo"
	"lalela-backend/internal/pkg/organizations"
	orgs "lalela-backend/internal/pkg/organizations"
	orgsStore "lalela-backend/internal/pkg/organizations/store"
	mongoOrgsStore "lalela-backend/internal/pkg/organizations/store/mongo"
	"lalela-backend/internal/pkg/security/casbin"
	"lalela-backend/internal/pkg/users"
	userStore "lalela-backend/internal/pkg/users/store"
	mongoUsersStore "lalela-backend/internal/pkg/users/store/mongo"
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

	orgId1 := uuid.NewV4().String()
	orgId2 := uuid.NewV4().String()
	_, err = MongoOrgStore.CreateOne(orgsStore.CreateOneRequest{
		organizations.Organizations{
			ID:   orgId1,
			Name: "Spar Yellow Wood Park",
		},
	})

	if err != nil {
		log.Error().Err(err).Msg("Already seeded")
	}
	_, err = MongoOrgStore.CreateOne(orgsStore.CreateOneRequest{
		organizations.Organizations{
			ID:   orgId2,
			Name: "Spar Woodlands East",
		},
	})

	if err != nil {
		log.Error().Err(err).Msg("Already seeded")
	}

	MongoUserStore := mongoUsersStore.New(
		mongoDb,
	)

	pwdHash, err := bcrypt.GenerateFromPassword(
		[]byte("test"),
		bcrypt.DefaultCost,
	)

	imranId := uuid.NewV4().String()
	_, err = MongoUserStore.CreateOne(userStore.CreateOneRequest{
		User: users.User{
			ID:         imranId,
			Name:       "Imran",
			Email:      "imran@paruk.com",
			Password:   pwdHash,
			ResetToken: "",
		},
	})

	if err != nil {
		log.Error().Err(err).Msg("Already seeded")
	}

	ciciId := uuid.NewV4().String()
	_, err = MongoUserStore.CreateOne(userStore.CreateOneRequest{
		User: users.User{
			ID:         ciciId,
			Name:       "Cici",
			Email:      "ci@ci.com",
			Password:   pwdHash,
			ResetToken: "",
		},
	})

	if err != nil {
		log.Error().Err(err).Msg("Already seeded")
	}

	casbinEnforcer := casbin.NewCasbinEnforcer(config.CasbinModelFile, mongoDb)
	casbinEnforcer.Enforcer.EnableAutoSave(true)
	err = casbinEnforcer.Enforcer.LoadPolicy()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading policy")
	}

	servicesList := []string{
		orgsStore.OrgsServiceProvider,
		userStore.UserServiceProvider,
	}

	canDoList := []string{
		"CreateOne",
		"FindOne",
		"UpdateOne",
	}

	var orgz []orgs.Organizations
	_, err = mongoDb.Collection("organizations").FindMany(&orgz, bson.D{}, mongo.Query{})

	for _, org := range orgz {
		_, err := casbinEnforcer.Enforcer.AddGroupingPolicy(imranId, "g:admin", org.ID)
		if err != nil {
			fmt.Println(err)
		}
		_ = casbinEnforcer.Enforcer.SavePolicy()
		for _, service := range servicesList {
			for _, canDo := range canDoList {
				_, _ = casbinEnforcer.Enforcer.AddPolicy("g:admin", org.ID, service, canDo)
				_ = casbinEnforcer.Enforcer.SavePolicy()
			}
		}
	}

	for _, servicez := range servicesList {
		fmt.Println(orgz, canDoList)
		_, _ = casbinEnforcer.Enforcer.AddPolicy(ciciId, orgz[0].ID, servicez, canDoList[1])
		_ = casbinEnforcer.Enforcer.SavePolicy()
	}

	//user, domain, eft, resource, action
	//p, g:admin, RB, kanbanCards, edit
	//g, imran, g:admin, RB

	//casbinEnforcer.Enforcer.AddGroupingPolicy()
}
