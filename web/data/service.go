package data

import (
	"context"
	"email-specter/database"
	"email-specter/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"sort"
	"time"
)

func GetAggregatedDataByRange(from string, to string) []map[string]interface{} {
	
	collection := database.MongoConn.Collection("aggregated_statistics")

	matchStage := bson.D{}

	if from != "" || to != "" {

		dateRange := bson.D{}

		if from != "" {
			dateRange = append(dateRange, bson.E{Key: "$gte", Value: util.ConvertYmdToTime(from)})
		}

		if to != "" {
			toTime := util.ConvertYmdToTime(to).Add(23*time.Hour + 59*time.Minute + 59*time.Second + 999*time.Millisecond)
			dateRange = append(dateRange, bson.E{Key: "$lte", Value: toTime})
		}

		matchStage = append(matchStage, bson.E{Key: "date", Value: dateRange})

	}

	pipeline := mongo.Pipeline{}

	if len(matchStage) > 0 {
		pipeline = append(pipeline, bson.D{{"$match", matchStage}})
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

	pipeline = append(pipeline, groupStage, sortStage)

	cursor, err := collection.Aggregate(context.TODO(), pipeline)

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

func getProviderClassificationData(requestData ProviderClassificationRequest) []map[string]interface{} {

	collection := database.MongoConn.Collection("aggregated_statistics")

	filter := bson.D{
		{"date", bson.D{
			{"$gte", requestData.From},
			{"$lte", requestData.To},
		}},
		{
			"event_type", requestData.EventType,
		},
	}

	if requestData.DestinationDomain != "" {
		filter = append(filter, bson.E{Key: "destination_domain", Value: requestData.DestinationDomain})
	}

	if requestData.DestinationService != "" {
		filter = append(filter, bson.E{Key: "destination_service", Value: requestData.DestinationService})
	}

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

	// Group results by both kumo_mta_classification and email_specter_classification

	grouped := make(map[string]map[string]int)

	for _, result := range results {

		kumoClassification := result["kumo_mta_classification"].(string)
		emailClassification := result["email_specter_classification"].(string)

		count := util.EnforceInt(result["count"])

		if _, ok := grouped["kumo_mta_classification"]; !ok {
			grouped["kumo_mta_classification"] = make(map[string]int)
		}

		if _, ok := grouped["email_specter_classification"]; !ok {
			grouped["email_specter_classification"] = make(map[string]int)
		}

		grouped["kumo_mta_classification"][kumoClassification] += count
		grouped["email_specter_classification"][emailClassification] += count

	}

	var finalResults []map[string]interface{}

	for classificationType, classifications := range grouped {

		for classification, count := range classifications {

			finalResults = append(finalResults, map[string]interface{}{
				"classification_type": classificationType,
				"classification":      classification,
				"count":               count,
			})

		}

	}

	sort.Slice(finalResults, func(i, j int) bool {
		return finalResults[i]["classification"].(string) < finalResults[j]["classification"].(string)
	})

	return finalResults

}

func getDateRangeForLastMonth() (string, string) {

	endDate := time.Now().Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -30).Format("2006-01-02")

	return startDate, endDate

}

func aggregateEntityData(startDate, endDate, fieldName string, limit int, requireNonEmpty bool) []string {

	collection := database.MongoConn.Collection("aggregated_statistics")

	matchStage := bson.D{
		{"date", bson.D{
			{"$gte", startDate},
			{"$lte", endDate},
		}},
	}

	if requireNonEmpty {

		matchStage = append(matchStage, bson.E{
			Key:   fieldName,
			Value: bson.D{{"$ne", ""}},
		})

	}

	pipeline := mongo.Pipeline{
		{{"$match", matchStage}},
		{{"$group", bson.D{
			{"_id", "$" + fieldName},
			{"count", bson.D{{"$sum", "$count"}}},
		}}},
		{{"$sort", bson.D{
			{"count", -1},
		}}},
	}

	if limit > 0 {
		pipeline = append(pipeline, bson.D{{"$limit", limit}})
	}

	cursor, err := collection.Aggregate(context.TODO(), pipeline)

	if err != nil {
		return nil
	}

	defer cursor.Close(context.TODO())

	type aggregationResult struct {
		ID    string `bson:"_id"`
		Count int    `bson:"count"`
	}

	var rawResults []aggregationResult

	if err = cursor.All(context.TODO(), &rawResults); err != nil {
		return nil
	}

	results := make([]string, 0, len(rawResults))

	for _, item := range rawResults {

		if item.ID != "" {
			results = append(results, item.ID)
		}

	}

	return results

}

func getTopEntities() map[string]interface{} {

	startDate, endDate := getDateRangeForLastMonth()

	result := make(map[string]interface{})

	entityConfigs := []struct {
		resultKey       string
		fieldName       string
		limit           int
		requireNonEmpty bool
	}{
		{"destination_services", "destination_service", 1000, true},
		{"source_ips", "source_ip", 0, true},
		{"source_domains", "source_domain", 1000, true},
		{"destination_domains", "destination_domain", 1000, true},
		{"kumo_mta_bounce_classification", "kumo_mta_bounce_classification", 1000, true},
	}

	for _, config := range entityConfigs {

		if data := aggregateEntityData(startDate, endDate, config.fieldName, config.limit, config.requireNonEmpty); data != nil && len(data) > 0 {
			result[config.resultKey] = data
		}

	}

	return result

}

func buildMessageFilter(request GetMessagesRequest) bson.D {

	filter := bson.D{}

	if request.From != "" || request.To != "" {

		dateRange := bson.D{}

		if request.From != "" {
			dateRange = append(dateRange, bson.E{Key: "$gte", Value: util.ConvertYmdToTime(request.From)})
		}

		if request.To != "" {
			toTime := util.ConvertYmdToTime(request.To).Add(23*time.Hour + 59*time.Minute + 59*time.Second + 999*time.Millisecond)
			dateRange = append(dateRange, bson.E{Key: "$lte", Value: toTime})
		}

		filter = append(filter, bson.E{Key: "updated_at", Value: dateRange})

	}

	if request.MtaId > 0 {
		filter = append(filter, bson.E{Key: "mta_id", Value: request.MtaId})
	}

	if request.SourceIP != "" {
		filter = append(filter, bson.E{Key: "source_ip", Value: request.SourceIP})
	}

	if request.SourceDomain != "" {
		filter = append(filter, bson.E{Key: "source_domain", Value: request.SourceDomain})
	}

	if request.DestinationService != "" {
		filter = append(filter, bson.E{Key: "destination_service", Value: request.DestinationService})
	}

	if request.DestinationDomain != "" {
		filter = append(filter, bson.E{Key: "destination_domain", Value: request.DestinationDomain})
	}

	if request.LastStatus != "" {
		filter = append(filter, bson.E{Key: "last_status", Value: request.LastStatus})
	}

	if request.EmailSpecterBounceClassification != "" {
		filter = append(filter, bson.E{Key: "email_specter_bounce_classification", Value: request.EmailSpecterBounceClassification})
	}

	if request.KumoMtaBounceClassification != "" {
		filter = append(filter, bson.E{Key: "kumo_mta_bounce_classification", Value: request.KumoMtaBounceClassification})
	}

	return filter

}

func countMessages(filter bson.D) int64 {

	collection := database.MongoConn.Collection("messages")

	totalCount, err := collection.CountDocuments(context.TODO(), filter)

	if err != nil {
		log.Println("Error counting messages:", err)
		return 0
	}

	return totalCount

}

func getMessages(request GetMessagesRequest) map[string]interface{} {

	collection := database.MongoConn.Collection("messages")

	filter := buildMessageFilter(request)

	totalCount := countMessages(filter)

	const messagesPerPage = 100

	page := request.Page

	if page < 1 {
		page = 1
	}

	skip := (page - 1) * messagesPerPage
	totalPages := (totalCount + messagesPerPage - 1) / messagesPerPage

	findOptions := options.Find()
	findOptions.SetLimit(messagesPerPage)
	findOptions.SetSkip(int64(skip))
	findOptions.SetSort(bson.D{{"updated_at", -1}})

	cursor, err := collection.Find(context.TODO(), filter, findOptions)

	if err != nil {

		log.Println("Error retrieving messages:", err)

		return map[string]interface{}{
			"messages":    []map[string]interface{}{},
			"total_count": 0,
			"page":        page,
			"pages":       0,
		}

	}

	defer cursor.Close(context.TODO())

	var messages []map[string]interface{}

	if err = cursor.All(context.TODO(), &messages); err != nil {

		log.Println("Error retrieving messages:", err)

		return map[string]interface{}{
			"messages":    []map[string]interface{}{},
			"total_count": 0,
			"page":        page,
			"pages":       0,
		}

	}

	return map[string]interface{}{
		"messages":    messages,
		"total_count": totalCount,
		"page":        page,
		"pages":       totalPages,
	}

}
