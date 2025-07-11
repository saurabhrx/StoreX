package handler

import (
	"database/sql"
	"encoding/json"
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

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var body models.CreateAssetRequest
	var laptopSpecs models.LaptopSpecs
	var mobileSpecs models.MobileSpecs
	var mouseSpecs models.MouseSpecs
	var monitorSpecs models.MonitorSpecs
	var hardDriveSpecs models.HardDiskSpecs
	var penDriveSpecs models.PenDriveSpecs
	var simSpecs models.SimSpecs
	var accessoriesSpecs models.AccessoriesSpecs

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	
	if body.Type == "laptop" {
		err := json.Unmarshal(body.Specifications, &laptopSpecs)
		fmt.Println(laptopSpecs)
		if err != nil {
			fmt.Println(err)
		}
	}
	if body.Type == "mobile" {
		err := json.Unmarshal(body.Specifications, &mobileSpecs)
		fmt.Println(mobileSpecs)
		if err != nil {
			fmt.Println(err)
		}
	}
	if body.Type == "mouse" {
		err := json.Unmarshal(body.Specifications, &mouseSpecs)
		fmt.Println(mouseSpecs)
		if err != nil {
			fmt.Println(err)
		}
	}
	if body.Type == "monitor" {
		err := json.Unmarshal(body.Specifications, &monitorSpecs)
		fmt.Println(monitorSpecs)
		if err != nil {
			fmt.Println(err)
		}
	}
	if body.Type == "hard_drive" {
		err := json.Unmarshal(body.Specifications, &hardDriveSpecs)
		fmt.Println(hardDriveSpecs)
		if err != nil {
			fmt.Println(err)
		}
	}
	if body.Type == "pen_drive" {
		err := json.Unmarshal(body.Specifications, &penDriveSpecs)
		fmt.Println(penDriveSpecs)
		if err != nil {
			fmt.Println(err)
		}
	}
	if body.Type == "sim" {
		err := json.Unmarshal(body.Specifications, &simSpecs)
		fmt.Println(simSpecs)
		if err != nil {
			fmt.Println(err)
		}
	}
	if body.Type == "accessories" {
		err := json.Unmarshal(body.Specifications, &accessoriesSpecs)
		fmt.Println(accessoriesSpecs)
		if err != nil {
			fmt.Println(err)
		}
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
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		assetID, err = dbhelper.CreateAsset(tx, &body)
		if err != nil {
			return err
		}
		fmt.Println(assetID)
		switch body.Type {
		case "laptop":
			laptopSpecs.AssetID = assetID
			if err := dbhelper.CreateLaptopSpecs(tx, &laptopSpecs); err != nil {
				return err
			}
		case "mobile":
			mobileSpecs.AssetID = assetID
			if err := dbhelper.CreateMobileSpecs(tx, &mobileSpecs); err != nil {
				return err
			}
		case "mouse":
			mouseSpecs.AssetID = assetID
			if err := dbhelper.CreateMouseSpecs(tx, &mouseSpecs); err != nil {
				return err
			}
		case "monitor":
			monitorSpecs.AssetID = assetID
			if err := dbhelper.CreateMonitorSpecs(tx, &monitorSpecs); err != nil {
				return err
			}
		case "hard_disk":
			hardDriveSpecs.AssetID = assetID
			if err := dbhelper.CreateHardDiskSpecs(tx, &hardDriveSpecs); err != nil {
				return err
			}
		case "pen_drive":
			penDriveSpecs.AssetID = assetID
			if err := dbhelper.CreatePenDriveSpecs(tx, &penDriveSpecs); err != nil {
				return err
			}
		case "sim":
			simSpecs.AssetID = assetID
			if err := dbhelper.CreateSimSpecs(tx, &simSpecs); err != nil {
				return err
			}
		case "accessories":
			accessoriesSpecs.AssetID = assetID
			if err := dbhelper.CreateAccessoriesSpecs(tx, &accessoriesSpecs); err != nil {
				return err
			}

		}

		return nil
	})
	if txErr != nil {
		fmt.Println(txErr)
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
