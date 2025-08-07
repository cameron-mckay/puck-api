package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/joho/godotenv"
	"otherworldly.dev/puck-api/db"
	"otherworldly.dev/puck-api/monnitapi"
	"otherworldly.dev/puck-api/mqtt"
	"otherworldly.dev/puck-api/routes"
	"otherworldly.dev/puck-api/websocket"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	accId, err := strconv.ParseInt(os.Getenv("API_ACCOUNT_ID"), 10, 32)
	if err != nil {
		log.Fatalf("Error parsing account id: %v", err)
	}

	err = db.Init(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		log.Fatal(err)
	}
	websocket.Init()
	monnitapi.Init(os.Getenv("API_BASE_URL"), os.Getenv("API_KEY_ID"), os.Getenv("API_KEY_SECRET"), int(accId))
	mqtt.Init(os.Getenv("MQTT_BROKER"), os.Getenv("MQTT_USERNAME"), os.Getenv("MQTT_PASSWORD"), os.Getenv("MQTT_TOPIC"))
	routes.Init(os.Getenv("HTTP_LISTEN_ADDR"))

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Application running. Press Ctrl+C to exit.")

	// Block until interrupt
	<-sigChan

	routes.Close()
	mqtt.Close()
	monnitapi.Close()
	websocket.Close()
	db.Close()

	log.Println("Received interrupt signal. Exiting...")
}
