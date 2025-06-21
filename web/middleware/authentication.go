package middleware

import (
	"context"
	"email-specter/database"
	"email-specter/model"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

const ExpectedLoginTokenLength = 64

func getTokenFromHeader(c *fiber.Ctx) (string, error) {

	token := c.Get("Authorization")

	if token == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	token = strings.TrimSpace(token)

	if len(token) != ExpectedLoginTokenLength {
		return "", fmt.Errorf("invalid token length")
	}

	return token, nil

}

func OnlyAuthenticatedUsers(c *fiber.Ctx) error {

	token, err := getTokenFromHeader(c)

	if err != nil {

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "You must be logged in to access this resource.",
		})

	}

	if err = loadUserData(c, token); err != nil {

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"message": "It appears that your session has expired. Please log in again.",
		})

	}

	return c.Next()

}

func findUserByToken(token string) (model.User, error) {

	collection := database.MongoConn.Collection("login_tokens")

	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: bson.M{"token": token}}},
		{{
			Key: "$lookup",
			Value: bson.M{
				"from":         "users",
				"localField":   "user_id",
				"foreignField": "_id",
				"as":           "user",
			},
		}},
		{{Key: "$unwind", Value: "$user"}},
		{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$user"}}},
	}

	ctx := context.Background()

	cursor, err := collection.Aggregate(ctx, pipeline)

	if err != nil {
		return model.User{}, fmt.Errorf("aggregation error: %w", err)
	}

	defer cursor.Close(ctx)

	if cursor.Next(ctx) {

		var user model.User

		if err := cursor.Decode(&user); err != nil {
			return model.User{}, fmt.Errorf("failed to decode user: %w", err)
		}

		return user, nil

	}

	return model.User{}, fmt.Errorf("invalid token or user not found")

}

func loadUserData(c *fiber.Ctx, token string) error {

	user, err := findUserByToken(token)

	if err != nil {
		return fmt.Errorf("failed to find user by token: %w", err)
	}

	c.Locals("user_id", user.Id)
	c.Locals("full_name", user.FullName)
	c.Locals("email_address", user.EmailAddress)
	c.Locals("token", token)

	return nil

}
