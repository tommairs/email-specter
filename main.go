package main

import (
	"context"
	"email-specter/config"
	"email-specter/database"
	"email-specter/task"
	"email-specter/web/account"
	"email-specter/web/middleware"
	"email-specter/web/mta"
	"email-specter/web/webhook"
	"github.com/go-co-op/gocron/v2"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const maxBodySize = 5 * 1024 * 1024 * 1024

func runWebserver(shutdownCtx context.Context) {

	app := fiber.New(fiber.Config{
		Prefork:   false,
		BodyLimit: maxBodySize,
	})

	app.Use(func(c *fiber.Ctx) error {

		// Maybe replace with this an actual domain name later.

		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Method() == fiber.MethodOptions {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()

	})

	app.Get("/", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Welcome to the Email Specter API!",
		})

	})

	// Has the admin user been created?

	app.Get("/can-register", account.CanRegister)

	app.Post("/register", account.Register)
	app.Post("/login", account.Login)
	app.Get("/account", middleware.OnlyAuthenticatedUsers, account.GetAccount)
	app.Patch("/account/change-full-name", middleware.OnlyAuthenticatedUsers, account.ChangeFullName)
	app.Patch("/account/change-password", middleware.OnlyAuthenticatedUsers, account.ChangePassword)
	app.Patch("/account/change-email", middleware.OnlyAuthenticatedUsers, account.ChangeEmail)
	app.Post("/logout", middleware.OnlyAuthenticatedUsers, account.Logout)

	// Add an MTA

	app.Get("/mta", middleware.OnlyAuthenticatedUsers, mta.GetAllMTAs)
	app.Post("/mta", middleware.OnlyAuthenticatedUsers, mta.AddMTA)
	app.Patch("/mta/:id", middleware.OnlyAuthenticatedUsers, mta.EditMTA)
	app.Delete("/mta/:id", middleware.OnlyAuthenticatedUsers, mta.DeleteMTA)
	app.Post("/mta/:id/rotate-secret-token", middleware.OnlyAuthenticatedUsers, mta.RotateSecretToken)

	// Webhook Collector

	app.Post("/webhook/:id/:token", webhook.ProcessWebhook)

	// Not Found Handler

	app.All("*", func(c *fiber.Ctx) error {

		return c.JSON(fiber.Map{
			"success": false,
			"message": "Quo vadis, amicus?",
		})

	})

	go func() {

		if err := app.Listen(config.ListenAddress + ":" + config.HttpPort); err != nil {
			log.Fatalf("Error starting HTTP server: %v", err)
		}

	}()

	<-shutdownCtx.Done()

	_ = app.Shutdown()

	log.Println("Webserver shutting down...")

}

func handleScheduleError(_ gocron.Job, err error) {

	if err != nil {
		panic(err)
	}

}

func runScheduler(shutdownCtx context.Context) {

	s, err := gocron.NewScheduler()

	if err != nil {
		log.Fatalf("Error starting scheduler: %v", err)
	}

	handleScheduleError(s.NewJob(gocron.DurationJob(1*time.Hour), gocron.NewTask(task.CleanLoginTokens)))
	handleScheduleError(s.NewJob(gocron.DurationJob(1*time.Hour), gocron.NewTask(task.CleanMessages)))
	handleScheduleError(s.NewJob(gocron.DurationJob(24*time.Hour), gocron.NewTask(task.CleanAggregatedData)))

	s.Start()
	<-shutdownCtx.Done()
	s.Shutdown()

	log.Println("Scheduler shutting down...")

}

func boot() {
	database.CreateDatabaseConnections()
}

func main() {

	shutdownCtx, shutdownCancel := context.WithCancel(context.Background())
	defer shutdownCancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	defer signal.Stop(signalChan)

	go boot()
	go runScheduler(shutdownCtx)
	go runWebserver(shutdownCtx)

	<-signalChan

	shutdownCancel()

	log.Println("Received shutdown signal, shutting down...")

}
