package utils

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"io"
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
func ResponseJSON(w http.ResponseWriter, statusCode int, body interface{}) {
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		logrus.Errorf("failed to send the error %+v", err)
	}
}

func ParseBody(body io.Reader, out interface{}) error {
	err := json.NewDecoder(body).Decode(out)
	if err != nil {
		return err
	}
	return nil
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
func UserRoleArray(empRole string) []string {
	if empRole == "" {
		return []string{}
	}
	return strings.Split(empRole, ",")
}
func UserTypeArray(empType string) []string {
	if empType == "" {
		return []string{}
	}
	return strings.Split(empType, ",")
}
func AssetTypeArray(assetType string) []string {
	if assetType == "" {
		return []string{}
	}
	return strings.Split(assetType, ",")
}
func AssetStatusArray(status string) []string {
	if status == "" {
		return []string{}
	}
	return strings.Split(status, ",")
}
func OwnedByArray(ownedBy string) []string {
	if ownedBy == "" {
		return []string{}
	}
	return strings.Split(ownedBy, ",")
}

func Validate(i interface{}) error {
	var validate = validator.New()
	err := validate.Struct(i)
	if err == nil {
		return nil
	}
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var sb strings.Builder
		for _, fe := range ve {
			sb.WriteString(fmt.Sprintf("Field '%s' failed on '%s' validation; ", fe.Field(), fe.Tag()))
		}
		return fmt.Errorf(sb.String())
	}
	return err
}
