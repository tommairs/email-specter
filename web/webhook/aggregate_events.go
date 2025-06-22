package webhook

import (
	"context"
	"email-specter/database"
	"email-specter/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func upsertAggregatedEvent(mtaId primitive.ObjectID, message *model.Message, currentTime time.Time) {

	collection := database.MongoConn.Collection("aggregated_statistics")

	hour := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), 0, 0, 0, time.UTC,
	)

	date := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		0, 0, 0, 0, time.UTC,
	)

	filter := bson.M{
		"date":                         date.Format("2006-01-02"),
		"hour":                         hour.Format("2006-01-02 15:00:00"),
		"mta_id":                       mtaId,
		"source_ip":                    message.SourceIP,
		"source_domain":                message.SourceDomain,
		"destination_service":          message.DestinationService,
		"destination_domain":           message.DestinationDomain,
		"event_type":                   message.LastStatus,
		"kumo_mta_classification":      message.KumoMtaBounceClassification,
		"email_specter_classification": message.EmailSpecterBounceClassification,
	}

	update := bson.M{
		"$inc": bson.M{"count": 1},
	}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.Background(), filter, update, opts)

	if err != nil {
		log.Printf("upsertAggregatedEvent failed: %v", err)
	}

}
