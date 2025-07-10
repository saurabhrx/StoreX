package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"storeX/database/dbhelper"
	"storeX/middleware"
	"storeX/models"
	"storeX/utils"
)

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var body models.CreateAssetRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	if body.Brand == "" {
		utils.ResponseError(w, http.StatusBadRequest, "enter valid brand")
		return
	}
	if body.Model == "" {
		utils.ResponseError(w, http.StatusBadRequest, "enter valid model")
		return
	}
	if body.Serial == "" {
		utils.ResponseError(w, http.StatusBadRequest, "enter valid serial")
		return
	}
	if body.AssetType == "" {
		utils.ResponseError(w, http.StatusBadRequest, "enter valid asset type")
		return
	}
	if body.PurchasedAt == "" {
		utils.ResponseError(w, http.StatusBadRequest, "enter valid date")
		return
	}
	if body.Price < 0 {
		utils.ResponseError(w, http.StatusBadRequest, "price must be greater than or equal to 0")
		return
	}
	assetID, err := dbhelper.IsAssetExists(body.Serial)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, "failed to check asset exists")
		return
	}
	if assetID != "" {
		utils.ResponseError(w, http.StatusConflict, "asset already exists")
		return
	}
	CreatorId := middleware.UserContext(r)
	body.CreatedBy = CreatorId
	if err := dbhelper.CreateAsset(&body); err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create asset")
		return
	}
	utils.ResponseJSON(w, http.StatusOK, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "asset created successfully",
	})

}

func AssignAsset(w http.ResponseWriter, r *http.Request) {
	var body models.AssignAssetRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse body")
		return
	}
	emplID, err := dbhelper.IsAssetAssign(body.AssetID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, "failed to check is asset assigned")
		return
	}
	if emplID != "" {
		utils.ResponseError(w, http.StatusConflict, "asset assigned to someone else")
		return
	}
	if err := dbhelper.AssignAsset(&body); err != nil {
		fmt.Println(err)
		utils.ResponseError(w, http.StatusInternalServerError, "failed to assign asset")
		return
	}

	utils.ResponseJSON(w, http.StatusCreated, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusCreated,
		Message: "asset assigned successfully",
	})
}
