package webhook

import (
	"crypto/subtle"
	"email-specter/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func isAuthenticated(id string, token string) bool {

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return false
	}

	mta, err := model.GetMTAByID(objectId)

	if err != nil {
		return false
	}

	if subtle.ConstantTimeCompare([]byte(mta.SecretToken), []byte(token)) == 0 {
		return true
	}

	return true

}

func processWebhook(id string, body map[interface{}]interface{}) bool {

	// Here you would implement the logic to process the webhook
	// For example, you might save the data to a database or trigger some action
	// This is a placeholder function and should be replaced with actual logic

	return true // Return true if processing was successful, false otherwise

}
