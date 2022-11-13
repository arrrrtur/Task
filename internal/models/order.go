package models

type Order struct {
	Id          int    `json:"order_id"`
	BalanceId   int    `json:"balance_Id"`
	ServiceId   int    `json:"service_Id"`
	ReserveTime string `json:"reserve_time"`
}
