//package utils
//
//import (
//	"lalela-backend/internal/pkg/middleware"
//	"fmt"
//	"github.com/jinzhu/gorm"
//	_ "github.com/jinzhu/gorm/dialects/postgres"
//	"log"
//	// _ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
//)
//
//type DBController struct {
//	DB *gorm.DB
//}
//
//var db = &DBController{}
//var dbURIString = ""
//
////ConnectDB function: Make database connection
//func InitDB() *gorm.DB {
//
//	username := GetEnvVar("databaseUser")
//	password := GetEnvVar("databasePassword")
//	databaseName := GetEnvVar("databaseName")
//	databaseHost := GetEnvVar("databaseHost")
//	databasePort := GetEnvVar("databasePort")
//
//	//Define DB connection string
//	var dbURI string
//
//	if databasePort == "" {
//		dbURI = fmt.Sprintf("host=/cloudsql/%s user=%s password=%s dbname=%s sslmode=disable", databaseHost, username, password, databaseName)
//	} else {
//		dbURI = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, databasePort, username, databaseName, password)
//	}
//	dbURIString = dbURI
//	fmt.Println("dbURI: ", dbURI)
//	err := *new(error)
//
//	//connect to db URI
//	db.DB, err = gorm.Open("postgres", dbURI)
//
//	if err != nil {
//		log.Printf("Error: %s", err)
//		log.Print(middleware.NewError(err))
//		panic(err)
//	}
//
//	// Run Seeder
//	InitSeed(db.DB)
//
//	fmt.Println("Successfully connected!", db.DB)
//	return db.DB
//}
//
//func GetDB() *gorm.DB {
//	return db.DB
//}
