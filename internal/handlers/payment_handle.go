package handlers

import (
	"Task/internal"
	"Task/internal/apperror"
	"Task/internal/models"
	"net/http"
	"strconv"
)

const (
	REPLENISHMENT = iota + 1
	TRANSFER
	PAY
)

// В случае возникновения ошибок, возвращает баланс в исходное состояние
func backUpBalance(app *application, r *http.Request, blc *models.Balance, amount float64) error {
	blc.AmountOnBalance += amount
	blc.AmountOnReserve -= amount
	if err := app.services.BalanceService.UpdateFull(r.Context(), blc); err != nil {
		return err
	}
	return nil
}

// ReserveFromBalance Перечислить деньги в резерв
func (app *application) ReserveFromBalance(w http.ResponseWriter, r *http.Request) error {
	jsn, err := internal.ParsJSON(r)
	if err != nil {
		return err
	}

	balanceId, err := strconv.Atoi(jsn["balance_id"])
	if err != nil {
		return err
	}

	serviceId, err := strconv.Atoi(jsn["service_id"])
	if err != nil {
		return err
	}

	orderId, err := strconv.Atoi(jsn["order_id"])
	if err != nil {
		return err
	}

	amount, err := strconv.ParseFloat(jsn["amount"], 64)
	if err != nil {
		return err
	}

	blc, err := app.services.BalanceService.FindOne(r.Context(), balanceId)
	if blc == nil {
		return apperror.NewAppError(err, "Bad request", "id not found", "some-code")
	}

	if amount > blc.AmountOnBalance {
		return apperror.NewAppError(err, "Bad request", "not enough money", "some-code")
	}

	blc.AmountOnBalance -= amount
	blc.AmountOnReserve += amount
	if err := app.services.BalanceService.UpdateFull(r.Context(), blc); err != nil {
		return err
	}

	avsrvs, err := app.services.AvitoServiceService.FindOne(r.Context(), serviceId)
	if err != nil || avsrvs == nil {
		avsrvs = &models.AvitoService{
			Id:    serviceId,
			Name:  "Какая-то услуга",
			Price: amount,
		}
		if err := app.services.AvitoServiceService.Create(r.Context(), avsrvs); err != nil {
			if err := backUpBalance(app, r, blc, amount); err != nil {
				return err
			}
			return err
		}
	}

	ordr, err := app.services.OrderService.FindOne(r.Context(), orderId)
	if ordr != nil {
		if err := backUpBalance(app, r, blc, amount); err != nil {
			return err
		}
		return apperror.NewAppError(err, "bad request", "such order already exist", "some code")
	}

	ordr = &models.Order{
		Id:        orderId,
		BalanceId: balanceId,
		ServiceId: serviceId,
	}
	trn := &models.Transaction{
		SenderId:         balanceId,
		ReceiverId:       balanceId,
		TransactionPrice: amount,
		OperationId:      PAY,
		Status:           IN_PROCESSING,
	}

	if err := app.services.OrderService.Create(r.Context(), ordr); err != nil {
		if err := backUpBalance(app, r, blc, amount); err != nil {
			return err
		}
		return err
	}
	if err := app.services.TransactionService.Create(r.Context(), trn); err != nil {
		if err := backUpBalance(app, r, blc, amount); err != nil {
			return err
		}
		return err
	}

	if err = internal.WriteJSON(w, trn, ordr, avsrvs, blc); err != nil {
		return err
	}

	return nil
}

// WriteOffMoney Списать деньги с резерва
func (app *application) WriteOffMoney(w http.ResponseWriter, r *http.Request) error {
	jsn, err := internal.ParsJSON(r)
	if err != nil {
		return err
	}

	balanceId, err := strconv.Atoi(jsn["balance_id"])
	if err != nil {
		return err
	}

	serviceId, err := strconv.Atoi(jsn["service_id"])
	if err != nil {
		return err
	}

	orderId, err := strconv.Atoi(jsn["order_id"])
	if err != nil {
		return err
	}

	amount, err := strconv.ParseFloat(jsn["amount"], 64)
	if err != nil {
		return err
	}

	blc, err := app.services.BalanceService.FindOne(r.Context(), balanceId)
	if blc == nil {
		return apperror.NewAppError(err, "Bad request", "balance not found", "some-code")
	}

	if amount > blc.AmountOnReserve {
		return apperror.NewAppError(err, "Bad request", "not enough reserved money", "some-code")
	}

	_, err = app.services.OrderService.FindOne(r.Context(), orderId)
	if err != nil {
		return apperror.NewAppError(err, "Bad request", "order not found", "some-code")
	}

	_, err = app.services.AvitoServiceService.FindOne(r.Context(), serviceId)
	if err != nil {
		return apperror.NewAppError(err, "Bad request", "service not found", "some-code")
	}

	blc.AmountOnReserve -= amount
	if err := app.services.BalanceService.UpdateFull(r.Context(), blc); err != nil {
		return err
	}

	trn := &models.Transaction{
		SenderId:         balanceId,
		ReceiverId:       balanceId,
		TransactionPrice: amount,
		OperationId:      PAY,
		Status:           COMPLETED,
	}

	if err := app.services.TransactionService.Update(r.Context(), trn); err != nil {
		if err := backUpBalance(app, r, blc, amount); err != nil {
			return err
		}
		return err
	}

	rep := &models.Report{
		BalanceId: balanceId,
		ServiceId: serviceId,
		Amount:    amount,
	}
	if err := app.services.ReportService.Create(r.Context(), rep); err != nil {
		if err := backUpBalance(app, r, blc, amount); err != nil {
			return err
		}
		return err
	}

	if err := internal.WriteJSON(w, trn, blc); err != nil {
		return err
	}
	return nil
}

// CancelOrder Если заказ отменен, возвращает деньги на основной баланс
func (app *application) CancelOrder(w http.ResponseWriter, r *http.Request) error {
	jsn, err := internal.ParsJSON(r)
	if err != nil {
		return err
	}

	balanceId, err := strconv.Atoi(jsn["balance_id"])
	if err != nil {
		return err
	}

	orderId, err := strconv.Atoi(jsn["order_id"])
	if err != nil {
		return err
	}

	amount, err := strconv.ParseFloat(jsn["amount"], 64)
	if err != nil {
		return err
	}

	blc, err := app.services.BalanceService.FindOne(r.Context(), balanceId)
	if err != nil || blc == nil {
		return apperror.NewAppError(err, "bad request", "balance not found", "some-code")
	}

	rdr, err := app.services.OrderService.FindOne(r.Context(), orderId)
	if rdr == nil {
		return apperror.NewAppError(err, "bad request", "order not found", "some code")
	}

	trn := &models.Transaction{
		SenderId:         balanceId,
		ReceiverId:       balanceId,
		TransactionPrice: amount,
		OperationId:      PAY,
		Status:           CANCELED,
	}
	if err := app.services.TransactionService.Update(r.Context(), trn); err != nil {
		return err
	}

	blc.AmountOnBalance += amount
	blc.AmountOnReserve -= amount
	if err := app.services.BalanceService.UpdateFull(r.Context(), blc); err != nil {
		return err
	}

	if err = internal.WriteJSON(w, trn, blc); err != nil {
		return err
	}

	return nil
}
