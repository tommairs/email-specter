package model

import (
	"context"
	"email-specter/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	Type     string    `json:"type" bson:"type"`
	Content  string    `json:"content" bson:"content"`
	Datetime time.Time `json:"datetime" bson:"datetime"`
}

type Message struct {
	ID                               primitive.ObjectID `json:"id" bson:"_id"`
	MtaId                            primitive.ObjectID `json:"mta_id" bson:"mta_id"`
	KumoMtaID                        string             `json:"kumo_mta_id" bson:"kumo_mta_id"`
	SourceIP                         string             `json:"source_ip" bson:"source_ip"`
	SourceDomain                     string             `json:"source_domain" bson:"source_domain"`
	DestinationService               string             `json:"destination_service" bson:"destination_service"`
	DestinationDomain                string             `json:"destination_domain" bson:"destination_domain"`
	Sender                           string             `json:"sender" bson:"sender"`
	Recipient                        string             `json:"recipient" bson:"recipient"`
	Events                           []Event            `json:"events" bson:"events"`
	KumoMtaBounceClassification      string             `json:"kumo_mta_bounce_classification" bson:"kumo_mta_bounce_classification"`
	EmailSpecterBounceClassification string             `json:"email_specter_bounce_classification" bson:"email_specter_bounce_classification"`
	LastStatus                       string             `json:"last_status" bson:"last_status"`
	CreatedAt                        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt                        time.Time          `json:"updated_at" bson:"updated_at"`
}

func GetMessageByKumoMtaID(kumoMtaID string) (*Message, error) {

	collection := database.MongoConn.Collection("messages")

	var message Message
	err := collection.FindOne(context.Background(), bson.M{"kumo_mta_id": kumoMtaID}).Decode(&message)

	if err != nil {
		return nil, err
	}

	return &message, nil

}

func (m *Message) Save() error {

	collection := database.MongoConn.Collection("messages")

	_, err := collection.UpdateOne(context.Background(), bson.M{"_id": m.ID}, bson.M{"$set": m})

	if err != nil {
		return err
	}

	return nil

}

func (m *Message) Insert() error {

	collection := database.MongoConn.Collection("messages")

	_, err := collection.InsertOne(context.Background(), m)

	return err

}
