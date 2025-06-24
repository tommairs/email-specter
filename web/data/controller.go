package data

import (
	"email-specter/util"
	"github.com/gofiber/fiber/v2"
)

func GetAggregatedData(c *fiber.Ctx) error {

	from := c.Query("from")
	to := c.Query("to")

	if from == "" || to == "" {

		return c.JSON(fiber.Map{
			"success": false,
			"message": "Both 'from' and 'to' query parameters are required.",
			"data":    nil,
		})

	}

	if util.ValidateDate(from) == false || util.ValidateDate(to) == false {

		return c.JSON(fiber.Map{
			"success": false,
			"message": "Invalid date format. Please use YYYY-MM-DD.",
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
