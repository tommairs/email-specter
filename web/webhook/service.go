package webhook

import (
	"crypto/subtle"
	"email-specter/model"
	"encoding/json"
	"errors"
	"github.com/dlclark/regexp2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const bounceClassificationsPath = "config/bounce_categories/bounce_categories.json"
const serviceProvidersSeedListPath = "config/service_providers/service_providers.json"

var ServiceProviders []model.ServiceProvider
var BounceCategories model.BounceCategories

func init() {
	loadServiceProviders()
	loadBounceCategories()
}

func toFieldName(key string) string {

	parts := strings.Split(key, "_")

	for i, part := range parts {
		parts[i] = strings.Title(part)
	}

	return strings.Join(parts, "")

}

func loadBounceCategories() {

	contents, err := os.ReadFile(bounceClassificationsPath)

	if err != nil {
		panic("Failed to read bounce classifications: " + err.Error())
	}

	var raw map[string][]string

	err = json.Unmarshal(contents, &raw)

	if err != nil {
		panic("Failed to parse bounce classifications: " + err.Error())
	}

	var result model.BounceCategories

	val := reflect.ValueOf(&result).Elem()

	for key, patterns := range raw {

		fieldName := toFieldName(key)

		field := val.FieldByName(fieldName)

		if !field.IsValid() || !field.CanSet() {
			panic("Unknown or unassignable field: " + fieldName)
		}

		var compiled []regexp2.Regexp

		for _, pattern := range patterns {

			re, err := regexp2.Compile(pattern, 0)

			if err != nil {
				panic("Invalid regex in " + key + ": " + err.Error())
			}

			compiled = append(compiled, *re)

		}

		field.Set(reflect.ValueOf(compiled))

	}

	BounceCategories = result

}

func loadServiceProviders() {

	contents, err := os.ReadFile(serviceProvidersSeedListPath)

	if err != nil {
		panic("Failed to read service providers seed list: " + err.Error())
	}

	err = json.Unmarshal(contents, &ServiceProviders)

	if err != nil {
		panic("Failed to parse service providers seed list: " + err.Error())
	}

	for i, provider := range ServiceProviders {
		ServiceProviders[i].CompiledRegex = regexp.MustCompile(provider.Regex)
	}

}

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
		return false
	}

	return true

}

func findCaseInsensitiveKey(headers map[string]string, key string) string {

	for k, v := range headers {

		if strings.ToLower(k) == strings.ToLower(key) {
			return v
		}

	}

	return ""

}

func processWebhook(mtaId string, webhookData model.WebhookEvent) bool {

	mtaIdObject, err := primitive.ObjectIDFromHex(mtaId)

	if err != nil {
		return false
	}

	switch webhookData.Type {

	case "Reception":
		return handleReceptionEvent(mtaIdObject, webhookData)

	case "Delivery":
		return handleDeliveryEvent(mtaIdObject, webhookData)

	case "TransientFailure":
		return handleTransientFailureEvent(mtaIdObject, webhookData)

	case "Bounce":
		return handleBounceEvent(mtaIdObject, webhookData)

	}

	return false

}

func getIPAddress(ipPort string) string {

	host, _, err := net.SplitHostPort(ipPort)

	if err != nil {

		if strings.Contains(ipPort, ":") && strings.Count(ipPort, ":") > 1 {
			return ipPort
		}

		return strings.Split(ipPort, ":")[0]

	}

	return host

}

func getDomain(email string) string {

	atIndex := strings.Index(email, "@")

	if atIndex == -1 {
		return ""
	}

	return email[atIndex+1:]

}

func getServiceName(peerName string, domain string) string {

	peerName = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(peerName), "."))
	domain = strings.ToLower(strings.TrimSuffix(strings.TrimSpace(domain), "."))

	for _, provider := range ServiceProviders {

		if provider.CompiledRegex.MatchString(peerName) || provider.CompiledRegex.MatchString(domain) {
			return provider.Name
		}

	}

	return "Unknown"

}

func handleBounceEvent(mtaId primitive.ObjectID, webhookData model.WebhookEvent) bool {

	currentTime := time.Now()

	event := model.Event{
		Type:     webhookData.Type,
		Content:  webhookData.Response.Content,
		Datetime: currentTime,
	}

	message, err := getOrCreateMessage(mtaId, webhookData, currentTime)

	if err != nil {
		return false
	}

	return updateMessageStatus(webhookData, message, event, webhookData.Type, currentTime)

}

func handleTransientFailureEvent(mtaId primitive.ObjectID, webhookData model.WebhookEvent) bool {

	currentTime := time.Now()

	event := model.Event{
		Type:     webhookData.Type,
		Content:  webhookData.Response.Content,
		Datetime: currentTime,
	}

	message, err := getOrCreateMessage(mtaId, webhookData, currentTime)

	if err != nil {
		return false
	}

	return updateMessageStatus(webhookData, message, event, webhookData.Type, currentTime)

}

func handleDeliveryEvent(mtaId primitive.ObjectID, webhookData model.WebhookEvent) bool {

	currentTime := time.Now()

	event := model.Event{
		Type:     webhookData.Type,
		Content:  webhookData.Response.Content,
		Datetime: currentTime,
	}

	message, err := getOrCreateMessage(mtaId, webhookData, currentTime)

	if err != nil {
		return false
	}

	return updateMessageStatus(webhookData, message, event, webhookData.Type, currentTime)

}

func handleReceptionEvent(mtaId primitive.ObjectID, webhookData model.WebhookEvent) bool {

	currentTime := time.Now()

	events := []model.Event{
		{
			Type:     webhookData.Type,
			Content:  "",
			Datetime: currentTime,
		},
	}

	message := createMessage(mtaId, currentTime, webhookData)

	message.Events = events

	err := message.Insert()

	if err != nil {
		return false
	}

	return true

}

func createMessage(mtaId primitive.ObjectID, currentTime time.Time, webhookData model.WebhookEvent) *model.Message {

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

func getFullBounceMessage(webhookData model.WebhookEvent) string {

	if webhookData.Type != "TransientFailure" && webhookData.Type != "Bounce" {
		return ""
	}

	smtpResponseCode := webhookData.Response.Code
	enhancedCode := webhookData.Response.EnhancedCode // C.S.D
	content := webhookData.Response.Content

	bounceMessage := strconv.Itoa(smtpResponseCode) + " " + content

	if enhancedCode != nil {
		bounceMessage = strconv.Itoa(smtpResponseCode) + " " + strconv.Itoa(enhancedCode.Class) + "." + strconv.Itoa(enhancedCode.Subject) + "." + strconv.Itoa(enhancedCode.Detail) + " " + content
	}

	return bounceMessage

}

func categorizeBounce(webhookData model.WebhookEvent) string {

	bounceMessage := getFullBounceMessage(webhookData)

	if bounceMessage == "" {
		return ""
	}

	for category, patterns := range BounceCategories.GetCategories() {

		for _, pattern := range patterns {

			if match, _ := pattern.MatchString(bounceMessage); match {
				return category
			}

		}

	}

	return "Other"

}

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

	return message.Save() == nil

}

func getOrCreateMessage(mtaId primitive.ObjectID, webhookData model.WebhookEvent, currentTime time.Time) (*model.Message, error) {

	message, err := model.GetMessageByKumoMtaID(webhookData.ID)

	if err != nil {

		if errors.Is(err, mongo.ErrNoDocuments) {

			message = createMessage(mtaId, currentTime, webhookData)

			if err := message.Insert(); err != nil {
				return nil, err
			}

			return message, nil

		}

		return nil, err

	}

	return message, nil

}
