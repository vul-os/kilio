package utils

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	// _ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
)

type DBController struct {
	DB *gorm.DB
}

var db = &DBController{}
var dbURIString = ""

//ConnectDB function: Make database connection
func InitPostgreDB() *gorm.DB {

	username := GetEnvVar("databaseUser")
	password := GetEnvVar("databasePassword")
	databaseName := GetEnvVar("databaseName")
	databaseHost := GetEnvVar("databaseHost")
	databasePort := GetEnvVar("databasePort")

	//Define DB connection string
	var dbURI string
	// todo: better way to do this?
	if databasePort == "" {
		dbURI = fmt.Sprintf(
			"host=/cloudsql/%s user=%s password=%s dbname=%s sslmode=disable",
			databaseHost, username, password, databaseName)
	} else {
		dbURI = fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
			databaseHost, databasePort, username, databaseName, password)
	}
	dbURIString = dbURI
	err := *new(error)

	//connect to db URI
	db.DB, err = gorm.Open("postgres", dbURI)

	if err != nil {
		log.Printf("Error: %s", err)
		log.Print(NewError(err))
		panic(err)
	}

	// Run Seeder
	InitSeed(db.DB)

	fmt.Println("Successfully connected!", db.DB)
	return db.DB
}

func GetPostgreDB() *gorm.DB {
	return db.DB
}
