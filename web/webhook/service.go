package webhook

import (
	"crypto/subtle"
	"email-specter/model"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net"
	"os"
	"regexp"
	"strings"
	"time"
)

const serviceProvidersSeedListPath = "config/service_providers.json"

var ServiceProviders []model.ServiceProvider

func init() {
	loadServiceProviders()
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
		return true
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

func getServiceName(mxRecord string, domain string) string {

	for _, provider := range ServiceProviders {

		if provider.CompiledRegex.MatchString(mxRecord) || provider.CompiledRegex.MatchString(domain) {
			return provider.Name
		}

	}

	return "Unknown"

}

func handleDeliveryEvent(mtaId primitive.ObjectID, webhookData model.WebhookEvent) bool {

}

func handleReceptionEvent(mtaId primitive.ObjectID, webhookData model.WebhookEvent) bool {

	currentTime := time.Now()

	sourceIpAddress := getIPAddress(webhookData.SourceAddress.Address)
	sourceDomain := getDomain(webhookData.Sender)
	receiverDomain := getDomain(webhookData.Recipient)

	events := []model.Event{
		{
			Type:     "Accepted",
			Content:  "",
			Datetime: currentTime,
		},
	}

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
		Events:                           events,
		KumoMtaBounceClassification:      "",
		EmailSpecterBounceClassification: "",
		LastStatus:                       "Accepted",
		CreatedAt:                        currentTime,
		UpdatedAt:                        currentTime,
	}

	err := message.Save()

	if err != nil {
		return false
	}

	return true

}
