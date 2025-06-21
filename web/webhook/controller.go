package webhook

import (
	"email-specter/util"
	"email-specter/web/shared"
	"github.com/gofiber/fiber/v2"
)

func ProcessWebhook(c *fiber.Ctx) error {

	var body map[interface{}]interface{}

	id := c.Params("id")
	token := c.Params("token")

	if !isAuthenticated(id, token) {
		return c.JSON(shared.ResponseMessage{Success: false, Message: "You are not authorized to access this resource."})
	}

	if err := util.ParseBodyRequest(c, &body); err != nil {
		return c.JSON(shared.ResponseMessage{Success: false, Message: util.FormatError(err)})
	}

	response := processWebhook(id, body)

	if response {
		return c.JSON(shared.ResponseMessage{Success: true, Message: "The webhook has been processed."})
	} else {
		return c.JSON(shared.ResponseMessage{Success: false, Message: "We could not process the webhook."})
	}

}
