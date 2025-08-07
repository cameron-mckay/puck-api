package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"otherworldly.dev/puck-api/db"
	"otherworldly.dev/puck-api/monnitapi"
	"otherworldly.dev/puck-api/websocket"
)

var routeLog *log.Logger
var routeError *log.Logger

func reqBodyParser[T any](r *http.Request) (*T, error) {
	defer r.Body.Close()
	decodedReq := new(T)
	err := json.NewDecoder(r.Body).Decode(decodedReq)
	if err != nil {
		return nil, err
	}
	return decodedReq, nil
}

func Init(addr string) {
	routeLog = log.New(os.Stdout, "- [routes][DEBUG]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
	routeError = log.New(os.Stderr, "- [routes][ERROR]: ", log.Ldate|log.Ltime|log.Lmsgprefix)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", websocket.Serve)
	http.HandleFunc("/api/getNetworks", getNetworks)
	http.HandleFunc("/api/deleteAllSensorsOnNetwork", deleteAllSensorsOnNetwork)
	http.HandleFunc("/api/getSensorsOnNetwork", getSensorsOnNetwork)
	http.HandleFunc("/api/addBinToNetwork", addBinToNetwork)

	go http.ListenAndServe(addr, nil)

	routeLog.Println("Listening on " + addr)
}

func Close() {
	// nothing here yet
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
	routeLog.Println("Served index")
}

func getNetworks(w http.ResponseWriter, r *http.Request) {
	networks, err := monnitapi.GetNetworkList()

	if err != nil {
		routeError.Println(err)
		w.Write([]byte(fmt.Sprintf("Error fetching network info: %v", err)))
		return
	}

	res, err := json.Marshal(networks)

	if err != nil {
		routeError.Println(err)
		w.Write([]byte(fmt.Sprintf("Error encoding network json: %v", err)))
		return
	}

	w.Write(res)
}

type networkReq struct {
	NetworkID int `json:"networkId"`
}

func getSensorsOnNetwork(w http.ResponseWriter, r *http.Request) {
	routeLog.Println("Sensor list requested")
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	decodedReq := new(networkReq)

	err := json.NewDecoder(r.Body).Decode(decodedReq)

	// Check for errors
	if err != nil {
		routeError.Println(err)
		w.Write([]byte(fmt.Sprintf("Error decoding request: %v", err)))
		return
	}
	routeLog.Printf("Network ID %d", decodedReq.NetworkID)
	// Create request body
	sensors, err := monnitapi.GetSensorsOnNetwork(decodedReq.NetworkID)
	if err != nil {
		routeError.Println(err)
		w.Write([]byte(fmt.Sprintf("Error getting sensors: %v", err)))
		return
	}

	res, err := json.Marshal(sensors)

	if err != nil {
		routeError.Println(err)
		w.Write([]byte(fmt.Sprintf("Error encoding sensor json: %v", err)))
		return
	}

	w.Write(res)
}

func deleteAllSensorsOnNetwork(w http.ResponseWriter, r *http.Request) {
	routeLog.Println("Sensor delete requested")
	w.Header().Set("Content-Type", "application/json")

	decodedReq, err := reqBodyParser[networkReq](r)
	// Check for errors
	if err != nil {
		routeError.Println(err)
		w.Write([]byte(fmt.Sprintf("Error decoding request: %v", err)))
		return
	}

	err = db.DeleteAllSensorsOnNetwork(decodedReq.NetworkID)

	if err != nil {
		routeError.Printf("can't delete sensors, error: %v", err)
		w.Write([]byte(fmt.Sprintf("Error deleting sensors: %v", err)))
		return
	}

	w.Write([]byte("Success"))
}

type binReq struct {
	BinID     int `json:"binId"`
	NetworkID int `json:"networkId"`
}

func addBinToNetwork(w http.ResponseWriter, r *http.Request) {
	routeLog.Println("Sensor delete requested")
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	decodedReq, err := reqBodyParser[binReq](r)
	// Check for errors
	if err != nil {
		routeError.Println(err)
		w.Write([]byte(fmt.Sprintf("Error decoding request: %v", err)))
		return
	}

	err = db.AddBinToNetwork(decodedReq.BinID, decodedReq.NetworkID)

	if err != nil {
		routeError.Printf("can't add sensors, error: %v", err)
		w.Write([]byte(fmt.Sprintf("Error adding sensors: %v", err)))
		return
	}

	w.Write([]byte("Success"))
}
