package monnitapi

type WebhookMessage struct {
	GatewayMessage GatewayMessage  `json:"gatewayMessage"`
	SensorMessages []SensorMessage `json:"sensorMessages,omitempty"`
}

type GatewayMessage struct {
	GatewayID      string `json:"gatewayID"`
	GatewayName    string `json:"gatewayName"`
	AccountID      string `json:"accountID"`
	NetworkID      string `json:"networkID"`
	MessageType    string `json:"messageType"`
	Power          string `json:"power"`
	BatteryLevel   string `json:"batteryLevel"`
	Date           string `json:"date"`
	Count          string `json:"count"`
	SignalStrength string `json:"signalStrength"`
	PendingChange  string `json:"pendingChange"`
}

type SensorMessage struct {
	SensorID        string `json:"sensorID"`
	SensorName      string `json:"sensorName"`
	ApplicationID   string `json:"applicationID"`
	NetworkID       string `json:"networkID"`
	DataMessageGUID string `json:"dataMessageGUID"`
	State           string `json:"state"`
	RawData         string `json:"rawData"`
	DataType        string `json:"dataType"`
	DataValue       string `json:"dataValue"`
	PlotValues      string `json:"plotValues"`
	PlotLabels      string `json:"plotLabels"`
	BatteryLevel    string `json:"batteryLevel"`
	SignalStrength  string `json:"signalStrength"`
	PendingChange   string `json:"pendingChange"`
	Voltage         string `json:"voltage"`
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

type genericApiResponse struct {
	Method string `json:"Method"`
	Result string `json:"Result"`
}

type Network struct {
	NetworkID           int    `json:"NetworkID"`
	NetworkName         string `json:"NetworkName"`
	SendNotifications   bool   `json:"SendNotifications"`
	ExternalAccessUntil string `json:"ExternalAccessUntil"`
}

type Gateway struct {
	GatewayID            int     `json:"GatewayID"`
	NetworkID            int     `json:"NetworkID"`
	Name                 string  `json:"Name"`
	GatewayType          string  `json:"GatewayType"`
	Heartbeat            float64 `json:"Heartbeat"`
	IsDirty              bool    `json:"IsDirty"`
	LastCommuncationDate string  `json:"LastCommunicationDate"`
	LastInboundIPAddress string  `json:"LastInboundIPAddress"`
	MacAddress           string  `json:"MacAddress"`
	IsUnlocked           bool    `json:"IsUnlocked"`
	CheckDigit           string  `json:"CheckDigit"`
	AccountID            int     `json:"AccountID"`
	SignalStrength       int     `json:"SignalStrength"`
	BatteryLevel         int     `json:"BatteryLevel"`
	ResetInterval        int     `json:"ResetInterval"`
	GatewayPowerMode     string  `json:"GatewayPowerMode"`
}
