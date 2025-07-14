package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"storeX/database"
	"storeX/database/dbhelper"
	"storeX/middleware"
	"storeX/models"
	"storeX/utils"
	"strings"
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
	name := body.FirstName + " " + body.LastName
	accessToken, accessErr := middleware.GenerateAccessToken(userID, empRole, name)
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

func GetUsers(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query()
	search := queryParam.Get("search")
	empRole := queryParam.Get("role")
	empType := queryParam.Get("type")
	empRoleArray := utils.UserRoleArray(empRole)
	empTypeArray := utils.UserTypeArray(empType)

	var filters models.UserFilter
	filters.Search = search
	filters.EmpType = empTypeArray
	filters.EmpRole = empRoleArray
	filters.Limit, filters.Offset = utils.Pagination(r)
	filters.IsSearchText = strings.TrimSpace(search) != ""
	body, err := dbhelper.GetUsers(&filters)
	if err != nil {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, "failed to get users")
		return
	}
	utils.ResponseJSON(w, http.StatusOK, body)
}

func UserTimeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user-id"]
	fmt.Println(userID)
	body, err := dbhelper.UserTimeline(userID)
	if err != nil {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, "failed to get user timeline")
		return
	}
	utils.ResponseJSON(w, http.StatusOK, body)
}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.UserContext(r)
	name := middleware.NameContext(r)
	var body models.DashboardResponse
	body.ID = userID
	body.Name = name
	err := dbhelper.Dashboard(userID, &body)
	if err != nil {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, "failed to get dashboard")
		return
	}
	utils.ResponseJSON(w, http.StatusOK, body)
}

func UpdateUserDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user-id"]
	var body models.UpdateUserDetails
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		fmt.Println(parseErr)
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	body.UserID = userID
	err := dbhelper.UpdateUserDetails(&body)
	if err != nil {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, "failed to update user details")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "user updated successfully",
	})
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user-id"]

	err := dbhelper.DeleteUser(userID)
	if err != nil {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "user deleted successfully",
	})
}
