package mta

import (
	"email-specter/util"
	"email-specter/web/shared"
	"github.com/gofiber/fiber/v2"
)

func GetAllMTAs(c *fiber.Ctx) error {

	response := getAllMTAs()

	return c.JSON(response)

}

func AddMTA(c *fiber.Ctx) error {

	var request struct {
		Name string `json:"name"`
	}

	if err := util.ParseBodyRequest(c, &request); err != nil {
		return c.JSON(shared.ResponseMessage{Success: false, Message: util.FormatError(err)})
	}

	response := addMTA(request.Name)

	return c.JSON(response)

}

func EditMTA(c *fiber.Ctx) error {

	var request struct {
		Name string `json:"name"`
	}

	if err := util.ParseBodyRequest(c, &request); err != nil {
		return c.JSON(shared.ResponseMessage{Success: false, Message: util.FormatError(err)})
	}

	mtaID := c.Params("id")
	response := editMTA(mtaID, request.Name)

	return c.JSON(response)

}

func DeleteMTA(c *fiber.Ctx) error {

	mtaID := c.Params("id")
	response := deleteMTA(mtaID)

	return c.JSON(response)

}

func RotateSecretToken(c *fiber.Ctx) error {

	mtaID := c.Params("id")
	response := rotateSecretToken(mtaID)

	return c.JSON(response)

}
