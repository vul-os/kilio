package main

import (
	controllers "lalela-backend/internal/pkg/controllers"
	"lalela-backend/internal/pkg/middleware"
	"github.com/gorilla/handlers"
	"log"
	//"lalela-backend/internal/pkg/routes"
	"lalela-backend/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	//"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	//"log"
	"net/http"
)

func main() {

	// Viper Get Requires Type
	viper.SetConfigFile("./internal/config/env/env.toml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Print(middleware.NewError(err))
	}

	// Get Port
	port := controllers.GetEnvVar("port")
	utils.InitDB()

	// Init RBCA Util
	utils.InitRBCA()

	// Create a new RPC server
	s := rpc.NewServer()

	s.RegisterCodec(json.NewCodec(), "application/json")

	// Register the service by creating a new JSON server
	controllers.InitRPC(s)

	r := mux.NewRouter()

	r.Handle("/rpc", s)

	err = http.ListenAndServe(":"+port, handlers.CORS(
		handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin",
			"x-access-token", "Access-Control-Allow-Origin"}),
	)(r))
	if err != nil {
		log.Print(middleware.NewError(err))
	}

}
