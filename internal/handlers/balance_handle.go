package handlers

import (
	"Task/internal"
	"Task/internal/apperror"
	"Task/internal/models"
	"Task/pkg/api"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

const (
	ASC  = "ASC"
	DESC = "DESC"
)

// GetBalance Показать баланс
func (app *application) GetBalance(w http.ResponseWriter, r *http.Request) error {
	jsn, err := internal.ParsJSON(r)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(jsn["balance_id"])
	if err != nil {
		return err
	}

	blc, err := app.services.BalanceService.FindOne(r.Context(), id)
	if err != nil {
		return err
	}

	if err := internal.WriteJSON(w, blc); err != nil {
		return err
	}

	return nil
}

// TopUpBalance Пополнить баланс
func (app *application) TopUpBalance(w http.ResponseWriter, r *http.Request) error {

	jsn, err := internal.ParsJSON(r)
	if err != nil {
		return err
	}

	balanceId, err := strconv.Atoi(jsn["balance_id"])
	if err != nil {
		return err
	}

	amount, err := strconv.ParseFloat(jsn["amount"], 64)
	if err != nil {
		return err
	}

	blc, err := app.services.BalanceService.FindOne(r.Context(), balanceId)

	var flag int
	if blc != nil {
		flag = Updated
		blc.AmountOnBalance += amount
		if err := app.services.BalanceService.Update(r.Context(), blc); err != nil {
			return err
		}
	} else {
		flag = Created
		blc = &models.Balance{
			Id:              balanceId,
			AmountOnBalance: amount,
			AmountOnReserve: 0,
		}
		if err := app.services.BalanceService.Create(r.Context(), blc); err != nil {
			return err
		}
	}

	trn := &models.Transaction{
		SenderId:         balanceId,
		ReceiverId:       balanceId,
		TransactionPrice: amount,
		OperationId:      REPLENISHMENT,
		Status:           COMPLETED,
	}

	err = app.services.TransactionService.Create(r.Context(), trn)
	if err != nil {
		switch flag {
		case Created:
			if err2 := app.services.BalanceService.Delete(r.Context(), balanceId); err != nil {
				return err2
			}
		case Updated:
			blc.AmountOnBalance -= amount
			if err2 := app.services.BalanceService.Update(r.Context(), blc); err != nil {
				return err2
			}
		}
		return apperror.NewAppError(err, "internal error", "transaction was not valid", "some-code")
	}

	if err := internal.WriteJSON(w, trn, blc); err != nil {
		return err
	}
	return nil
}

// GetTransactionHistory Получить список транзакций (принимает id пользователя и параметры сортировки и пагинации)
func (app *application) GetTransactionHistory(w http.ResponseWriter, r *http.Request) error {
	var params api.Params

	jsn, err := internal.ParsJSON(r)
	if err != nil {
		return err
	}

	balanceId, err := strconv.Atoi(jsn["balance_id"])
	if err != nil {
		return err
	}

	params.SortByDateOrder = strings.ToUpper(jsn["sort_by_date_order"])
	params.SortByAmountOrder = strings.ToUpper(jsn["sort_by_amount_order"])
	params.PerPage, err = strconv.Atoi(jsn["per_page"])
	params.Page, err = strconv.Atoi(jsn["page"])
	if err != nil {
		return err
	}

	switch params.SortByDateOrder {
	case ASC, DESC, "":
		break
	default:
		return errors.New("так нельзя сортировать")
	}

	switch params.SortByAmountOrder {
	case ASC, DESC, "":
		break
	default:
		return errors.New("так нельзя сортировать")
	}

	if params.Page == 0 {
		params.Page = 1
	}
	if params.PerPage == 0 {
		params.PerPage = 1
	}

	transactions, err := app.services.TransactionService.FindAll(r.Context(), balanceId, params)
	if err != nil {
		return err
	}

	err = internal.WriteJSON(w, transactions)
	if err != nil {
		return err
	}

	return nil
}
