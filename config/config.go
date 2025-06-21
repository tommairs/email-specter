package config

import (
	"github.com/joho/godotenv"
	"os"
	"time"
)

var MongoConnStr string
var MongoDb string

var SessionLength time.Duration

var FrontendUrl string

var BackendUrl string
var HttpPort string
var ListenAddress string

func loadConfig() {

	err := godotenv.Load(".env")

	if err != nil {
		panic("Error loading .env file: " + err.Error())
		return
	}

	MongoConnStr = os.Getenv("MONGO_CONN_STR")
	MongoDb = os.Getenv("MONGO_DB")

	SessionLength, _ = time.ParseDuration(os.Getenv("SESSION_LENGTH"))

	FrontendUrl = os.Getenv("FRONTEND_URL")

	BackendUrl = os.Getenv("BACKEND_URL")
	HttpPort = os.Getenv("HTTP_PORT")
	ListenAddress = os.Getenv("LISTEN_ADDRESS")

}

func init() {
	loadConfig()
}
