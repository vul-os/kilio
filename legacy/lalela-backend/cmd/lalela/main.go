package main

import (
	"flag"
	"fmt"
	jsonRpcHttpServer "lalela-backend/internal/pkg/api/jsonRpc/server/http"
	jsonRPCServiceProvider "lalela-backend/internal/pkg/api/jsonRpc/service/provider"
	"lalela-backend/internal/pkg/middleware"
	"lalela-backend/internal/pkg/security/casbin"
	formsJSONRPCAdaptor "lalela-backend/internal/pkg/forms/store/adapter"
	mongoFormsStore "lalela-backend/internal/pkg/forms/store/mongo"
	"gopkg.in/square/go-jose.v2"
	"lalela-backend/internal/pkg/logs"
	"lalela-backend/internal/pkg/mongo"
	authenticatorJSONRPCAdaptor "lalela-backend/internal/pkg/security/authenticator/adaptor/jsonRpc"
	basicAuthenticator "lalela-backend/internal/pkg/security/authenticator/basic"
	"lalela-backend/internal/pkg/security/key"
	basicTokenGenerator "lalela-backend/internal/pkg/security/token/generator/basic"
	basicTokenValidator "lalela-backend/internal/pkg/security/token/validator/basic"
	userStoreJsonRpcAdaptor "lalela-backend/internal/pkg/users/store/adaptor/jsonRPC"
	mongoUsersStore "lalela-backend/internal/pkg/users/store/mongo"
	"os"
	"os/signal"

	//orgAdapter "lalela-backend/internal/pkg/organizations/store/adapter"
	//submissionAdapter "lalela-backend/internal/pkg/submissions/store/adapter"
	//userAdapter "lalela-backend/internal/pkg/users/store/adapter"
	"github.com/rs/zerolog/log"
	"net/http"
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

	//
	// Service Providers
	//
	FormsStore := mongoFormsStore.New(
		mongoDb,
	)
	//// Mongo
	//BasicMongoAdmin := mongoAdmin.NewBasicAdmin(
	//	RequestValidator,
	//	mongoDb,
	//)
	//
	//// Role
	//MongoRoleStore := mongoRoleStore.New(
	//	RequestValidator,
	//	mongoDb,
	//)
	//
	//// User
	MongoUserStore := mongoUsersStore.New(
		mongoDb,
	)
	//BasicUserValidator := basicUserValidator.New(
	//	RequestValidator,
	//	MongoUserStore,
	//	MongoRoleStore,
	//)
	//BasicUserAdmin := basicUserAdmin.New(
	//	RequestValidator,
	//	BasicUserValidator,
	//	MongoUserStore,
	//	MongoRoleStore,
	//)

	//
	// authentication and authorisation middleware
	//
	// fetch or generate RSA key pair
	rsaKeyPair, err := key.ParseRSAPrivateKeyFromString(config.PrivateKeyString)
	if err != nil {
		log.Fatal().Err(err).Msg("parsing private key")
	}

	// create a new signer using RSASSA-PSS (SHA512) with the given private key.
	joseSigner, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.PS512, Key: rsaKeyPair}, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("generating new jose signer")
	}
	fmt.Println(joseSigner)

	BasicTokenGenerator := basicTokenGenerator.New(
		joseSigner,
	)
	BasicTokenValidator := basicTokenValidator.New(
		rsaKeyPair,
	)

	casbinEnforcer := casbin.NewCasbinEnforcer(config.CasbinModelFile,  mongoDb)
	err = casbinEnforcer.Enforcer.LoadPolicy()
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading policy")
	}

	BasicAuthenticator := basicAuthenticator.New(
		MongoUserStore,
		BasicTokenGenerator,
		mongoDb,
		casbinEnforcer,
	)

	authenticationMiddleware := middleware.NewAuthentication(
		BasicTokenValidator,
	)
	authorisationMiddleware := middleware.NewAuthorisation(
		BasicAuthenticator,
	)




	// create rpc http server
	server := jsonRpcHttpServer.New(
		"0.0.0.0",
		config.ServerPort,
		[]jsonRpcHttpServer.RPCServerConfig{
			//
			// Public API Server
			//
			{
				Name:       "Public",
				Path:       "/api/public",
				Middleware: []func(http.Handler) http.Handler{},
				ServiceProviders: []jsonRPCServiceProvider.Provider{
					authenticatorJSONRPCAdaptor.New(BasicAuthenticator),
				},
			},

			//
			// Authenticated API Server
			//
			{
				Name: "Authenticated",
				Path: "/api/authenticated",
				Middleware: []func(http.Handler) http.Handler{
					authenticationMiddleware.Apply,
				},
				ServiceProviders: []jsonRPCServiceProvider.Provider{

				},
			},

			//
			// Authenticated and Authorised API Server
			//
			{
				Name: "Authenticated and Authorised",
				Path: "/api/authorised",
				Middleware: []func(http.Handler) http.Handler{
					authenticationMiddleware.Apply,
					authorisationMiddleware.Apply,
				},
				ServiceProviders: []jsonRPCServiceProvider.Provider{
					//// mongo
					//mongoAdmin.NewJSONRPCAdaptor(BasicMongoAdmin),
					//
					//// role
					//roleStoreJsonRpcAdaptor.New(MongoRoleStore),
					//
					//// user
					userStoreJsonRpcAdaptor.New(MongoUserStore),
					formsJSONRPCAdaptor.New(FormsStore),

					//userAdminJSONRPCAdaptor.New(BasicUserAdmin),

				},
			},
		},
	)

	// start server
	go func() {
		if err := server.Start(); err != nil {
			log.Error().Err(err).Msg("json rpc http api server has stopped")
		}
	}()

	// wait for interrupt signal to stop
	systemSignalsChannel := make(chan os.Signal, 1)
	signal.Notify(systemSignalsChannel, os.Interrupt)
	for s := range systemSignalsChannel {
		log.Info().Msgf("Application is shutting down.. ( %s )", s)
		return
	}
}
