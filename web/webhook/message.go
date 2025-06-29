package webhook

import (
	"email-specter/model"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func updateMessageStatus(webhookData model.WebhookEvent, message *model.Message, event model.Event, status string, currentTime time.Time) bool {

	message.Events = append(message.Events, event)

	message.LastStatus = status
	message.UpdatedAt = currentTime

	if status == "Bounce" || status == "TransientFailure" {

		message.KumoMtaBounceClassification = webhookData.BounceClassification
		message.EmailSpecterBounceClassification = categorizeBounce(webhookData)

	} else {

		message.KumoMtaBounceClassification = ""
		message.EmailSpecterBounceClassification = ""

	}

	// Good to overwrite it because the destination service, source IP cannot be determined at the time of reception

	message.DestinationService = getServiceName(webhookData.PeerAddress.Name, message.DestinationDomain)
	message.SourceIP = getIPAddress(webhookData.SourceAddress.Address)

	result := message.Save() == nil

	if result {
		go upsertAggregatedEvent(message.MtaId, message, currentTime)
	}

	return result

}

func getOrCreateMessage(mtaId primitive.ObjectID, webhookData model.WebhookEvent, currentTime time.Time) (*model.Message, error) {

	message, err := model.GetMessageByKumoMtaID(webhookData.ID)

	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {

			message = createMessageObject(mtaId, currentTime, webhookData)

			if err := message.Insert(); err != nil {
				return nil, err
			}

			return message, nil

		}

		return nil, err

	}

	return message, nil

}

func createMessageObject(mtaId primitive.ObjectID, currentTime time.Time, webhookData model.WebhookEvent) *model.Message {

	sourceIpAddress := getIPAddress(webhookData.SourceAddress.Address)
	sourceDomain := getDomain(webhookData.Sender)
	receiverDomain := getDomain(webhookData.Recipient)

	message := model.Message{
		ID:                               primitive.NewObjectID(),
		MtaId:                            mtaId,
		KumoMtaID:                        webhookData.ID,
		SourceIP:                         sourceIpAddress,
		SourceDomain:                     sourceDomain,
		DestinationService:               getServiceName(webhookData.PeerAddress.Name, sourceDomain),
		DestinationDomain:                receiverDomain,
		Sender:                           webhookData.Sender,
		Recipient:                        webhookData.Recipient,
		Events:                           []model.Event{},
		KumoMtaBounceClassification:      webhookData.BounceClassification,
		EmailSpecterBounceClassification: categorizeBounce(webhookData),
		LastStatus:                       webhookData.Type,
		CreatedAt:                        currentTime,
		UpdatedAt:                        currentTime,
	}

	return &message

}
