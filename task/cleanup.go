package task

import (
	"context"
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
