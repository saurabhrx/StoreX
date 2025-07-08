package utils

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary
var json = JSON

type clientError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func ResponseError(w http.ResponseWriter, statusCode int, message string) {
	logrus.Errorf("status : %d, message : %s", statusCode, message)
	clientErr := &clientError{
		StatusCode: statusCode,
		Message:    message,
	}
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(clientErr); err != nil {
		logrus.Errorf("failed to send the error %+v", err)
	}
}

func Pagination(r *http.Request) (int, int) {
	page := 1
	limit := 10
	queryParams := r.URL.Query()
	if pageValue := queryParams.Get("page"); pageValue != "" {
		if p, err := strconv.Atoi(queryParams.Get("page")); err == nil && p > 0 {
			page = p
		}
	}
	if limitValue := queryParams.Get("limit"); limitValue != "" {
		if l, err := strconv.Atoi(queryParams.Get("limit")); err == nil && l > 0 {
			limit = l
		}
	}
	offset := (page - 1) * limit

	return limit, offset
}

func IsValidEmail(email string) bool {
	parts := strings.Split(email, "@")
	if parts[1] == "remotestate.com" {
		return true
	}
	return false
}

func SplitName(email string) (string, string) {
	parts := strings.Split(strings.Split(email, "@")[0], ".")
	firstName := parts[0]
	lastName := parts[1]
	return firstName, lastName
}
