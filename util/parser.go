package util

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"time"
)

func ParseBodyRequest[T any](c *fiber.Ctx, request *T) error {

	if err := c.BodyParser(request); err != nil {
		return NewValidationError("it seems that the request body is not in the correct format")
	}

	return nil

}

func ParseJson(jsonStr string) (map[string]interface{}, error) {

	var result map[string]interface{}

	err := json.Unmarshal([]byte(jsonStr), &result)

	if err != nil {
		return nil, err
	}

	return result, nil

}

func ParseDuration(s string) (time.Duration, error) {

	if strings.HasSuffix(s, "d") {

		daysStr := strings.TrimSuffix(s, "d")
		days, err := strconv.Atoi(daysStr)

		if err != nil {
			return 0, err
		}

		return time.Duration(days) * 24 * time.Hour, nil

	} else if strings.HasSuffix(s, "y") {

		yearsStr := strings.TrimSuffix(s, "y")
		years, err := strconv.Atoi(yearsStr)

		if err != nil {
			return 0, err
		}

		return time.Duration(years) * 365 * 24 * time.Hour, nil

	}

	return time.ParseDuration(s)

}
