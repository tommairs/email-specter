package data

import (
	"email-specter/util"
	"github.com/gofiber/fiber/v2"
)

func GetAggregatedDataRange(c *fiber.Ctx) error {

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

func GetBounceDataRange(c *fiber.Ctx) error {

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

	data := GetBounceDataByRange(from, to)

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
