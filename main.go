package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"otherworldly.dev/puck-api/monnitapi"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	monnitapi.InitApiHandler(os.Getenv("BASE_URL"), os.Getenv("API_KEY_ID"), os.Getenv("API_KEY_SECRET"))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/deleteOnNetwork", DeleteAllSensorsOnNetwork)
	http.ListenAndServe(":42069", nil)
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
}


type DeleteSensorsOnNetwork struct {
	NetworkID int `json:"networkId"`
}
func DeleteAllSensorsOnNetwork(w http.ResponseWriter, r *http.Request) { 
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	decodedReq := new(DeleteSensorsOnNetwork)
	
	err := json.NewDecoder(r.Body).Decode(decodedReq)

	// Check for errors
	if err != nil {
		log.Fatalf("Error decoding request: %v", err)
		w.Write([]byte("Error decoding request"))
	}
	// Create request body
	sensors := monnitapi.GetSensorsOnNetwork(decodedReq.NetworkID)

	for i := range len(sensors.Result) {
		id := sensors.Result[i].SensorID
		res := monnitapi.RemoveSensor(id)
		if res == false {
			log.Fatalf("Error on sensor ID %d", id)
		}
	}

	w.Write([]byte("Success"))
}

