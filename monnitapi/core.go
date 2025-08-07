package monnitapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var requestQueue chan *http.Request
var resultChannel chan *http.Response
var httpClient *http.Client
var baseUrl string
var accountID int

var debugLog *log.Logger
var errorLog *log.Logger

// Initializes variables required for API communcation
func Init(url string, apiKeyId string, apiKeySecret string, accountId int) {

	debugLog = log.New(os.Stdout, "- [monnitapi][DEBUG]: ", log.Ldate|log.Ltime|log.Lmsgprefix)
	errorLog = log.New(os.Stdout, "- [monnitapi][ERROR]: ", log.Ldate|log.Ltime|log.Lmsgprefix)

	baseUrl = url
	accountID = accountId
	// Create non default http client
	httpClient = &http.Client{Timeout: 10 * time.Second}

	// Create API queue channel
	requestQueue = make(chan *http.Request, 1)
	resultChannel = make(chan *http.Response, 1)

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

func Close() {
	close(requestQueue)
	close(resultChannel)
}

func apiCall[T any](method string, route string, body any) (*T, error) {
	debugLog.Printf("Calling %s", route)
	payloadBuffer := new(bytes.Buffer)
	json.NewEncoder(payloadBuffer).Encode(body)

	req, err := http.NewRequest(method, baseUrl+route, payloadBuffer)

	if err != nil {
		return nil, fmt.Errorf("could not create http request: %v", err)
	}

	requestQueue <- req
	response := <-resultChannel

	if response == nil {
		return nil, errors.New("api response is null")
	}

	defer response.Body.Close()

	decodedRes := new(T)

	err = json.NewDecoder(response.Body).Decode(decodedRes)

	if err != nil {
		return nil, fmt.Errorf("error decoding api response: %v", err)
	}

	return decodedRes, nil
}
