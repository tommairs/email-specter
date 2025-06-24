package data

import (
	"context"
	"email-specter/database"
	"email-specter/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"sort"
)

func GetAggregatedDataByRange(from string, to string) []map[string]interface{} {

	collection := database.MongoConn.Collection("aggregated_statistics")

	matchStage := bson.D{
		{"$match", bson.D{
			{"date", bson.D{
				{"$gte", from},
				{"$lte", to},
			}},
		}},
	}

	groupStage := bson.D{
		{"$group", bson.D{
			{"_id", bson.D{
				{"date", "$date"},
				{"event_type", "$event_type"},
			}},
			{"count", bson.D{{"$sum", "$count"}}},
		}},
	}

	sortStage := bson.D{
		{"$sort", bson.D{
			{"_id.date", 1},
		}},
	}

	cursor, err := collection.Aggregate(context.TODO(), mongo.Pipeline{matchStage, groupStage, sortStage})

	if err != nil {
		log.Println("Aggregation error:", err)
		return nil
	}

	defer cursor.Close(context.TODO())

	var rawResults []struct {
		ID struct {
			Date      string `bson:"date"`
			EventType string `bson:"event_type"`
		} `bson:"_id"`
		Count int `bson:"count"`
	}

	if err = cursor.All(context.TODO(), &rawResults); err != nil {
		return nil
	}

	grouped := make(map[string][]map[string]interface{})

	for _, r := range rawResults {
		grouped[r.ID.Date] = append(grouped[r.ID.Date], map[string]interface{}{
			"event_type": r.ID.EventType,
			"count":      r.Count,
		})

	}

	var results []map[string]interface{}

	for date, events := range grouped {

		results = append(results, map[string]interface{}{
			"date":   date,
			"events": events,
		})

	}

	sort.Slice(results, func(i, j int) bool {
		return results[i]["date"].(string) < results[j]["date"].(string)
	})

	return results

}

func filterData(requestData ReportRequest) []map[string]interface{} {

	filter := bson.D{
		{"date", bson.D{
			{"$gte", requestData.From},
			{"$lte", requestData.To},
		}},
	}

	if requestData.SourceIP != "" {
		filter = append(filter, bson.E{Key: "source_ip", Value: requestData.SourceIP})
	}

	if requestData.SourceDomain != "" {
		filter = append(filter, bson.E{Key: "source_domain", Value: requestData.SourceDomain})
	}

	if requestData.DestinationDomain != "" {
		filter = append(filter, bson.E{Key: "destination_domain", Value: requestData.DestinationDomain})
	}

	if requestData.DestinationService != "" {
		filter = append(filter, bson.E{Key: "destination_service", Value: requestData.DestinationService})
	}

	if requestData.KumoMtaClassification != "" {
		filter = append(filter, bson.E{"kumo_mta_classification", requestData.KumoMtaClassification})
	}

	if requestData.EmailSpecterClassification != "" {
		filter = append(filter, bson.E{"email_specter_classification", requestData.EmailSpecterClassification})
	}

	if requestData.EventType != "" {
		filter = append(filter, bson.E{"event_type", requestData.EventType})
	}

	collection := database.MongoConn.Collection("aggregated_statistics")

	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		log.Println("Aggregation error:", err)
		return nil
	}

	defer cursor.Close(context.TODO())

	var results []map[string]interface{}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil
	}

	if requestData.GroupBy != "" {

		grouped := make(map[string]int)

		for _, r := range results {

			groupValRaw, ok := r[requestData.GroupBy]

			if !ok {
				continue
			}

			groupVal, ok := groupValRaw.(string)
			if !ok {
				continue
			}

			countRaw, ok := r["count"]

			if !ok {
				continue
			}

			count := util.EnforceInt(countRaw)

			grouped[groupVal] += count

		}

		final := make([]map[string]interface{}, 0, len(grouped))

		for key, total := range grouped {

			final = append(final, map[string]interface{}{
				requestData.GroupBy: key,
				"count":             total,
			})

		}

		sort.Slice(final, func(i, j int) bool {
			return final[i][requestData.GroupBy].(string) < final[j][requestData.GroupBy].(string)
		})

		return final

	}

	return results

}

func getProviderData(requestData ProviderDataRequest) map[string]map[string]int {

	filter := bson.D{
		{"date", bson.D{
			{"$gte", requestData.From},
			{"$lte", requestData.To},
		}},
	}

	if requestData.DestinationDomain != "" {
		filter = append(filter, bson.E{Key: "destination_domain", Value: requestData.DestinationDomain})
	}

	if requestData.DestinationService != "" {
		filter = append(filter, bson.E{Key: "destination_service", Value: requestData.DestinationService})
	}

	collection := database.MongoConn.Collection("aggregated_statistics")

	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		log.Println("Aggregation error:", err)
		return nil
	}

	defer cursor.Close(context.TODO())

	var results []map[string]interface{}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil
	}

	aggregated := make(map[string]map[string]int)

	for _, result := range results {

		date := result["date"].(string)
		eventType := result["event_type"].(string)
		count := util.EnforceInt(result["count"])

		if _, ok := aggregated[date]; !ok {
			aggregated[date] = make(map[string]int)
		}

		aggregated[date][eventType] += count

	}

	return aggregated

}
