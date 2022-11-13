package internal

import (
	"Task/internal/apperror"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func ParsJSON(r *http.Request) (map[string]string, error) {
	var result map[string]string
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&result)
	if err != nil {
		logrus.WithError(err).Errorf("error parcing JSON")
		logrus.Trace(fmt.Sprintf("values - %s", result))
		return nil, errors.New("error parcing JSON")
	}

	return result, nil
}

func FormatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

func WriteJSON(w http.ResponseWriter, objs ...interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	JSONSlice, err := json.Marshal(objs)
	if err != nil {
		return apperror.NewAppError(err, "bebebe", "error in writeJSON", "some code")
	}

	w.Write(JSONSlice)

	return nil
}
