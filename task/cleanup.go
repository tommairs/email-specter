package task

import (
	"context"
	"email-specter/config"
	"email-specter/database"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func CleanLoginTokens() error {

	collection := database.MongoConn.Collection("login_tokens")

	filter := bson.M{
		"expires_at": bson.M{"$lt": time.Now()},
	}

	_, err := collection.DeleteMany(context.Background(), filter)

	if err != nil {
		return err
	}

	return nil

}

func CleanMessages() error {

	collection := database.MongoConn.Collection("messages")

	filter := bson.M{
		"created_at": bson.M{"$lt": time.Now().Add(-config.LogRetentionPeriod)},
	}

	_, err := collection.DeleteMany(context.Background(), filter)

	if err != nil {
		return err
	}

	return nil

}

func CleanAggregatedData() error {

	collection := database.MongoConn.Collection("data")

	filter := bson.M{
		"date": bson.M{"$lt": time.Now().Add(-config.DataRetentionPeriod).Format("2006-01-02")},
	}

	_, err := collection.DeleteMany(context.Background(), filter)

	if err != nil {
		return err
	}

	return nil

}
