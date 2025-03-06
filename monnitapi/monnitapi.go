package monnitapi

import (
	"log"
	"net/http"
	"time"
)

var requestChannel chan *http.Request
var responseChannel chan *http.Response
var httpClient *http.Client

func InitApiHandler(apiKeyId string, apiKeySecret string, timeOutSeconds int) {
	// Create non default http client
	httpClient = &http.Client{Timeout: time.Duration(timeOutSeconds) * time.Second}

	// Create API queue channel
	requestChannel = make(chan *http.Request, 1)
	responseChannel = make(chan *http.Response, 1)

	// Create request handler goroutine
	go func() {
		for query := range requestChannel {
			query.Header.Add("APIKeyID", apiKeyId)
			query.Header.Add("APISecretKey", apiKeySecret)
			response, err := httpClient.Do(query)
			if err != nil {
				log.Fatalln("Error processing request: ", err)
				responseChannel <- nil
			} else {
				responseChannel <- response
			}
		}
	}()
}

func CloseApiHandler() {
	// Close request queues
	close(requestChannel)
	close(responseChannel)
}

func GatewayReform(gatwayId int) bool {

	// resp, err := http.Get("")
	// if err != nil {

	// }
	return true
}
