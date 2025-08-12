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

var routeInfo *log.Logger
var routeDebug *log.Logger
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
	routeInfo = log.New(os.Stdout, "- [routes][INFO]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
	routeDebug = log.New(os.Stdout, "- [routes][DEBUG]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
	routeError = log.New(os.Stderr, "- [routes][ERROR]: ", log.Ldate|log.Ltime|log.Lmsgprefix|log.Lshortfile)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", websocket.Serve)
	http.HandleFunc("/api/getNetworks", getNetworks)
	http.HandleFunc("/api/deleteAllSensorsOnNetwork", deleteAllSensorsOnNetwork)
	http.HandleFunc("/api/getSensorsOnNetwork", getSensorsOnNetwork)
	http.HandleFunc("/api/getMessageCountsOnNetwork", getMessageCountsOnNetwork)
	http.HandleFunc("/api/addBinToNetwork", addBinToNetwork)
	http.HandleFunc("/api/setHeartbeat", updateHeartbeat)

	go http.ListenAndServe(addr, nil)

	routeInfo.Println("Listening on " + addr)
}

func Close() {
	// nothing here yet
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/index.html")
	routeInfo.Println("Served index")
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
	routeInfo.Println("Sensor list requested")
	w.Header().Set("Content-Type", "application/json")

	decodedReq := new(networkReq)

	err := json.NewDecoder(r.Body).Decode(decodedReq)

	if err != nil {
		routeError.Println(err)
		w.Write([]byte(fmt.Sprintf("Error decoding request: %v", err)))
		return
	}

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

func getMessageCountsOnNetwork(w http.ResponseWriter, r *http.Request) {
	routeInfo.Println("Message counts requested")
	w.Header().Set("Content-Type", "application/json")

	decodedReq, err := reqBodyParser[networkReq](r)

	if err != nil {
		routeError.Println(err)
		w.Write([]byte(fmt.Sprintf("Error decoding request: %v", err)))
		return
	}

	mcs, err := db.GetMessageCounts(decodedReq.NetworkID)

	if err != nil {
		routeError.Println(err)
		w.Write([]byte(fmt.Sprintf("Error getting message counts: %v", err)))
		return
	}

	res, err := json.Marshal(mcs)

	if err != nil {
		routeError.Println(err)
		w.Write([]byte(fmt.Sprintf("Error encoding message count json: %v", err)))
		return
	}

	w.Write(res)
}

func deleteAllSensorsOnNetwork(w http.ResponseWriter, r *http.Request) {
	routeInfo.Println("Sensor delete requested")
	w.Header().Set("Content-Type", "application/json")

	decodedReq, err := reqBodyParser[networkReq](r)

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

	err = monnitapi.ReformNetwork(decodedReq.NetworkID)
	if err != nil {
		routeError.Printf("can't reform network, error: %v", err)
		w.Write([]byte(fmt.Sprintf("Error reforming network: %v", err)))
		return
	}

	w.Write([]byte("Success"))
}

type binReq struct {
	BinID     int `json:"binId"`
	NetworkID int `json:"networkId"`
}

func addBinToNetwork(w http.ResponseWriter, r *http.Request) {
	routeInfo.Println("Sensor add requested")
	w.Header().Set("Content-Type", "application/json")

	decodedReq, err := reqBodyParser[binReq](r)

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

	err = monnitapi.ReformNetwork(decodedReq.NetworkID)
	if err != nil {
		routeError.Printf("can't reform network, error: %v", err)
		w.Write([]byte(fmt.Sprintf("Error reforming network: %v", err)))
		return
	}

	w.Write([]byte("Success"))
}

type sensorReq struct {
	SensorID int `json:"sensorId"`
}

func updateHeartbeat(w http.ResponseWriter, r *http.Request) {
	routeInfo.Println("Sensor add requested")
	w.Header().Set("Content-Type", "application/json")

	decodedReq, err := reqBodyParser[sensorReq](r)

	if err != nil {
		routeError.Println(err)
		http.Error(w, fmt.Sprintf("Error decoding request: %v", err), 400)
		return
	}

	err = monnitapi.SetHeartbeat(decodedReq.SensorID, 120.0001, 120.0001)

	if err != nil {
		http.Error(w, fmt.Sprintf("Error setting heartbeat: %v", err), 500)
		return
	}

	w.Write([]byte("Success"))
}
