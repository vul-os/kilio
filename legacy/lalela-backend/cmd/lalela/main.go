package main

import (
	"flag"
	. "fmt"
	jsonRpcHttpServer "lalela-backend/internal/pkg/api/jsonRpc/server/http"
	jsonRPCServiceProvider "lalela-backend/internal/pkg/api/jsonRpc/service/provider"

	"gopkg.in/square/go-jose.v2"
	formsJSONRPCAdapter "lalela-backend/internal/pkg/forms/store/adapter"
	mongoFormsStore "lalela-backend/internal/pkg/forms/store/mongo"
	"lalela-backend/internal/pkg/key"
	"lalela-backend/internal/pkg/logs"
	"lalela-backend/internal/pkg/mongo"
	userJSONRPCAdapter "lalela-backend/internal/pkg/users/store/adapter"
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
	flag.Parse()			//used to parse command line input for usage
	logs.Setup()			//logging with JSON output via zerologs package

	//// get config
	config, err := GetConfig(configFileName) //
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
	UsersStore := mongoUsersStore.New(
		mongoDb,
	)
	log.Info().Msg("FUUUUUUUUUUUUUUUUUUUUUUUUUUUCK")
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
	//MongoUserStore := mongoUserStore.New(
	//	RequestValidator,
	//	mongoDb,
	//)
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
	Println(joseSigner)
	//BasicTokenGenerator := basicTokenGenerator.New(
	//	joseSigner,
	//	RequestValidator,
	//)
	//BasicTokenValidator := basicTokenValidator.New(
	//	rsaKeyPair,
	//	RequestValidator,
	//)
	//BasicAuthenticator := basicAuthenticator.New(
	//	MongoUserStore,
	//	MongoRoleStore,
	//	BasicTokenGenerator,
	//	RequestValidator,
	//	mongoDb,
	//)
	//authenticationMiddleware := middleware.NewAuthentication(
	//	BasicTokenValidator,
	//)
	//authorisationMiddleware := middleware.NewAuthorisation(
	//	BasicAuthenticator,
	//)

	// create rpc http server
	server := jsonRpcHttpServer.New(
		"localhost",
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
					formsJSONRPCAdapter.New(FormsStore),
					userJSONRPCAdapter.New(UsersStore),
				},
			},

			////
			//// Authenticated API Server
			////
			//{
			//	Name: "Authenticated",
			//	Path: "/api/authenticated",
			//	Middleware: []func(http.Handler) http.Handler{
			//		authenticationMiddleware.Apply,
			//	},
			//	ServiceProviders: []jsonRPCServiceProvider.Provider{
			//		configurationAdminJSONRPCAdaptor.NewAuthenticatedAdaptor(BasicConfigurationAdmin),
			//		shopConfigAdminJSONRPCAdaptor.NewAuthenticated(BasicShopConfigAdmin),
			//		shopAnalyserJSONRPCAdaptor.NewAuthenticated(BasicShopAnalyser),
			//	},
			//},
			//
			////
			//// Authenticated and Authorised API Server
			////
			//{
			//	Name: "Authenticated and Authorised",
			//	Path: "/api/authorised",
			//	Middleware: []func(http.Handler) http.Handler{
			//		authenticationMiddleware.Apply,
			//		authorisationMiddleware.Apply,
			//	},
			//	ServiceProviders: []jsonRPCServiceProvider.Provider{
			//		// mongo
			//		mongoAdmin.NewJSONRPCAdaptor(BasicMongoAdmin),
			//
			//		// role
			//		roleStoreJsonRpcAdaptor.New(MongoRoleStore),
			//
			//		// user
			//		userStoreJsonRpcAdaptor.New(MongoUserStore),
			//		userAdminJSONRPCAdaptor.New(BasicUserAdmin),
			//
			//		// shop
			//		shopStoreJSONRPCAdaptor.New(MongoShopStore),
			//		shopAdminJSONRPCAdaptor.New(BasicShopAdmin),
			//		// shop config
			//		shopConfigAdminJSONRPCAdaptor.NewAuthorised(BasicShopConfigAdmin),
			//
			//		// product
			//		productStoreJSONRPCAdaptor.New(MongoProductStore),
			//		productAdminJSONRPCAdaptor.New(BasicProductAdmin),
			//
			//		// product dataPoint
			//		productDataPointAdminJSONRPCAdaptor.New(BasicProductDataPointAdmin),
			//
			//		// job
			//		jobStoreJSONRPCAdaptor.New(MongoJobStore),
			//		jobAdminJSONRPCAdaptor.New(BasicJobAdmin),
			//	},
			//},
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
