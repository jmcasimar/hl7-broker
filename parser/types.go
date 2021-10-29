package parser

import "time"

type Patient struct {
	FirstName string
	LastName  string
	MRN       string
	AdmitDate time.Time
}

type PeriodicParam struct {
	ParameterID string `json:"parameter_id"`
	Parameter   string `json:"parameter"`
	Value       string `json:"value"`
}

type AperiodicParam struct {
	ParameterID string    `json:"parameter_id"`
	Parameter   string    `json:"parameter"`
	Value       string    `json:"value"`
	Time        time.Time `json:"time"`
}

type PhysiologicalAlarm struct {
	Level   string    `json:"level"`
	AlarmID string    `json:"id"`
	Message string    `json:"message"`
	Time    time.Time `json:"time"`
}

type TechnicalAlarm struct {
	Level   string `json:"level"`
	AlarmID string `json:"id"`
	Message string `json:"message"`
}

type AlarmLimit struct {
	Type      string `json:"type"`
	Parameter string `json:"parameter"`
	Value     string `json:"value"`
}

type Record struct {
	MRN                 string                `json:"mrn"`
	MonitorID           string                `json:"monitorId"`
	Time                time.Time             `json:"time"`
	PeriodicParams      []*PeriodicParam      `json:"periodic_params,omitempty"`
	AperiodicParams     []*AperiodicParam     `json:"aperiodic_params,omitempty"`
	PhysiologicalAlarms []*PhysiologicalAlarm `json:"physiological_alarms,omitempty"`
	TechnicalAlarms     []*TechnicalAlarm     `json:"technical_alarms,omitempty"`
	AlarmLimits         []*AlarmLimit         `json:"alarm_limits,omitempty"`
}
