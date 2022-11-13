package handlers

import (
	"Task/internal/apperror"
	"Task/internal/services"
	"Task/pkg/logging"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const (
	Created = iota
	Updated
)

const (
	CANCELED = iota + 1
	IN_PROCESSING
	COMPLETED
)

type Application interface {
	Register(router *httprouter.Router)
}

type application struct {
	logger   *logging.Logger
	services services.Services
}

func NewApplication(services services.Services, logger *logging.Logger) Application {
	return &application{
		services: services,
		logger:   logger,
	}
}

func (app *application) Register(router *httprouter.Router) {

	// work with balance
	router.HandlerFunc(http.MethodPost, "/balance/top-up-balance", apperror.Middleware(app.TopUpBalance))
	router.HandlerFunc(http.MethodPatch, "/balance/reserve-from-balance", apperror.Middleware(app.ReserveFromBalance))
	router.HandlerFunc(http.MethodPatch, "/balance/has-passed", apperror.Middleware(app.WriteOffMoney))
	router.HandlerFunc(http.MethodGet, "/balance/get-balance", apperror.Middleware(app.GetBalance))
	router.HandlerFunc(http.MethodPatch, "/balance/cancel-order", apperror.Middleware(app.CancelOrder))
	router.HandlerFunc(http.MethodGet, "/balance/get-history", apperror.Middleware(app.GetTransactionHistory))

	// work with report
	router.HandlerFunc(http.MethodGet, "/report/get-link-report", apperror.Middleware(app.GetReportLink))
	router.HandlerFunc(http.MethodGet, "/file/:filename", apperror.Middleware(app.GetReportFile))

}
