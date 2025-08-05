package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"otherworldly.dev/puck-api/monnitapi"
	"otherworldly.dev/puck-api/mqtt"
	"otherworldly.dev/puck-api/routes"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	monnitapi.InitApiHandler(os.Getenv("BASE_URL"), os.Getenv("API_KEY_ID"), os.Getenv("API_KEY_SECRET"))
	mqtt.SetupMQTT(os.Getenv("MQTT_BROKER"), os.Getenv("MQTT_USERNAME"), os.Getenv("MQTT_PASSWORD"), os.Getenv("MQTT_TOPIC"))

	routes.InitRoutes()
}
