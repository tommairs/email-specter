package util

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
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
