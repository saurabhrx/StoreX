package handler

import (
	"github.com/gorilla/mux"
	"net/http"
	"storeX/database/dbhelper"
	"storeX/models"
	"storeX/utils"
)

func Service(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["asset-id"]
	var body models.ServiceRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse body")
		return
	}
	body.AssetID = assetID
	err := dbhelper.CreateService(&body)
	if err != nil {
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
