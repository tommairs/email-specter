package model

import (
	"context"
	"email-specter/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type AggregatedData struct {
	Id                         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Date                       string             `json:"date" bson:"date"`
	Hour                       string             `json:"hour" bson:"hour"`
	MtaId                      primitive.ObjectID `json:"mta_id" bson:"mta_id"`
	SourceIP                   string             `json:"source_ip" bson:"source_ip"`
	SourceDomain               string             `json:"source_domain" bson:"source_domain"`
	DestinationService         string             `json:"destination_service" bson:"destination_service"`
	DestinationDomain          string             `json:"destination_domain" bson:"destination_domain"`
	EventType                  string             `json:"event_type" bson:"event_type"`
	KumoMtaClassification      string             `json:"kumo_mta_classification" bson:"kumo_mta_classification"`
	EmailSpecterClassification string             `json:"email_specter_classification" bson:"email_specter_classification"`
	Count                      int                `bson:"count,omitempty"`
}

func (a *AggregatedData) Upsert() {

	collection := database.MongoConn.Collection("aggregated_statistics")

	update := bson.M{
		"$inc": bson.M{"count": 1},
	}

	opts := options.Update().SetUpsert(true)

	_, err := collection.UpdateOne(context.Background(), a, update, opts)

	if err != nil {
		log.Printf("upsertAggregatedEvent failed: %v", err)
	}

}
