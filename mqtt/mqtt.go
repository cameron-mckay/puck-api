package mqtt

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"otherworldly.dev/puck-api/monnitapi"
	"otherworldly.dev/puck-api/websocket"
)

var mqttLog *log.Logger
var client mqtt.Client

func setupLogs() {
	mqtt.ERROR = log.New(os.Stderr, "- [mqtt][ERROR]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
	mqtt.CRITICAL = log.New(os.Stderr, "- [mqtt][CRIT]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
	mqtt.WARN = log.New(os.Stdout, "- [mqtt][WARN]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
	//mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
	mqttLog = log.New(os.Stdout, "- [mqtt][DEBUG]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
}

func createClient(broker string, username string, password string) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("coyote_" + strconv.Itoa(rand.Intn(999)))
	opts.SetUsername(username)
	opts.SetPassword(password)
	client := mqtt.NewClient(opts)
	mqttLog.Println("Connecting to broker...")
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}
	mqttLog.Println("Connected to broker.")
	return client, nil
}

func Init(broker string, username string, password string, topic string) error {
	setupLogs()
	client, err := createClient(broker, username, password)

	if err != nil {
		return err
	}

	mqttLog.Println("Subscribing to topic...")
	if token := client.Subscribe(topic, 0x0, messageHandler); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	mqttLog.Println("Subscribed to topic.")
	return nil
}

func Close() {
	client.Disconnect(0)
}

func messageHandler(client mqtt.Client, msg mqtt.Message) {
	mqttLog.Println("Parsing message...")
	raw := msg.Payload()
	var parsed monnitapi.WebhookMessage
	err := json.Unmarshal(raw, &parsed)
	if err != nil {
		mqtt.ERROR.Println("Unabled to parse json message")
		return
	}
	mqttLog.Println("Parsed message.")

	message, err := json.Marshal(parsed)

	if err != nil {
		mqtt.ERROR.Println("Unabled to encode json message")
		return
	}

	mqttLog.Println("Sending to hub...")
	websocket.MessageHub.Broadcast <- message
	mqttLog.Println("Sent to hub.")
}
