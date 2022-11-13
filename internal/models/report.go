package models

type Report struct {
	ID         int     `json:"report_id"`
	BalanceId  int     `json:"balance_id"`
	ServiceId  int     `json:"service_id"`
	Amount     float64 `json:"amount"`
	ReportTime string  `json:"report_time"`
}
