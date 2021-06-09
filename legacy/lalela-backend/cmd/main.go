package main

import (
	"flag"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/spf13/viper"
	jsonRPCServiceProvider "lalela-backend/internal/pkg/api/jsonRpc/service/provider"
	"lalela-backend/internal/pkg/database"
	formsAdapter "lalela-backend/internal/pkg/forms/store/adapter"
	"lalela-backend/internal/pkg/forms/store/mongo"
	"lalela-backend/internal/pkg/logs"
	"os"
	"os/signal"

	//orgAdapter "lalela-backend/internal/pkg/organizations/store/adapter"
	//submissionAdapter "lalela-backend/internal/pkg/submissions/store/adapter"
	//userAdapter "lalela-backend/internal/pkg/users/store/adapter"
	"log"
	"net/http"
)

var allowedHeaders = []string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin",
	"x-access-token", "Access-Control-Allow-Origin"}

func main() {



	database.Client = database.InitDB(mongoUri)
	//auth.SecretKey = viper.Get("secretKey").(string)
	database.DbName = dbName

	//auth.InitCasbin(mongoUri, dbName, "./cmd/casbin/model.conf")

	// Now we enable the AutoSave.
	//auth.Enforcer.EnableAutoSave(true)

	SeedDB()
	// Create a new RPC serverr
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")

	// Register the service by creating a new JSON server
	InitRPC(s)
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	err = http.ListenAndServe(":" + port, handlers.CORS(handlers.AllowedHeaders(allowedHeaders))(r))
	if err != nil {
		log.Print(err)
	}

}

func main() {
	logs.Setup()

	// Viper Get Requires Type
	viper.SetConfigFile("./cmd/env/env.toml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Print(err)
	}

	// Get Port
	port := viper.Get("port").(string)
	mongoUri := viper.Get("mongoDbUrl").(string)
	dbName := viper.Get("dbName").(string)

	//// get config
	//config, err := GetConfig(configFileName)
	//if err != nil {
	//	log.Fatal().Err(err).Msg("getting config from file")
	//}

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
	BasicTokenGenerator := basicTokenGenerator.New(
		joseSigner,
		RequestValidator,
	)
	BasicTokenValidator := basicTokenValidator.New(
		rsaKeyPair,
		RequestValidator,
	)
	BasicAuthenticator := basicAuthenticator.New(
		MongoUserStore,
		MongoRoleStore,
		BasicTokenGenerator,
		RequestValidator,
		mongoDb,
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
					configurationAdminJSONRPCAdaptor.NewAuthenticatedAdaptor(BasicConfigurationAdmin),
					shopConfigAdminJSONRPCAdaptor.NewAuthenticated(BasicShopConfigAdmin),
					shopAnalyserJSONRPCAdaptor.NewAuthenticated(BasicShopAnalyser),
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
					// mongo
					mongoAdmin.NewJSONRPCAdaptor(BasicMongoAdmin),

					// role
					roleStoreJsonRpcAdaptor.New(MongoRoleStore),

					// user
					userStoreJsonRpcAdaptor.New(MongoUserStore),
					userAdminJSONRPCAdaptor.New(BasicUserAdmin),

					// shop
					shopStoreJSONRPCAdaptor.New(MongoShopStore),
					shopAdminJSONRPCAdaptor.New(BasicShopAdmin),
					// shop config
					shopConfigAdminJSONRPCAdaptor.NewAuthorised(BasicShopConfigAdmin),

					// product
					productStoreJSONRPCAdaptor.New(MongoProductStore),
					productAdminJSONRPCAdaptor.New(BasicProductAdmin),

					// product dataPoint
					productDataPointAdminJSONRPCAdaptor.New(BasicProductDataPointAdmin),

					// job
					jobStoreJSONRPCAdaptor.New(MongoJobStore),
					jobAdminJSONRPCAdaptor.New(BasicJobAdmin),
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
