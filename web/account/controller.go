package account

import (
	"email-specter/util"
	"email-specter/web/shared"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Login(c *fiber.Ctx) error {

	var request struct {
		EmailAddress string `json:"email_address"`
		Password     string `json:"password"`
	}

	if err := util.ParseBodyRequest(c, &request); err != nil {
		return c.JSON(shared.ResponseMessage{Success: false, Message: util.FormatError(err)})
	}

	response := authenticateUser(request.EmailAddress, request.Password)

	return c.JSON(response)

}

func Register(c *fiber.Ctx) error {

	var request struct {
		FullName     string `json:"full_name"`
		EmailAddress string `json:"email_address"`
		Password     string `json:"password"`
	}

	if err := util.ParseBodyRequest(c, &request); err != nil {
		return c.JSON(shared.ResponseMessage{Success: false, Message: util.FormatError(err)})
	}

	response := createUser(request.FullName, request.EmailAddress, request.Password)

	return c.JSON(response)

}

func Logout(c *fiber.Ctx) error {

	userId := c.Locals("user_id").(primitive.ObjectID)
	token := c.Locals("token").(string)

	logout(userId, token)

	return c.JSON(shared.ResponseMessage{Success: true, Message: "You have been successfully logged out."})

}
