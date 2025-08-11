package monitoringdto

import "time"

type CustomerPanelStatusResponse struct {
	DatalogSerial string    `json:"datalog_serial"`
	PVSerial      string    `json:"pv_serial"`
	PVStatus      int       `json:"pv_status"`
	PVPowerIn     float64   `json:"pv_power_in"`
	PV1Voltage    float64   `json:"pv1_voltage"`
	PV1Current    float64   `json:"pv1_current"`
	PV2Voltage    float64   `json:"pv2_voltage"`
	PV2Current    float64   `json:"pv2_current"`
	PVPowerOut    float64   `json:"pv_power_out"`
	ACFreq        float64   `json:"ac_freq"`
	ACVoltage     float64   `json:"ac_voltage"`
	ACOutputPower float64   `json:"ac_output_power"`
	Temperature   float64   `json:"temperature"`
	BatVoltage    float64   `json:"bat_voltage"`
	BatCurrent    float64   `json:"bat_current"`
	BatPower      float64   `json:"bat_power"`
	GridExport    float64   `json:"grid_export"`
	GridImport    float64   `json:"grid_import"`
	EnergyToday   float64   `json:"energy_today"`
	EnergyTotal   float64   `json:"energy_total"`
	Timestamp     time.Time `json:"timestamp"`
}
