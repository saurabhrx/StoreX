package handler

import (
	"database/sql"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"os"
	"storeX/database"
	"storeX/database/dbhelper"
	"storeX/middleware"
	"storeX/models"
	"storeX/utils"
	"strings"
)

var jwtSecret = []byte(os.Getenv("SECRET_KEY"))

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var body models.LoginUserRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	if err := utils.Validate(body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
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
	if err := utils.Validate(body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	if !utils.IsValidEmail(body.Email) {
		utils.ResponseError(w, http.StatusBadRequest, "enter valid email")
		return
	}
	if !models.IsValidRole(body.EmployeeRole) {
		utils.ResponseError(w, http.StatusBadRequest, "enter valid user role")
		return
	}
	if !models.IsValidType(body.EmployeeType) {
		utils.ResponseError(w, http.StatusBadRequest, "enter valid user type")
		return
	}

	creatorID := middleware.UserContext(r)
	body.CreatedBy = creatorID
	userID, err := dbhelper.IsUserExists(body.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
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
		utils.ResponseError(w, http.StatusInternalServerError, "failed to get users")
		return
	}
	utils.ResponseJSON(w, http.StatusOK, body)
}

func UserTimeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user-id"]
	body, err := dbhelper.UserTimeline(userID)
	if err != nil {
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
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	body.UserID = userID
	err := dbhelper.UpdateUserDetails(&body)
	if err != nil {
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

	assetID, err := dbhelper.GetAssignedAsset(userID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to check is asset assigned")
		return
	}
	if assetID != "" {
		utils.ResponseError(w, http.StatusConflict, "asset is assigned to this user")
		return
	}

	err = dbhelper.DeleteUser(userID)
	if err != nil {
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

func RoleChange(w http.ResponseWriter, r *http.Request) {
	var body models.UserRoleChangeRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	err := dbhelper.RoleChange(&body)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to change role")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "user role changed successfully",
	})
}

func TypeChange(w http.ResponseWriter, r *http.Request) {
	var body models.UserTypeChangeRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	err := dbhelper.TypeChange(&body)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to change type")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "user type changed successfully",
	})
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	var body models.RefreshToken
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	if body.UserID == "" || body.Token == "" {
		utils.ResponseError(w, http.StatusUnauthorized, "session expired login again")
		return
	}
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(body.Token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		utils.ResponseError(w, http.StatusUnauthorized, "session expired login again")
		return

	}
	res, resErr := dbhelper.GetNameRoleByUserID(body.UserID)
	if resErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "error while getting role")
		return
	}
	accessToken, accessErr := middleware.GenerateAccessToken(body.UserID, res.RoleType, res.Name)
	if accessErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "could not generate access token")
		return
	}
	refreshToken, refreshErr := middleware.GenerateRefreshToken(body.UserID)
	if refreshErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "could not generate refresh token")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, struct {
		Status       int    `json:"status"`
		Message      string `json:"message"`
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
	}{
		Status:       http.StatusOK,
		Message:      "user type changed successfully",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
