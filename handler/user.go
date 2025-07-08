package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"storeX/database"
	"storeX/database/dbhelper"
	"storeX/middleware"
	"storeX/models"
	"storeX/utils"
)

var json = utils.JSON
var secretKey = []byte(os.Getenv("SECRET_KEY"))

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var body models.LoginUserRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	if body.Email == "" {
		utils.ResponseError(w, http.StatusBadRequest, "enter the email")
		return
	}
	isValid := utils.IsValidEmail(body.Email)
	if !isValid {
		utils.ResponseError(w, http.StatusBadRequest, "invalid email domain")
		return
	}
	firstName, lastName := utils.SplitName(body.Email)
	body.FirstName = firstName
	body.LastName = lastName
	var userID string
	fmt.Println(body)
	userID, err = dbhelper.IsUserExists(body.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, "failed to check user exists")
		return
	}
	var empRole string
	if userID == "" {
		txErr := database.Tx(func(tx *sqlx.Tx) error {
			userID, err = dbhelper.Register(tx, &body)
			if err != nil {
				return err
			}
			empRole = models.RoleEmployee
			err = dbhelper.CreateEmployeeRole(tx, userID)
			if err != nil {
				return err
			}
			err = dbhelper.CreateEmployeeType(tx, userID)
			if err != nil {
				return err
			}
			return nil
		})
		if txErr != nil {
			utils.ResponseError(w, http.StatusInternalServerError, "failed to login")
			return
		}
	} else {
		empRole, err = dbhelper.GetEmployeeRole(userID)
		if err != nil {
			utils.ResponseError(w, http.StatusInternalServerError, "failed to get user role")
			return
		}
	}
	accessToken, accessErr := middleware.GenerateAccessToken(userID, empRole)
	refreshToken, refreshErr := middleware.GenerateRefreshToken(userID)
	if accessErr != nil || refreshErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "could not generate jwt token")
		return
	}

	EncodeErr := json.NewEncoder(w).Encode(map[string]string{
		"message":       "user logged in successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
	if EncodeErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to send response")
	}
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var body models.CreateUserRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	if body.Email == "" {
		utils.ResponseError(w, http.StatusBadRequest, "enter the email")
		return
	}
	if body.FirstName == "" {
		utils.ResponseError(w, http.StatusBadRequest, "enter the first name")
		return
	}
	if body.Phone == "" {
		utils.ResponseError(w, http.StatusBadRequest, "enter the phone number")
		return
	}
	if body.EmployeeRole == "" {
		utils.ResponseError(w, http.StatusBadRequest, "enter the user role")
		return
	}
	if body.EmployeeType == "" {
		utils.ResponseError(w, http.StatusBadRequest, "enter the user type")
		return
	}

	userID := middleware.UserContext(r)
	body.CreatedBy = userID

}
