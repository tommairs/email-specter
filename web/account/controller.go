package account

import (
	"email-specter/util"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Login(c *fiber.Ctx) error {

	var request struct {
		EmailAddress string `json:"email_address"`
		Password     string `json:"password"`
	}

	if err := util.ParseBodyRequest(c, &request); err != nil {

		return c.JSON(map[string]interface{}{
			"success": false,
			"message": util.FormatError(err),
		})

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

		return c.JSON(map[string]interface{}{
			"success": false,
			"message": util.FormatError(err),
		})

	}

	response := createUser(request.FullName, request.EmailAddress, request.Password)

	return c.JSON(response)

}

func CanRegister(c *fiber.Ctx) error {

	canRegister := isFirstUser()

	return c.JSON(map[string]interface{}{
		"success": true,
		"message": "",
		"data": map[string]interface{}{
			"can_register": canRegister,
		},
	})

}

func Logout(c *fiber.Ctx) error {

	userId := c.Locals("user_id").(primitive.ObjectID)
	token := c.Locals("token").(string)

	logout(userId, token)

	return c.JSON(map[string]interface{}{
		"success": true,
		"message": "You have been logged out successfully.",
		"data":    nil,
	})

}

func GetAccount(c *fiber.Ctx) error {

	return c.JSON(map[string]interface{}{
		"success": true,
		"message": "",
		"data": map[string]interface{}{
			"user_id":       c.Locals("user_id").(primitive.ObjectID).Hex(),
			"full_name":     c.Locals("full_name"),
			"email_address": c.Locals("email_address"),
		},
	})

}

func ChangeFullName(c *fiber.Ctx) error {

	var request struct {
		FullName string `json:"full_name"`
	}

	if err := util.ParseBodyRequest(c, &request); err != nil {

		return c.JSON(map[string]interface{}{
			"success": false,
			"message": util.FormatError(err),
		})

	}

	response := updateFullName(c.Locals("user_id").(primitive.ObjectID), request.FullName)

	return c.JSON(response)

}

func ChangePassword(c *fiber.Ctx) error {

	var request struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}

	if err := util.ParseBodyRequest(c, &request); err != nil {

		return c.JSON(map[string]interface{}{
			"success": false,
			"message": util.FormatError(err),
		})

	}

	response := changeUserPassword(c.Locals("user_id").(primitive.ObjectID), request.CurrentPassword, request.NewPassword)

	return c.JSON(response)

}

func ChangeEmail(c *fiber.Ctx) error {

	var request struct {
		NewEmailAddress string `json:"new_email_address"`
	}

	if err := util.ParseBodyRequest(c, &request); err != nil {

		return c.JSON(map[string]interface{}{
			"success": false,
			"message": util.FormatError(err),
		})

	}

	response := changeUserEmail(c.Locals("user_id").(primitive.ObjectID), request.NewEmailAddress)

	return c.JSON(response)

}
