package data

import (
	"context"
	"email-specter/database"
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
