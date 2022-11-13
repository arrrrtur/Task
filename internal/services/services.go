package services

import (
	"Task/internal/models"
	"Task/pkg/api"
	"Task/pkg/client/postgresql"
	"Task/pkg/logging"
	"context"
)

type AvitoServiceService interface {
	Create(ctx context.Context, srv *models.AvitoService) error
	FindOne(ctx context.Context, id int) (*models.AvitoService, error)
}

type BalanceService interface {
	Create(ctx context.Context, balance *models.Balance) error
	FindOne(ctx context.Context, id int) (*models.Balance, error)
	Update(ctx context.Context, balance *models.Balance) error
	UpdateFull(ctx context.Context, balance *models.Balance) error
	Delete(ctx context.Context, id int) error
}

type OrderService interface {
	Create(ctx context.Context, order *models.Order) error
	FindOne(ctx context.Context, id int) (*models.Order, error)
}

type ReportService interface {
	Create(ctx context.Context, report *models.Report) error
	FindAll(ctx context.Context, year, month int) ([][]string, error)
}

type TransactionService interface {
	Create(ctx context.Context, transaction *models.Transaction) error
	Update(ctx context.Context, transaction *models.Transaction) error
	FindAll(ctx context.Context, balanceId int, params api.Params) ([]models.TransactionDTO, error)
}

type Services struct {
	BalanceService      BalanceService
	TransactionService  TransactionService
	AvitoServiceService AvitoServiceService
	OrderService        OrderService
	ReportService       ReportService
}

func NewServices(client postgresql.Client, logger *logging.Logger) *Services {
	return &Services{
		BalanceService:      NewBalanceService(client, logger),
		TransactionService:  NewTransactionService(client, logger),
		AvitoServiceService: NewAvitoServiceService(client, logger),
		OrderService:        NewOrderService(client, logger),
		ReportService:       NewReportService(client, logger),
	}
}
