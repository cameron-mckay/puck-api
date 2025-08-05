package monnitapi

type WebhookMessage struct {
	GatewayMessage GatewayMessage  `json:"gatewayMessage"`
	SensorMessages []SensorMessage `json:"sensorMessage"`
}

type GatewayMessage struct {
	GatewayID      int    `json:"gatewayID"`
	GatewayName    string `json:"gatewayName"`
	AccountID      int    `json:"accountID"`
	NetworkID      int    `json:"networkID"`
	MessageType    int    `json:"messageType"`
	Power          int    `json:"power"`
	BatteryLevel   int    `json:"batteryLevel"`
	Date           string `json:"date"`
	Count          int    `json:"count"`
	SignalStrength int    `json:"signalStrength"`
	PendingChange  bool   `json:"pendingChange"`
}

type SensorMessage struct {
	SensorID        int     `json:"sensorID"`
	SensorName      string  `json:"sensorName"`
	ApplicationID   int     `json:"applicationID"`
	NetworkID       int     `json:"networkID"`
	DataMessageGUID string  `json:"dataMessageGUID"`
	State           int     `json:"state"`
	RawData         string  `json:"rawData"`
	DataType        string  `json:"dataType"`
	DataValue       string  `json:"dataValue"`
	PlotValues      string  `json:"plotValues"`
	PlotLabels      string  `json:"plotLabels"`
	BatteryLevel    int     `json:"batteryLevel"`
	SignalStrength  int     `json:"signalStrength"`
	PendingChange   bool    `json:"pendingChange"`
	Voltage         float32 `json:"voltage"`
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
