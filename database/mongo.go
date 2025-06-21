package database

import (
	"context"
	"email-specter/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoConn *mongo.Database

func getMongoConnection() *mongo.Database {

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mongoOptions := options.Client().ApplyURI(config.MongoConnStr).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.Background(), mongoOptions)

	if err != nil {
		panic(err)
	}

	return client.Database(config.MongoDb)

}
