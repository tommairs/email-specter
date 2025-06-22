package webhook

import (
	"email-specter/model"
	"net"
	"strconv"
	"strings"
)

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

func findCaseInsensitiveKey(headers map[string]string, key string) string {

	for k, v := range headers {

		if strings.ToLower(k) == strings.ToLower(key) {
			return v
		}

	}

	return ""

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

	for _, entry := range BounceCategories.GetCategories() {

		for _, pattern := range entry.Patterns {

			if match, _ := pattern.MatchString(bounceMessage); match {
				return entry.Name
			}

		}

	}

	return "Other"

}
