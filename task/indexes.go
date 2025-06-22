package task

import (
	"context"
	"email-specter/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func CreateIndexes() {
	ensureMessagesIndex()
	ensureAggregatedStatisticsIndex()
}

func ensureMessagesIndex() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := database.MongoConn.Collection("messages")

	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "kumo_mta_id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{Keys: bson.D{{Key: "mta_id", Value: 1}}},
		{Keys: bson.D{{Key: "source_ip", Value: 1}}},
		{Keys: bson.D{{Key: "source_domain", Value: 1}}},
		{Keys: bson.D{{Key: "destination_service", Value: 1}}},
		{Keys: bson.D{{Key: "destination_domain", Value: 1}}},
		{Keys: bson.D{{Key: "last_status", Value: 1}}},
		{Keys: bson.D{{Key: "email_specter_bounce_classification", Value: 1}}},
		{Keys: bson.D{{Key: "kumo_mta_bounce_classification", Value: 1}}},
		{Keys: bson.D{{Key: "updated_at", Value: -1}}},
	}

	_, _ = collection.Indexes().CreateMany(ctx, indexModels)

}

func ensureAggregatedStatisticsIndex() {

	collection := database.MongoConn.Collection("aggregated_statistics")

	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{Key: "date", Value: 1},
			{Key: "hour", Value: 1},
			{Key: "mta_id", Value: 1},
			{Key: "source_ip", Value: 1},
			{Key: "source_domain", Value: 1},
			{Key: "destination_service", Value: 1},
			{Key: "destination_domain", Value: 1},
			{Key: "event_type", Value: 1},
			{Key: "kumo_mta_classification", Value: 1},
			{Key: "email_specter_classification", Value: 1},
		},
		Options: options.Index().SetUnique(true).SetName("aggregated_event_key"),
	}

	_, _ = collection.Indexes().CreateOne(context.Background(), indexModel)

}
