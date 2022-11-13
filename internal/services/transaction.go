package services

import (
	"Task/internal"
	"Task/internal/models"
	"Task/pkg/api"
	"Task/pkg/client/postgresql"
	"Task/pkg/logging"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"time"
)

type transactionService struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewTransactionService(client postgresql.Client, logger *logging.Logger) TransactionService {
	return &transactionService{
		client: client,
		logger: logger,
	}
}

func (r *transactionService) Create(ctx context.Context, transaction *models.Transaction) error {
	q := `
			insert into transaction 
			    (sender_id, receiver_id, transaction_price, operation_id, status)
			values ($1, $2, $3, $4, $5);
	`

	r.logger.Trace(fmt.Sprintf("SQL Quert: %s", internal.FormatQuery(q)))

	tag, err := r.client.Exec(ctx, q, transaction.SenderId, transaction.ReceiverId,
		transaction.TransactionPrice, transaction.OperationId, transaction.Status)

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

func (r *transactionService) Update(ctx context.Context, transaction *models.Transaction) error {
	q := `
			update transaction
			set status = $3
			where sender_id = $1 and receiver_id = $2 and status = 2;
	`

	r.logger.Trace(fmt.Sprintf("SQL Quert: %s", internal.FormatQuery(q)))

	tag, err := r.client.Exec(ctx, q, transaction.SenderId, transaction.ReceiverId, transaction.Status)
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

func (r *transactionService) FindAll(ctx context.Context, balanceId int, params api.Params) ([]models.TransactionDTO, error) {
	q := `
			select 
			    t.id, 
			    t.sender_id, 
			    t.receiver_id, 
			    o.name,
			    t.transaction_time, 
			    t.transaction_price,
			    t.status
			from transaction t
			inner join operation_type o 
				on t.operation_id = o.id
			where t.sender_id = $1 OR t.receiver_id = $1
	`

	q = fmt.Sprintf("%s order by transaction_time %s, transaction_price %s limit %d offset %d;",
		q, params.SortByDateOrder, params.SortByAmountOrder, params.PerPage, (params.Page-1)*params.PerPage)

	r.logger.Trace(fmt.Sprintf("SQL Query: %s", internal.FormatQuery(q)))

	rows, err := r.client.Query(ctx, q, balanceId)
	if err != nil {
		return nil, err
	}

	transactions := make([]models.TransactionDTO, 0)

	for rows.Next() {
		var transaction models.TransactionDTO
		var tmpstm time.Time
		if err := rows.Scan(&transaction.Id, &transaction.SenderId, &transaction.ReceiverId,
			&transaction.OperationType, &tmpstm, &transaction.TransactionPrice,
			&transaction.Status); err != nil {
			return nil, err
		}
		transaction.TransactionTime = tmpstm.String()
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}
