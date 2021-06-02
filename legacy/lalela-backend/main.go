package main

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/rpc"
	"github.com/gorilla/rpc/json"
	"lalela-backend/internal/pkg/controllers"
	"lalela-backend/internal/pkg/utils"
	"log"
	//"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	//"log"
	"net/http"
)

func main() {

	// Viper Get Requires Type
	viper.SetConfigFile("./configs/env/env.toml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Print(utils.NewError(err))
	}

	// Get Port
	port := viper.Get("port").(string)

	// Create a new RPC server
	s := rpc.NewServer()

	s.RegisterCodec(json.NewCodec(), "application/json")

	// Register the service by creating a new JSON server
	InitRPC(s)

	r := mux.NewRouter()

	r.Handle("/rpc", s)

	err = http.ListenAndServe(":"+port, handlers.CORS(
		handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "x-access-token", "Access-Control-Allow-Origin"}),
	)(r))
	if err != nil {
		log.Print(utils.NewError(err))
	}

}

func InitRPC(s *rpc.Server) {
	err := s.RegisterService(new(controllers.UserCon), "UserCon")
	if err != nil {
		log.Print(utils.NewError(err))
	}
}

