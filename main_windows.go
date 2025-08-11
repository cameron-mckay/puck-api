package main

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"otherworldly.dev/puck-api/db"
	"otherworldly.dev/puck-api/monnitapi"
	"otherworldly.dev/puck-api/mqtt"
	"otherworldly.dev/puck-api/routes"
	"otherworldly.dev/puck-api/websocket"
)

var mainLog *log.Logger
var mainError *log.Logger

type puckService struct{}

func runService(name string, isDebug bool) {
	if isDebug {
		err := debug.Run(name, &puckService{})
		if err != nil {
			log.Fatalln("Error running service in debug mode.")
		}
	} else {
		err := svc.Run(name, &puckService{})
		if err != nil {
			log.Fatalln("Error running service in Service Control mode.")
		}
	}
}

func (m *puckService) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (bool, uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

	status <- svc.Status{State: svc.StartPending}

	status <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		c := <-r
		switch c.Cmd {
		case svc.Interrogate:
			status <- c.CurrentStatus
		case svc.Stop, svc.Shutdown:
			log.Print("Shutting service...!")
			break loop
		default:
			log.Printf("Unexpected service control request #%d", c)
		}
	}

	status <- svc.Status{State: svc.StopPending}
	return false, 1
}

func main() {
	mainLog = log.New(os.Stdout, "- [main][DEBUG]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
	mainError = log.New(os.Stderr, "- [main][ERROR]: ", log.Ldate|log.Ltime|log.Lmsgprefix)

	err := godotenv.Load(".env")
	if err != nil {
		mainError.Fatalf("Error loading .env file: %v", err)
	}

	accId, err := strconv.ParseInt(os.Getenv("API_ACCOUNT_ID"), 10, 32)
	if err != nil {
		mainError.Fatalf("Error parsing account id: %v", err)
	}

	err = db.Init(os.Getenv("DB_CONNECTION_STRING"))
	if err != nil {
		mainError.Fatal(err)
	}
	websocket.Init()
	monnitapi.Init(os.Getenv("API_BASE_URL"), os.Getenv("API_KEY_ID"), os.Getenv("API_KEY_SECRET"), int(accId))
	mqtt.Init(os.Getenv("MQTT_BROKER"), os.Getenv("MQTT_USERNAME"), os.Getenv("MQTT_PASSWORD"), os.Getenv("MQTT_TOPIC"))
	routes.Init(os.Getenv("HTTP_LISTEN_ADDR"))

	defer routes.Close()
	defer mqtt.Close()
	defer monnitapi.Close()
	defer websocket.Close()
	defer db.Close()

	runService("myservice", false)
}
