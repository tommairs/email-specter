package model

import (
	"context"
	"email-specter/database"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id" json:"id"`
	FullName     string             `json:"full_name" bson:"full_name"`
	EmailAddress string             `json:"email_address" bson:"email_address"`
	PasswordHash string             `json:"password_hash,omitempty" bson:"password_hash"`
}

func GetUserBy(column string, value string) (*User, error) {

	collection := database.MongoConn.Collection("users")

	filter := bson.M{column: value}

	if column == "_id" {

		id, err := primitive.ObjectIDFromHex(value)

		if err != nil {
			return nil, err
		}

		filter = bson.M{"_id": id}
		
	}

	var user User

	err := collection.FindOne(context.Background(), filter).Decode(&user)

	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}

		return nil, err

	}

	return &user, nil

}
