package services

import (
	"Task/internal"
	"Task/internal/models"
	"Task/pkg/client/postgresql"
	"Task/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
)

type orderService struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewOrderService(client postgresql.Client, logger *logging.Logger) OrderService {
	return &orderService{
		client: client,
		logger: logger,
	}
}

func (r *orderService) Create(ctx context.Context, order *models.Order) error {
	q := `
			insert into "order"
			    (id, balance_id, service_id)
			values ($1, $2, $3);
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", internal.FormatQuery(q)))

	tag, err := r.client.Exec(ctx, q, order.Id, order.BalanceId, order.ServiceId)
	if err != nil || tag.RowsAffected() == 0 {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail: %s, Where: %s, Code: %s SQLState: %s",
				pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			r.logger.Error(newErr)
			return newErr
		}
		return err
	}
	return nil
}

func (r *orderService) FindOne(ctx context.Context, id int) (*models.Order, error) {
	if id < 1 {
		return nil, errors.New("incorrect id")
	}

	q := `
			select id, balance_id, service_id, reserve_time
			from "order"
			where id=$1;
	`

	r.logger.Trace(fmt.Sprintf("SQL Quert: %s", internal.FormatQuery(q)))

	var ordr models.Order
	var tmpstmp pgtype.Timestamp

	err := r.client.QueryRow(ctx, q, id).Scan(&ordr.Id, &ordr.BalanceId, &ordr.ServiceId, &tmpstmp)
	if err != nil {
		return nil, err
	}
	ordr.ReserveTime = tmpstmp.Time.String()

	return &ordr, nil
}
