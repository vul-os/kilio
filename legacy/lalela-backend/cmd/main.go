package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"github.com/spf13/viper"
	"lalela-backend/internal/pkg/database"
	formsAdapter "lalela-backend/internal/pkg/forms/store/adapter"
	//orgAdapter "lalela-backend/internal/pkg/organizations/store/adapter"
	//submissionAdapter "lalela-backend/internal/pkg/submissions/store/adapter"
	//userAdapter "lalela-backend/internal/pkg/users/store/adapter"
	"log"
	"net/http"
)

var allowedHeaders = []string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin",
	"x-access-token", "Access-Control-Allow-Origin"}

func main() {

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

func InitRPC(s *rpc.Server) {

	services := []interface{}{
		//new(userAdapter.UserCon),
		new(formsAdapter.FormsCon),
		//new(orgAdapter.OrganizationsCon),
		//new(submissionAdapter.SubmissionsCon),
	}

	for _, service := range services {
		err := s.RegisterService(service, "")
		if err != nil {
			log.Print(err)
		}

	}
}

