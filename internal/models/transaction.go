package models

type Transaction struct {
	Id               int     `json:"transaction_id"`
	SenderId         int     `json:"sender_id"`
	ReceiverId       int     `json:"receiver_id"`
	TransactionTime  string  `json:"transaction_time"`
	TransactionPrice float64 `json:"transaction_price"`
	OperationId      int     `json:"operation_id"`
	Status           int     `json:"status"`
}

type TransactionDTO struct {
	Id               int     `json:"transaction_id"`
	SenderId         int     `json:"sender_id"`
	ReceiverId       int     `json:"receiver_id"`
	TransactionTime  string  `json:"transaction_time"`
	TransactionPrice float64 `json:"transaction_price"`
	OperationType    string  `json:"operation_type"`
	Status           int     `json:"status"`
}
