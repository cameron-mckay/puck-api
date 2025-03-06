package monnitapi

import (
	"log"
	"net/http"
	"time"
)

var requestQueue chan *http.Request
var resultChannel chan *http.Response
var httpClient *http.Client

func InitApiHandler(apiKeyId string, apiKeySecret string) {
	// Create non default http client
	httpClient = &http.Client{Timeout: 10 * time.Second}

	// Create API queue channel
	requestQueue = make(chan *http.Request, 1)
	resultChannel = make(chan *http.Response, 1)

	//
	go func() {
		for query := range requestQueue {
			query.Header.Add("APIKeyID", apiKeyId)
			query.Header.Add("APISecretKey", apiKeySecret)
			response, err := httpClient.Do(query)
			if err != nil {
				log.Fatalln("Error processing request: ", err)
				resultChannel <- nil
			} else {
				resultChannel <- response
			}
		}
	}()
}

func CloseApiHandler() {
	close(requestQueue)
	close(resultChannel)
}

func GatewayReform(gatwayId int) bool {

	// resp, err := http.Get("")
	// if err != nil {

	// }
	return true
}
