package handler

import (
	"fmt"
	"net/http"
	"storeX/database/dbhelper"
	"storeX/models"
	"storeX/utils"
)

func CreateVendor(w http.ResponseWriter, r *http.Request) {
	var body models.CreateVendorRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse body")
		return
	}
	if err := utils.Validate(body); err != nil {
		fmt.Println("Validation error:", err)
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	err := dbhelper.CreateVendor(&body)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create vendor")
		return
	}
	utils.ResponseJSON(w, http.StatusOK, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "vendor created successfully",
	})

}
