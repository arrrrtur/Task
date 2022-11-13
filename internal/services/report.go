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
	"strconv"
)

type reportService struct {
	client postgresql.Client
	logger *logging.Logger
}

func NewReportService(client postgresql.Client, logger *logging.Logger) ReportService {
	return &reportService{
		client: client,
		logger: logger,
	}
}

func (r *reportService) Create(ctx context.Context, report *models.Report) error {
	q := `
			insert into report 
			    (balance_id, service_id, amount)
			values ($1, $2, $3);
	`

	r.logger.Trace(fmt.Sprintf("SQL Quert: %s", internal.FormatQuery(q)))

	tag, err := r.client.Exec(ctx, q, report.BalanceId, report.ServiceId, report.Amount)
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

func (r *reportService) FindAll(ctx context.Context, year, month int) ([][]string, error) {
	q := `
			select 
			    s.id,
			    s.name,
			    sum(amount)
			from report
			inner join service s on s.id = report.service_id
			where (extract(year  from report_time)) = $1
			and (extract(month from report_time)) = $2
			group by s.id;
	`

	r.logger.Trace(fmt.Sprintf("SQL Quert: %s", internal.FormatQuery(q)))

	rep := make([][]string, 0)
	rep = append(rep, []string{"Id услуги", "Название услуги", fmt.Sprintf("Прибыль за %d-%d", year, month)})
	rows, err := r.client.Query(ctx, q, year, month)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			serviceId   int
			serviceName string
			income      float64
		)
		if err := rows.Scan(&serviceId, &serviceName, &income); err != nil {
			return nil, err
		}

		str := []string{strconv.Itoa(serviceId), serviceName, strconv.FormatFloat(income, 'f', 2, 64)}

		rep = append(rep, str)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return rep, nil
}
