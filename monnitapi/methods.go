package monnitapi

import (
	"fmt"
	"net/http"
)

type gatewayReformReq struct {
	GatewayID int `json:"GatewayID"`
}

func GatewayReform(gatewayId int) error {
	body := &gatewayReformReq{
		GatewayID: gatewayId,
	}

	response, err := apiCall[genericApiResponse](http.MethodPost, "/GatewayReform", body)

	if err != nil {
		return fmt.Errorf("could not reform gateway: %v", err)
	}

	if response.Result != "Success" {
		return fmt.Errorf("server did not send success.  Result: %s", response.Result)
	}
	return nil
}

type networkListReq struct {
	AccountID int `json:"accountID"`
}
type networkListResponse struct {
	Method string    `json:"Method"`
	Result []Network `json:"Result"`
}

func GetNetworkList() ([]Network, error) {
	body := &networkListReq{
		AccountID: accountID,
	}
	response, err := apiCall[networkListResponse](http.MethodPost, "/NetworkList", body)

	if err != nil {
		return nil, fmt.Errorf("could not fetch network list: %v", err)
	}

	return response.Result, nil
}

type sensorListReq struct {
	NetworkID int `json:"NetworkID"`
	AccountID int `json:"accountID"`
}

type sensorListResponse struct {
	Method string      `json:"Method"`
	Result []ApiSensor `json:"Result"`
}

func GetSensorsOnNetwork(networkId int) ([]ApiSensor, error) {
	body := &sensorListReq{
		NetworkID: networkId,
		AccountID: accountID,
	}
	response, err := apiCall[sensorListResponse](http.MethodPost, "/SensorList", body)

	if err != nil {
		return nil, fmt.Errorf("could not fetch network list: %v", err)
	}

	return response.Result, nil
}

type removeSensorReq struct {
	SensorID int `json:"SensorID"`
}

func RemoveSensor(sensorId int) error {
	body := &removeSensorReq{
		SensorID: sensorId,
	}
	response, err := apiCall[genericApiResponse](http.MethodPost, "/RemoveSensor", body)

	if err != nil {
		return fmt.Errorf("could not remove sensor: %v", err)
	}

	if response.Result != "Success" {
		return fmt.Errorf("server did not send success.  Result: \"%s\"", response.Result)
	}
	return nil
}

type gatewayListResponse struct {
	Method string    `json:"Method"`
	Result []Gateway `json:"Result"`
}

func GetGatewaysOnNetwork(networkId int) ([]Gateway, error) {
	body := &sensorListReq{
		NetworkID: networkId,
		AccountID: accountID,
	}
	response, err := apiCall[gatewayListResponse](http.MethodPost, "/GatewayList", body)

	if err != nil {
		return nil, fmt.Errorf("could not fetch gateway list: %v", err)
	}

	return response.Result, nil
}

func ReformNetwork(networkId int) error {
	gateways, err := GetGatewaysOnNetwork(networkId)
	if err != nil {
		return err
	}
	for _, g := range gateways {
		err = GatewayReform(g.GatewayID)
		if err != nil {
			return err
		}
	}
	return nil
}

type updateHeartbeatReq struct {
	SensorID            int     `json:"SensorID"`
	ReportInterval      float64 `json:"ReportInterval"`
	ActiveStateInterval float64 `json:"ActiveStateInterval"`
}

func SetHeartbeat(sensorId int, reportInterval float64, activeStateInterval float64) error {
	body := &updateHeartbeatReq{
		SensorID:            sensorId,
		ReportInterval:      reportInterval,
		ActiveStateInterval: activeStateInterval,
	}

	response, err := apiCall[genericApiResponse](http.MethodPost, "/SensorSetHeartbeat", body)

	if err != nil {
		return fmt.Errorf("could not update hearbeat: %v", err)
	}

	if response.Result != "Success" {
		return fmt.Errorf("could not update hearbeat: %v", err)
	}

	return nil
}
