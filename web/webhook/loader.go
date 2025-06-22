package webhook

import (
	"email-specter/model"
	"encoding/json"
	"github.com/dlclark/regexp2"
	"os"
	"reflect"
	"regexp"
	"strings"
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
