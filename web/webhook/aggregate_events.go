package webhook

import (
	"email-specter/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

func upsertAggregatedEvent(mtaId primitive.ObjectID, message *model.Message, currentTime time.Time) {
	
	hour := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		currentTime.Hour(), 0, 0, 0, time.UTC,
	)

	date := time.Date(
		currentTime.Year(), currentTime.Month(), currentTime.Day(),
		0, 0, 0, 0, time.UTC,
	)

	aggregatedData := model.AggregatedData{
		Date:                       date.Format("2006-01-02"),
		Hour:                       hour.Format("2006-01-02 15:00:00"),
		MtaId:                      mtaId,
		SourceIP:                   message.SourceIP,
		SourceDomain:               message.SourceDomain,
		DestinationService:         message.DestinationService,
		DestinationDomain:          message.DestinationDomain,
		EventType:                  message.LastStatus,
		KumoMtaClassification:      message.KumoMtaBounceClassification,
		EmailSpecterClassification: message.EmailSpecterBounceClassification,
	}

	aggregatedData.Upsert()

}
