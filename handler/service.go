package handler

import (
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"net/http"
	"storeX/database"
	"storeX/database/dbhelper"
	"storeX/models"
	"storeX/utils"
)

func CreateService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["asset-id"]
	var body models.ServiceRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse body")
		return
	}
	if err := utils.Validate(body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	body.AssetID = assetID

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err := dbhelper.CreateService(tx, &body)
		if err != nil {
			return err
		}
		err = dbhelper.ChangeAssetStatus(tx, assetID, "service")
		if err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create service")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "service created successfully",
	})

}
func UpdateService(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["asset-id"]
	var body models.ServiceRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse body")
		return
	}
	if err := utils.Validate(body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	body.AssetID = assetID

	err := dbhelper.UpdateService(&body)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to update service")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "service updated successfully",
	})

}
