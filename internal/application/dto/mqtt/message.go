package mqttdto

type Message struct {
	MessageType string      `json:"message_type"`
	Message     interface{} `json:"message"`
}

type StatusMessage struct {
	DatalogSerial string  `json:"datalogserial"`
	PVSerial      string  `json:"pvserial"`
	PVStatus      int     `json:"pvstatus"`
	PVPowerIn     float64 `json:"pvpowerin"`
	PV1Voltage    float64 `json:"pv1voltage"`
	PV1Current    float64 `json:"pv1current"`
	PV2Voltage    float64 `json:"pv2voltage"`
	PV2Current    float64 `json:"pv2current"`
	PVPowerOut    float64 `json:"pvpowerout"`
	ACFreq        float64 `json:"acfreq"`
	ACVoltage     float64 `json:"acvoltage"`
	ACOutputPower float64 `json:"acoutputpower"`
	Temperature   float64 `json:"temperature"`
	BatVoltage    float64 `json:"batvoltage"`
	BatCurrent    float64 `json:"batcurrent"`
	BatPower      float64 `json:"batpower"`
	GridExport    float64 `json:"gridexport"`
	GridImport    float64 `json:"gridimport"`
}

type HistoryMessage struct {
	DatalogSerial string  `json:"datalogserial"`
	PVSerial      string  `json:"pvserial"`
	Date          string  `json:"date"`
	EnergyToday   float64 `json:"energytoday"`
	EnergyTotal   float64 `json:"energytotal"`
}

type EventMessage struct {
	DatalogSerial string `json:"datalogserial"`
	PVSerial      string `json:"pvserial"`
	EventCode     string `json:"eventcode"`
	Description   string `json:"description"`
	Severity      string `json:"severity"`
}
