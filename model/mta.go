package model

import (
	"context"
	"email-specter/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MTA struct {
	ID            string `json:"id" bson:"_id,omitempty"`
	Name          string `json:"name" bson:"name"`
	SecretToken   string `json:"secret_token" bson:"secret_token"`
	CollectionUrl string `json:"collection_url,omitempty" bson:"-"`
}

func GetMTAByID(id primitive.ObjectID) (*MTA, error) {

	collection := database.MongoConn.Collection("mtas")

	var mta MTA

	err := collection.FindOne(context.Background(), primitive.M{"_id": id}).Decode(&mta)

	if err != nil {
		return nil, err
	}

	return &mta, nil

}
