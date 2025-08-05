package mqtt

import (
	"encoding/json"
	"log"
	"os"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"otherworldly.dev/puck-api/monnitapi"
)

func setupLogs() {
	mqtt.ERROR = log.New(os.Stderr, "[ERROR] ", 0)
	mqtt.CRITICAL = log.New(os.Stderr, "[CRIT] ", 0)
	mqtt.WARN = log.New(os.Stdout, "[WARN]  ", 0)
	mqtt.DEBUG = log.New(os.Stdout, "[DEBUG] ", 0)
}

func createClient(broker string, username string, password string) mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID("coyote")
	opts.SetUsername(username)
	opts.SetPassword(password)

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return client
}

func SetupMQTT(broker string, username string, password string, topic string) {
	client := createClient(broker, username, password)

	if token := client.Subscribe(topic, 0x0, messageHandler); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

func messageHandler(client mqtt.Client, msg mqtt.Message) {
	raw := msg.Payload()
	var parsed monnitapi.WebhookMessage
	err := json.Unmarshal(raw, &parsed)
	if err != nil {
		mqtt.ERROR.Println("Unabled to parse json message")
		return
	}

	// Relay the messages somewhere here

}
