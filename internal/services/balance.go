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
)

type balanceService struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewBalanceService(client postgresql.Client, logger *logging.Logger) BalanceService {
	return &balanceService{
		client: client,
		logger: logger,
	}
}

func (r *balanceService) Delete(ctx context.Context, id int) error {
	q := `
		delete from balance
		where id=$1
		returning id;
`
	r.logger.Trace(fmt.Sprintf("SQL Quert: %s", internal.FormatQuery(q)))

	tag, err := r.client.Exec(ctx, q, id)
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

func (r *balanceService) Create(ctx context.Context, balance *models.Balance) error {
	q := `
			insert into balance 
			    (id, amount_on_balance, amount_on_reserve)
			values ($1, $2, $3);
	`

	r.logger.Trace(fmt.Sprintf("SQL Quert: %s", internal.FormatQuery(q)))

	tag, err := r.client.Exec(ctx, q, balance.Id, balance.AmountOnBalance, balance.AmountOnReserve)
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

func (r *balanceService) FindOne(ctx context.Context, id int) (*models.Balance, error) {
	if id < 1 {
		return nil, errors.New("incorrect id")
	}

	q := `
			select id, amount_on_balance, amount_on_reserve
			from balance
			where id=$1;
	`

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", internal.FormatQuery(q)))

	var blc models.Balance
	err := r.client.QueryRow(ctx, q, id).Scan(&blc.Id, &blc.AmountOnBalance, &blc.AmountOnReserve)
	if err != nil {
		return nil, err
	}

	return &blc, nil
}

func (r *balanceService) Update(ctx context.Context, balance *models.Balance) error {
	q := `
			update balance
			set amount_on_balance = $2
			where id = $1;
	`

	r.logger.Trace(fmt.Sprintf("SQL Quert: %s", internal.FormatQuery(q)))

	tag, err := r.client.Exec(ctx, q, balance.Id, balance.AmountOnBalance)
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

func (r *balanceService) UpdateFull(ctx context.Context, balance *models.Balance) error {
	q := `
			update balance
			set amount_on_balance = $2, amount_on_reserve = $3
			where id = $1;
	`

	r.logger.Trace(fmt.Sprintf("SQL Quert: %s", internal.FormatQuery(q)))

	tag, err := r.client.Exec(ctx, q, balance.Id, balance.AmountOnBalance, balance.AmountOnReserve)
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
