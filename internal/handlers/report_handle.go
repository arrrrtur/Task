package handlers

import (
	"Task/internal"
	"encoding/csv"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// @Summary Get report link
// @Tags balance
// @Description Create report file and get a link to a report in which statistics on services for profit
// @Accept json
// @Produce json
// @Param input body models.UnmarshalGetReport true "date info"
// @Router /report/get-link-report [GET]
func (app *application) GetReportLink(w http.ResponseWriter, r *http.Request) error {
	jsn, err := internal.ParsJSON(r)
	if err != nil {
		return err
	}

	yearMonth := strings.Split(jsn["year-month"], "-")

	year, err := strconv.Atoi(yearMonth[0])
	if err != nil {
		return err
	}

	month, err := strconv.Atoi(yearMonth[1])
	if err != nil {
		return err
	}

	reports, err := app.services.ReportService.FindAll(r.Context(), year, month)
	if err != nil {
		return err
	}

	if err := os.MkdirAll("reports", 0644); err != nil {
		return err
	}

	fileName := fmt.Sprintf("report%d-%d.csv", year, month)
	filePath := "reports/" + fileName

	csvFile, err := os.Create(filePath)
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	err = writer.WriteAll(reports)
	if err != nil {
		return err
	}
	writer.Flush()

	link := fmt.Sprintf("http://localhost:8080/file/%s", fileName)

	if err := internal.WriteJSON(w, fmt.Sprintf(`{"report_link": "%s"}`, link)); err != nil {
		return err
	}

	return nil
}

// @Summary Get report file
// @Tags balance
// @Description get report file with statistic
// @Accept json
// @Produce json
// @Param filename path string true "filename"
// @Router /file/{filename} [GET]
func (app *application) GetReportFile(w http.ResponseWriter, r *http.Request) error {
	// TODO how to read param
	params := httprouter.ParamsFromContext(r.Context())
	fileName := params.ByName("filename")
	filePath := "reports/" + fileName
	http.ServeFile(w, r, filePath)
	return nil
}
