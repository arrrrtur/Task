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

type avitoServiceServiceService struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewAvitoServiceService(client postgresql.Client, logger *logging.Logger) AvitoServiceService {
	return &avitoServiceServiceService{
		client: client,
		logger: logger,
	}
}

func (r *avitoServiceServiceService) Create(ctx context.Context, srv *models.AvitoService) error {
	q := `
			insert into service 
			    (id, name, price )
			values ($1, $2, $3);
	`

	r.logger.Trace(fmt.Sprintf("SQL Quert: %s", internal.FormatQuery(q)))

	tag, err := r.client.Exec(ctx, q, srv.Id, srv.Name, srv.Price)
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

func (r *avitoServiceServiceService) FindOne(ctx context.Context, id int) (*models.AvitoService, error) {
	if id < 1 {
		return nil, errors.New("incorrect id")
	}

	q := `
			select id, name, price
			from service
			where id=$1;
	`

	r.logger.Trace(fmt.Sprintf("SQL Quert: %s", internal.FormatQuery(q)))

	var srv models.AvitoService
	err := r.client.QueryRow(ctx, q, id).Scan(&srv.Id, &srv.Name, &srv.Price)
	if err != nil {
		return nil, err
	}

	return &srv, nil
}
