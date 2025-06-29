package data

import (
	"email-specter/util"
	"github.com/gofiber/fiber/v2"
)

func GetAggregatedData(c *fiber.Ctx) error {

	from := c.Query("from")
	to := c.Query("to")

	if (from != "" && util.ValidateDate(from) == false) || (to != "" && util.ValidateDate(to) == false) {

		return c.JSON(fiber.Map{
			"success": false,
			"message": "When 'from' or 'to' is provided, they must be in the format YYYY-MM-DD.",
			"data":    nil,
		})

	}

	data := GetAggregatedDataByRange(from, to)

	if data == nil {

		return c.JSON(fiber.Map{
			"success": false,
			"message": "No data found for the specified date range.",
			"data":    nil,
		})

	}

	return c.JSON(fiber.Map{
		"success": true,
		"message": "",
		"data":    data,
	})

}

type ReportRequest struct {
	From                       string `json:"from"`
	To                         string `json:"to"`
	SourceIP                   string `json:"source_ip"`
	SourceDomain               string `json:"source_domain"`
	DestinationDomain          string `json:"destination_domain"`
	DestinationService         string `json:"destination_service"`
	KumoMtaClassification      string `json:"kumo_mta_classification"`
	EmailSpecterClassification string `json:"email_specter_classification"`
	EventType                  string `json:"event_type"`
	GroupBy                    string `json:"group_by"`
}

func GenerateReport(c *fiber.Ctx) error {

	var request ReportRequest

	if err := util.ParseBodyRequest(c, &request); err != nil {

		return c.JSON(map[string]interface{}{
			"success": false,
			"message": util.FormatError(err),
		})

	}

	if util.ValidateDate(request.From) == false || util.ValidateDate(request.To) == false {

		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid date format. Please use YYYY-MM-DD.",
			"data":    nil,
		})

	}

	data := filterData(request)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "",
		"data":    data,
	})

}

type ProviderDataRequest struct {
	From               string `json:"from"`
	To                 string `json:"to"`
	DestinationDomain  string `json:"destination_domain"`
	DestinationService string `json:"destination_service"`
}

func GetProviderData(c *fiber.Ctx) error {

	var request ProviderDataRequest

	if err := util.ParseBodyRequest(c, &request); err != nil {

		return c.JSON(map[string]interface{}{
			"success": false,
			"message": util.FormatError(err),
		})

	}

	if util.ValidateDate(request.From) == false || util.ValidateDate(request.To) == false {

		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid date format. Please use YYYY-MM-DD.",
			"data":    nil,
		})

	}

	if request.DestinationDomain == "" && request.DestinationService == "" {

		return c.JSON(fiber.Map{
			"success": false,
			"message": "Either 'destination_domain' or 'destination_service' must be provided.",
			"data":    nil,
		})

	}

	data := getProviderData(request)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "",
		"data":    data,
	})

}

type ProviderClassificationRequest struct {
	From               string `json:"from"`
	To                 string `json:"to"`
	EventType          string `json:"event_type"`
	DestinationDomain  string `json:"destination_domain"`
	DestinationService string `json:"destination_service"`
}

func GetProviderClassificationData(c *fiber.Ctx) error {

	var request ProviderClassificationRequest

	if err := util.ParseBodyRequest(c, &request); err != nil {

		return c.JSON(map[string]interface{}{
			"success": false,
			"message": util.FormatError(err),
		})

	} else if util.ValidateDate(request.From) == false || util.ValidateDate(request.To) == false {

		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid date format. Please use YYYY-MM-DD.",
			"data":    nil,
		})

	} else if request.EventType != "TransientFailure" && request.EventType != "Bounce" {

		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid event type. Allowed values are 'TransientFailure' or 'Bounce'.",
			"data":    nil,
		})

	} else {

		data := getProviderClassificationData(request)

		return c.JSON(fiber.Map{
			"success": true,
			"message": "",
			"data":    data,
		})

	}

}

func GetTopEntities(c *fiber.Ctx) error {

	data := getTopEntities()

	return c.JSON(fiber.Map{
		"success": true,
		"message": "",
		"data":    data,
	})

}

type GetMessagesRequest struct {
	From                             string `json:"from"`
	To                               string `json:"to"`
	MtaId                            int    `json:"mta_id"`
	SourceIP                         string `json:"source_ip"`
	SourceDomain                     string `json:"source_domain"`
	DestinationService               string `json:"destination_service"`
	DestinationDomain                string `json:"destination_domain"`
	LastStatus                       string `json:"last_status"`
	EmailSpecterBounceClassification string `json:"email_specter_bounce_classification"`
	KumoMtaBounceClassification      string `json:"kumo_mta_bounce_classification"`
	Page                             int    `json:"page"`
}

func GetMessages(c *fiber.Ctx) error {

	var request GetMessagesRequest

	if err := util.ParseBodyRequest(c, &request); err != nil {

		return c.JSON(map[string]interface{}{
			"success": false,
			"message": util.FormatError(err),
		})

	}

	if (request.From != "" && util.ValidateDate(request.From) == false) || (request.To != "" && util.ValidateDate(request.To) == false) {

		return c.JSON(fiber.Map{
			"success": false,
			"message": "If 'from' or 'to' is provided, they must be in the format YYYY-MM-DD.",
			"data":    nil,
		})

	}

	if request.Page < 1 {
		request.Page = 1
	}

	data := getMessages(request)

	return c.JSON(fiber.Map{
		"success": true,
		"message": "",
		"data":    data,
	})

}
