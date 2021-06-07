package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"lalela-backend/internal/pkg/controllers"
	"lalela-backend/internal/pkg/dbSeed"
	"lalela-backend/internal/pkg/services"
	"log"
	//"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	//"log"
	"net/http"
)

var allowedHeaders = []string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin",
	"x-access-token", "Access-Control-Allow-Origin"}

func main() {

	// Viper Get Requires Type
	viper.SetConfigFile("./configs/env/env.toml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Print(services.NewError(err))
	}

	// Get Port
	port := viper.Get("port").(string)
	mongoUri := viper.Get("mongoDbUrl").(string)
	dbName := viper.Get("dbName").(string)

	services.Client = services.InitDB(mongoUri)
	services.SecretKey = viper.Get("secretKey").(string)
	services.DbName = dbName

	services.InitCasbin(mongoUri, dbName, "./configs/casbin/model.conf")

	// Now we enable the AutoSave.
	services.Enforcer.EnableAutoSave(true)

	dbSeed.SeedDB()
	// Create a new RPC serverr
	s := rpc.NewServer()
	s.RegisterCodec(json.NewCodec(), "application/json")

	// Register the service by creating a new JSON server
	InitRPC(s)
	r := mux.NewRouter()
	r.Handle("/rpc", s)
	err = http.ListenAndServe(":" + port, handlers.CORS(handlers.AllowedHeaders(allowedHeaders))(r))
	if err != nil {
		log.Print(services.NewError(err))
	}

}

func InitRPC(s *rpc.Server) {
	err := s.RegisterService(new(controllers.UserCon), "UserCon")
	if err != nil {
		log.Print(services.NewError(err))
	}
	err = s.RegisterService(new(controllers.FormsCon), "FormsCon")
	if err != nil {
		log.Print(services.NewError(err))
	}
}

