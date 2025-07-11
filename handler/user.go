package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"net/http"
	"storeX/database"
	"storeX/database/dbhelper"
	"storeX/middleware"
	"storeX/models"
	"storeX/utils"
)

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var body models.LoginUserRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
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
	userID, err := dbhelper.IsUserExists(body.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, "failed to check user exists")
		return
	}
	var empRole string
	var empType string
	if userID == "" {
		txErr := database.Tx(func(tx *sqlx.Tx) error {
			userID, err = dbhelper.Register(tx, &body)
			if err != nil {
				return err
			}
			empRole = models.RoleEmployee
			err = dbhelper.CreateUserRole(tx, userID, empRole)
			if err != nil {
				return err
			}
			empType = models.TypeFullTime
			err = dbhelper.CreateUserType(tx, userID, empType)
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
		empRole, err = dbhelper.GetUserRole(userID)
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

	utils.ResponseJSON(w, http.StatusOK, struct {
		Status       int    `json:"status"`
		Message      string `json:"message"`
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		Status:       http.StatusOK,
		Message:      "user logged in successfully",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var body models.CreateUserRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	fmt.Println(body)
	if body.Email == "" || !utils.IsValidEmail(body.Email) {
		utils.ResponseError(w, http.StatusBadRequest, "enter valid email")
		return
	}
	if body.FirstName == "" {
		utils.ResponseError(w, http.StatusBadRequest, "enter valid first name")
		return
	}
	if body.Phone == "" || len(body.Phone) != 10 {
		utils.ResponseError(w, http.StatusBadRequest, "enter valid number")
		return
	}
	if body.EmployeeRole == "" || !models.IsValidRole(body.EmployeeRole) {
		utils.ResponseError(w, http.StatusBadRequest, "enter valid user role")
		return
	}
	if body.EmployeeType == "" || !models.IsValidType(body.EmployeeType) {
		utils.ResponseError(w, http.StatusBadRequest, "enter valid user type")
		return
	}

	creatorID := middleware.UserContext(r)
	body.CreatedBy = creatorID
	userID, err := dbhelper.IsUserExists(body.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, "failed to check user exists")
		return
	}
	if userID != "" {
		utils.ResponseError(w, http.StatusConflict, "user already exists")
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, err = dbhelper.CreateUser(tx, &body)
		if err != nil {
			return err
		}
		err = dbhelper.CreateUserRole(tx, userID, body.EmployeeRole)
		if err != nil {
			return err
		}
		err = dbhelper.CreateUserType(tx, userID, body.EmployeeType)
		if err != nil {
			return err
		}

		return nil
	})
	if txErr != nil {
		fmt.Println(txErr)
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusCreated,
		Message: "user created successfully",
	})

}

func GetEmployees(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query()
	search := queryParam.Get("search")
	empRole := queryParam.Get("role")
	empType := queryParam.Get("type")
	empRoleArray := utils.RoleArray(empRole)
	empTypeArray := utils.TypeArray(empType)

	fmt.Println(search)
	fmt.Println(empRoleArray)
	fmt.Println("type...", empTypeArray)

	body, err := dbhelper.GetUsers(search, empRoleArray, empTypeArray)
	if err != nil {
		fmt.Println(err)
	}
	utils.ResponseJSON(w, http.StatusOK, body)
}
