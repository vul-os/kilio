package mongo

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	lalelaException "lalela-backend/internal/pkg/exception"
	"strings"
	"time"
)

type Database struct {
	mongoClient *mongoDriver.Client
	database    *mongoDriver.Database
}

func New(
	mongoDBHosts []string,
	mongoDBUsername,
	mongoDBPassword,
	connectionString,
	databaseName string,
) (*Database, error) {

	var db *Database
	var err error

	// try connect with a connection string if one is provided
	if connectionString != "" {
		db, err = NewFromConnectionString(connectionString)
		if err != nil {
			log.Error().Err(err).Msg("connecting to mongo")
			return nil, ErrUnexpected{}
		}
	} else if len(mongoDBHosts) != 0 {
		db, err = NewFromHosts(mongoDBHosts, mongoDBUsername, mongoDBPassword)
		if err != nil {
			log.Error().Err(err).Msg("connecting to mongo")
			return nil, ErrUnexpected{}
		}
	} else {
		err = ErrInvalidConfig{Reasons: []string{"no hosts or connection string"}}
		log.Error().Err(err).Msg("connecting to mongo")
		return nil, err
	}

	// connection successful populate and return database
	db.database = db.mongoClient.Database(databaseName)

	return db, nil
}

func NewFromHosts(mongoDBHosts []string, mongoDBUsername, mongoDBPassword string) (*Database, error) {
	log.Info().Msg(fmt.Sprintf(
		"Connecting to mongo cluster on nodes: [%s]",
		strings.Join(mongoDBHosts, ","),
	))

	// create mongo client options
	mongoOptions := &options.ClientOptions{
		Hosts: mongoDBHosts,
	}

	// if a username is provided set auth on mongo client options
	if mongoDBUsername != "" {
		mongoOptions.SetAuth(options.Credential{
			Username:      mongoDBUsername,
			Password:      mongoDBPassword,
			AuthSource:    "admin",
			PasswordSet:   true,
			AuthMechanism: "SCRAM-SHA-1",
		})
	}

	ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	mongoClient, err := mongoDriver.Connect(
		ctx,
		mongoOptions,
	)
	if err != nil {
		log.Error().Err(err).Msg("error connecting to mongo")
		return nil, err
	}

	// confirm that the client is connected
	ctx, cancelFn = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelFn()
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Error().Err(err).Msg("could not ping mongo")
		return nil, err
	} else {
		log.Info().Msg("connected to mongo")
	}

	return &Database{
		mongoClient: mongoClient,
	}, nil
}

func NewFromConnectionString(connectionString string) (*Database, error) {
	log.Info().Msg("Connecting to mongo with connection string")

	// create a new mongo client
	ctx, cancelFn := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()
	mongoClient, err := mongoDriver.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Error().Err(err).Msg("connecting to mongo")
		return nil, err
	}

	// confirm that the client is connected
	ctx, cancelFn = context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFn()
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Error().Err(err).Msg("could not ping mongo")
		return nil, err
	} else {
		log.Info().Msg("connected to mongo")
	}

	return &Database{
		mongoClient: mongoClient,
	}, nil
}

func (d *Database) CloseConnection() error {
	if err := d.mongoClient.Disconnect(context.Background()); err != nil {
		log.Error().Err(err).Msg("disconnecting from mongo Database")
		return err
	}
	return nil
}

type DBStats struct {
	// name of database
	DB string `json:"db" bson:"db"`

	// number of collections in the database
	Collections int32 `json:"collections" bson:"collections"`

	// number of documents across all collections
	Objects int32 `json:"objects" bson:"objects"`

	// average size of each object in bytes
	// NOT AFFECTED BY SCALE FACTOR
	AvgObjSize float64 `json:"avgObjSize" bson:"avgObjSize"`

	// total size of the uncompressed data held in database in bytes/scaleFactor
	DataSize float64 `json:"dataSize" bson:"dataSize"`

	// total amount of space allocated to collections in this database for document storage
	StorageSize float64 `json:"storageSize" bson:"storageSize"`

	// number of extents in the database across all collections
	NumExtents int32 `json:"numExtents" bson:"numExtents"`

	// number of indexes across all collections in the database
	Indexes int32 `json:"indexes" bson:"indexes"`

	// total size of all indexes created on this database
	IndexSize float64 `json:"indexSize" bson:"indexSize"`

	// scale used by the command
	ScaleFactor float64 `json:"scaleFactor" bson:"scaleFactor"`

	// total sized of all disk space in use on the filesystem where MongoDB stores data
	FSUsedSize float64 `json:"fsUsedSize" bson:"fsUsedSize"`

	// total size of all disk capacity on the filesystem where MongoDB stores data
	FSTotalSize float64 `json:"fsTotalSize" bson:"fsTotalSize"`

	Views int32   `json:"views" bson:"views"`
	Ok    float64 `json:"ok" bson:"ok"`
}

func (d *Database) GetDBStats() (*DBStats, error) {
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	result := d.database.RunCommand(ctx, bson.D{
		{
			Key:   "dbStats",
			Value: 1,
		},
		{
			Key:   "scale",
			Value: 1024,
		},
	})
	dbStats := new(DBStats)
	if err := result.Decode(dbStats); err != nil {
		log.Error().Err(err).Msg("unable to decode db stats result")
		return nil, lalelaException.ErrUnexpected{}
	}

	return dbStats, nil
}


func (d *Database) Collection(collectionName string) *Collection {
	return &Collection{
		driverCollection: d.database.Collection(collectionName),
	}
}