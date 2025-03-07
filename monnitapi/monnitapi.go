package monnitapi

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

var requestQueue chan *http.Request
var resultChannel chan *http.Response
var httpClient *http.Client
var baseUrl string

// Initializes variables required for API communcation
func InitApiHandler(url string, apiKeyId string, apiKeySecret string) {
	baseUrl = url
	// Create non default http client
	httpClient = &http.Client{Timeout: 10 * time.Second}

	// Create API queue channel
	requestQueue = make(chan *http.Request, 1)
	resultChannel = make(chan *http.Response, 1)

	//
	go func() {
		// Watch the request queue
		for query := range requestQueue {
			// Add api auth
			query.Header.Add("APIKeyID", apiKeyId)
			query.Header.Add("APISecretKey", apiKeySecret)

			// Set content type
			query.Header.Set("Content-Type", "application/json")

			// Run the request
			response, err := httpClient.Do(query)

			// Return result to channel
			if err != nil {
				log.Fatalln("Error processing request: ", err)
				resultChannel <- nil
			} else {
				resultChannel <- response
			}
		}
	}()
}

// Closes request handler channels
func CloseApiHandler() {
	close(requestQueue)
	close(resultChannel)
}

type genericApiResponse struct {
	Method string `json:"Method"`
	Result string `json:"Result"`
}

type gatewayReformReq struct {
	GatewayID int `json:"GatewayID"`
}

func GatewayReform(gatewayId int) bool {
	// Create request body
	body := &gatewayReformReq{
		GatewayID: gatewayId,
	}

	// Create the payload buffer and copy to it
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(body)

	// Create the request
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/GatewayReform", payloadBuffer)

	// Check for errors
	if err != nil {
		log.Fatalf("Could not create gateway reform post request: %v", err)
		return false
	}
	requestQueue <- req
	response := <-resultChannel

	if response == nil {
		log.Fatalln("Gateway reform response is null")
		return false
	}

	// Close the body when done
	defer response.Body.Close()

	// Create the response object
	decodedRes := new(genericApiResponse)

	// Decode buffer to object
	err = json.NewDecoder(response.Body).Decode(decodedRes)

	// Check for errors
	if err != nil {
		log.Fatalf("Error decoding gateway reform response: %v", err)
		return false
	}

	return decodedRes.Result == "Success"
}

type sensorListReq struct {
	NetworkID int `json:"NetworkID"`
	AccountID int `json:"accountID"`
}

type SensorListResponse struct {
	Method string `json:"Method"`
	Result []ApiSensor `json:"Result"`
}

type ApiSensor struct {
	SensorID             int    `json:"SensorID"`
	ApplicationID        int    `json:"ApplicationID"`
	CSNetID              int    `json:"CSNetID"`
	SensorName           string `json:"SensorName"`
	LastCommuncationDate string `json:"LastCommunicationDate"`
	NextCommuncationDate string `json:"NextCommunicationDate"`
	LastDataMessageID    int    `json:"LastDataMessageID"`
	PowerSourceID        int    `json:"PowerSourceID"`
	Status               int    `json:"Status"`
	CanUpdate            bool   `json:"CanUpdate"`
	CurrentReading       string `json:"CurrentReading"`
	BatteryLevel         int    `json:"BatteryLevel"`
	SignalStrength       int    `json:"SignalStrength"`
	AlertsActive         bool   `json:"AlertsActive"`
	CheckDigit           string `json:"CheckDigit"`
	AccountID            int    `json:"AccountID"`
	MonnitApplicationID  int    `json:"MonnitApplicationID"`
}

func GetSensorsOnNetwork(networkId int) *SensorListResponse {
	// Create request body
	body := &sensorListReq{
		NetworkID: networkId,
		AccountID: 73221,
	}

	// Create the payload buffer and copy to it
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(body)

	// Create the request
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/SensorList", payloadBuffer)

	// Check for errors
	if err != nil {
		log.Fatalf("Could not create gateway reform post request: %v", err)
		return nil
	}
	requestQueue <- req
	response := <-resultChannel

	if response == nil {
		log.Fatalln("Gateway reform response is null")
		return nil
	}

	// Close the body when done
	defer response.Body.Close()

	// Create the response object
	decodedRes := new(SensorListResponse)

	res := response.Body

	// Decode buffer to object
	err = json.NewDecoder(res).Decode(decodedRes)

	// Check for errors
	if err != nil {
		log.Fatalf("Error decoding gateway reform response: %v", err)
		return nil
	}

	// Check if reform is successful
	return decodedRes
}


type removeSensorReq struct {
	SensorID int `json:"SensorID"`
}

func RemoveSensor(sensorId int) bool {
	// Create request body
	body := &removeSensorReq{
		SensorID: sensorId,
	}

	// Create the payload buffer and copy to it
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(body)

	// Create the request
	req, err := http.NewRequest(http.MethodPost, baseUrl+"/RemoveSensor", payloadBuffer)

	// Check for errors
	if err != nil {
		log.Fatalf("Could not create remove sensor post request: %v", err)
		return false
	}
	requestQueue <- req
	response := <-resultChannel

	if response == nil {
		log.Fatalln("Remove sensor response is null")
		return false
	}

	// Close the body when done
	defer response.Body.Close()

	// Create the response object
	decodedRes := new(genericApiResponse)

	// Decode buffer to object
	err = json.NewDecoder(response.Body).Decode(decodedRes)

	// Check for errors
	if err != nil {
		log.Fatalf("Error decoding sensor remove response: %v", err)
		return false
	}

	return decodedRes.Result == "Success"
}
