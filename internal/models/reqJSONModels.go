package models

type UnmarshallGetBalance struct {
	BalanceId string `json:"balance_id"`
}

type UnmarshallTopUpBalance struct {
	BalanceId string `json:"balance_id"`
	Amount    string `json:"amount"`
}

type UnmarshallGetTransactionHistory struct {
	BalanceId         string `json:"balance_id"`
	SortByDateOrder   string `json:"sort_by_date_order"`
	SortByAmountOrder string `json:"sort_by_amount_order"`
	PerPage           string `json:"per_page"`
	Page              string `json:"page"`
}

type UnmarshallServicePay struct {
	BalanceId string `json:"balance_id"`
	ServiceId string `json:"service_id"`
	OrderId   string `json:"order_id"`
	Amount    string `json:"amount"`
}

type UnmarshallCancel struct {
	BalanceId string `json:"balance_id"`
	OrderId   string `json:"order_id"`
	Amount    string `json:"amount"`
}

type UnmarshalGetReport struct {
	YearMonth string `json:"year_month"`
}
