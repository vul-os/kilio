package main

import (
	"cog-analytics-engine-go/internal/pkg/services"
	"github.com/gorilla/handlers"
	"log"
	//"cog-analytics-engine-go/internal/pkg/routes"
	"cog-analytics-engine-go/internal/pkg/utils"
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
	viper.SetConfigFile("./configs/env/env.toml")

	err := viper.ReadInConfig()
	if err != nil {
		log.Print(utils.NewError(err))
	}

	// Get Port
	port := utils.GetEnvVar("port")
	utils.InitPostgreDB()

	// Init RBCA Util
	utils.InitRBCA()

	// Create a new RPC server
	s := rpc.NewServer()

	s.RegisterCodec(json.NewCodec(), "application/json")

	// Register the service by creating a new JSON server
	InitRPC(s)

	r := mux.NewRouter()

	r.Handle("/rpc", s)

	// Choose the folder to serve
	staticDir := "/avatars/"

	// Cron Scheduler: Run Check of Date Processed Every 5 hours
	//c := gocron.NewScheduler(time.UTC)
	//c.Every(24).Hours().Do(cronRun)
	//c.StartAsync()

	// Create the route
	r.
		PathPrefix(staticDir).
		Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	err = http.ListenAndServe(":"+port, handlers.CORS(
		handlers.AllowedHeaders([]string{"Accept", "Accept-Language", "Content-Type", "Content-Language", "Origin", "x-access-token", "Access-Control-Allow-Origin"}),
	)(r))
	if err != nil {
		log.Print(utils.NewError(err))
	}

}

func InitRPC(s *rpc.Server) {
	err := s.RegisterService(new(services.AuthCon), "")
	if err != nil {
		log.Print(utils.NewError(err))
	}
	s.RegisterService(new(services.ConnectorsCon), "")
	s.RegisterService(new(services.DashCon), "")
	s.RegisterService(new(services.UserCon), "")
	s.RegisterService(new(services.ReportCon), "")
	s.RegisterService(new(services.PermissionCon), "")
	s.RegisterService(new(services.GroupsCon), "")
	s.RegisterService(new(services.OverviewCon), "")
	s.RegisterService(new(services.MaintenanceCon), "")
}


//// Cron: Run Check of Date Processed
//func cronRun() {
//	services.GetDateProcessedCron()
//}

// Todo :: Permissions :: Update Permissions, Remove Permissions (From Admin side)
// Todo :: Change Role Checks to include Name rather than Id
// Todo :: Refactor Pass On New Methods
// Todo :: UnitTest Coverage for RPC
