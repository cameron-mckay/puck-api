package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"otherworldly.dev/puck-api/monnitapi"
)

func InitRoutes() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/api/deleteOnNetwork", DeleteAllSensorsOnNetwork)
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
		w.Write([]byte(fmt.Sprintf("Error decoding request: %v", err)))
		return
	}
	log.Printf("Network ID %d", decodedReq.NetworkID)
	// Create request body
	sensors, err := monnitapi.GetSensorsOnNetwork(decodedReq.NetworkID)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("Error getting sensors: %v", err)))
		return
	}

	for i := range len(sensors) {
		id := sensors[i].SensorID
		log.Printf("Sensor ID %d", id)
		err = monnitapi.RemoveSensor(id)
		if err != nil {
			log.Fatalf("Error on sensor ID %d, error: %v", id, err)
		}
	}

	w.Write([]byte("Success"))
}
