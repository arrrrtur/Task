package models

type Balance struct {
	Id              int     `json:"balance_id,"`
	AmountOnBalance float64 `json:"amount_on_balance"`
	AmountOnReserve float64 `json:"amount_on_reserve"`
}
