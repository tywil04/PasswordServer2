package database

import (
	"context"
	"os"

	psErrors "PasswordServer2/lib/errors"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client
var Database *mongo.Database

func DatabaseConnect() {
	environment := os.Getenv("ENVIRONMENT")

	if environment == "" {
		panic(psErrors.ErrorEnvironmentEnvNotFound)
	} else if environment == "testing" || environment == "development" || environment == "production" {
		client, clientError := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_DB_URI")))
		if clientError != nil {
			panic(psErrors.ErrorLoadingDatabase)
		}

		Database = client.Database(environment)
		Client = client
	} else {
		panic(psErrors.ErrorEnvironmentEnvInvalid)
	}

	LoadCollections()
}
