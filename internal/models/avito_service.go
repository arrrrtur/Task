package models

type AvitoService struct {
	Id    int     `json:"service_id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}
