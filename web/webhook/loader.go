package webhook

import (
	"email-specter/model"
	"encoding/json"
	"github.com/dlclark/regexp2"
	"os"
	"regexp"
	"strings"
)

const bounceClassificationsPath = "config/bounce_categories/bounces.txt"
const serviceProvidersSeedListPath = "config/service_providers/service_providers.json"

var ServiceProviders []model.ServiceProvider
var BounceCategories []model.BounceTypeEntry

func init() {
	loadServiceProviders()
	loadBounceCategories()
}

func loadBounceCategories() {

	contents, err := os.ReadFile(bounceClassificationsPath)

	if err != nil {
		panic("Failed to read bounce classifications: " + err.Error())
	}

	lines := strings.Split(string(contents), "\n")

	for _, line := range lines {

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		
		parts := strings.SplitN(line, ",", 4)

		if len(parts) < 4 {
			panic("Invalid bounce classification line: " + line)
		}

		regex := strings.TrimSpace(parts[0])
		_ = strings.TrimSpace(parts[1])
		category := strings.TrimSpace(parts[2])
		_ = strings.TrimSpace(parts[3])

		BounceCategories = append(BounceCategories, model.BounceTypeEntry{
			CompiledRegex: regexp2.MustCompile(regex, 0),
			Category:      category,
		})

	}

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
