package monnitapi

type WebhookMessage struct {
	gatewayMessage GatewayMessage
	sensorMessages []SensorMessage
}

type GatewayMessage struct {
	gatewayID      int
	gatewayName    string
	accountID      int
	networkID      int
	messageType    int
	power          int
	batteryLevel   int
	date           string
	count          int
	signalStrength int
	pendingChange  bool
}

type SensorMessage struct {
	sensorID        int
	sensorName      string
	applicationID   int
	networkID       int
	dataMessageGUID string
	state           int
	rawData         string
	dataType        string
	dataValue       string
	plotValues      string
	plotLabels      string
	batteryLevel    int
	signalStrength  int
	pendingChange   bool
	voltage         float32
}
