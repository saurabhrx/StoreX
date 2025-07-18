package handler

import (
	"database/sql"
	"encoding/json"
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

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var body models.CreateAssetRequest
	var laptopSpecs models.LaptopSpecsRequest
	var mobileSpecs models.MobileSpecsRequest
	var mouseSpecs models.MouseSpecsRequest
	var monitorSpecs models.MonitorSpecsRequest
	var hardDriveSpecs models.HardDiskSpecsRequest
	var penDriveSpecs models.PenDriveSpecsRequest
	var simSpecs models.SimSpecsRequest
	var accessoriesSpecs models.AccessoriesSpecsRequest

	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	if err := utils.Validate(body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	switch body.AssetType {
	case "laptop":
		if err := json.Unmarshal(body.Specifications, &laptopSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid laptop specifications format")
			return
		}
		if err := utils.Validate(laptopSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}

	case "mobile":
		if err := json.Unmarshal(body.Specifications, &mobileSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid mobile specifications format")
			return
		}
		if err := utils.Validate(mobileSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}

	case "mouse":
		if err := json.Unmarshal(body.Specifications, &mouseSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid mouse specifications format")
			return
		}
		if err := utils.Validate(mouseSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}

	case "monitor":
		if err := json.Unmarshal(body.Specifications, &monitorSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid monitor specifications format")
			return
		}
		if err := utils.Validate(monitorSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}

	case "hard_disk":
		if err := json.Unmarshal(body.Specifications, &hardDriveSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid hard disk specifications format")
			return
		}
		if err := utils.Validate(hardDriveSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}

	case "pen_drive":
		if err := json.Unmarshal(body.Specifications, &penDriveSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid pen drive specifications format")
			return
		}
		if err := utils.Validate(penDriveSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}

	case "sim":
		if err := json.Unmarshal(body.Specifications, &simSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid sim specifications format")
			return
		}
		if err := utils.Validate(simSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}

	case "accessories":
		if err := json.Unmarshal(body.Specifications, &accessoriesSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid accessories specifications format")
			return
		}
		if err := utils.Validate(accessoriesSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}

	default:
		utils.ResponseError(w, http.StatusBadRequest, "unknown asset type")
		return
	}

	assetID, err := dbhelper.IsAssetExists(body.Serial)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
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
		if err := dbhelper.CreateWarranty(tx, assetID, body.WarrantyStartDate, body.WarrantyEndDate); err != nil {
			return err
		}
		switch body.AssetType {
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
		default:
			return fmt.Errorf("invalid asset type: %s", body.AssetType)
		}

		return nil
	})
	if txErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create asset")
		return
	}
	utils.ResponseJSON(w, http.StatusCreated, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusCreated,
		Message: "asset created successfully",
	})

}

func AssignAsset(w http.ResponseWriter, r *http.Request) {
	var body models.AssignAssetRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse body")
		return
	}
	if err := utils.Validate(body); err != nil {
		utils.ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}
	emplID, err := dbhelper.IsAssetAssign(body.AssetID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to check is asset assigned")
		return
	}
	if emplID != "" {
		utils.ResponseError(w, http.StatusConflict, "asset assigned to someone else")
		return
	}
	status, StatusErr := dbhelper.IsAssetAvailable(body.AssetID)
	if StatusErr != nil && !errors.Is(StatusErr, sql.ErrNoRows) {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to check is asset available")
		return
	}
	if status != "available" {
		utils.ResponseError(w, http.StatusConflict, "asset is not available")
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		if err := dbhelper.AssignAsset(tx, &body); err != nil {
			return err
		}
		if err := dbhelper.ChangeAssetStatus(tx, body.AssetID, "assigned"); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
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

func GetAssets(w http.ResponseWriter, r *http.Request) {
	queryParam := r.URL.Query()
	search := queryParam.Get("search")
	assetType := queryParam.Get("type")
	assetStatus := queryParam.Get("status")
	ownedBy := queryParam.Get("ownedBy")

	assetTypeArray := utils.AssetTypeArray(assetType)
	assetStatusArray := utils.AssetStatusArray(assetStatus)
	ownedByArray := utils.OwnedByArray(ownedBy)

	var filters models.AssetFilter

	filters.Search = search
	filters.AssetType = assetTypeArray
	filters.AssetStatus = assetStatusArray
	filters.OwnedType = ownedByArray
	filters.IsSearchText = strings.TrimSpace(search) != ""
	filters.Limit, filters.Offset = utils.Pagination(r)

	body, err := dbhelper.GetAllAssets(&filters)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to get assets")
		return
	}
	for i := range body {
		var spec interface{}
		switch body[i].AssetType {
		case "laptop":
			spec, err = dbhelper.GetLaptopSpec(body[i].ID)
		case "mobile":
			spec, err = dbhelper.GetMobileSpec(body[i].ID)
		case "mouse":
			spec, err = dbhelper.GetMouseSpec(body[i].ID)
		case "monitor":
			spec, err = dbhelper.GetMonitorSpec(body[i].ID)
		case "hard_disk":
			spec, err = dbhelper.GetHardDiskSpec(body[i].ID)
		case "pen_drive":
			spec, err = dbhelper.GetPenDriveSpec(body[i].ID)
		case "sim":
			spec, err = dbhelper.GetSimSpec(body[i].ID)
		case "accessories":
			spec, err = dbhelper.GetAccessoriesSpec(body[i].ID)
		default:
			continue
		}

		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			utils.ResponseError(w, http.StatusInternalServerError, "failed to fetch specifications")
			return
		}

		body[i].Specifications = spec
	}

	utils.ResponseJSON(w, http.StatusOK, body)

}

func AssetTimeline(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["asset-id"]
	body, err := dbhelper.AssetTimeline(assetID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to get asset timeline")
		return
	}
	utils.ResponseJSON(w, http.StatusOK, body)
}

func UnassignAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["asset-id"]
	var body models.ReasonOfRetrieve
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse body")
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		if err := dbhelper.UnassignAsset(tx, assetID, body.Reason); err != nil {
			return err
		}
		if err := dbhelper.ChangeAssetStatus(tx, assetID, "available"); err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to unassign the asset")
		return
	}
	utils.ResponseJSON(w, http.StatusOK, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusCreated,
		Message: "asset unassigned successfully",
	})
}

func AssetStats(w http.ResponseWriter, r *http.Request) {
	body, err := dbhelper.AssetStats()
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to get asset counts")
		return
	}
	utils.ResponseJSON(w, http.StatusOK, body)
}

func DeleteAsset(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["asset-id"]
	emplID, err := dbhelper.IsAssetAssign(assetID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to check is asset assigned")
		return
	}
	if emplID != "" {
		utils.ResponseError(w, http.StatusConflict, "asset assigned to someone")
		return
	}
	err = dbhelper.DeleteAsset(assetID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "asset deleted successfully",
	})

}

func ChangeStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["asset-id"]
	var body models.ChangeStatus
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse body")
		return
	}
	if !models.IsValidStatus(body.Status) {
		utils.ResponseError(w, http.StatusBadRequest, "incorrect status type")
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		err := dbhelper.ChangeAssetStatus(tx, assetID, body.Status)
		if err != nil {
			return err
		}
		return nil

	})
	if txErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to update status")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "asset status updated successfully",
	})

}

func UpdateAssetSpecs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	assetID := vars["asset-id"]
	var body models.UpdateAssetSpecsRequest
	if parseErr := utils.ParseBody(r.Body, &body); parseErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse body")
		return
	}
	var err error

	switch body.AssetType {
	case "laptop":
		var laptopSpecs models.LaptopSpecsRequest
		if err = json.Unmarshal(body.Specifications, &laptopSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid laptop specifications format")
			return
		}
		if err = utils.Validate(laptopSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}
		laptopSpecs.AssetID = assetID
		err = dbhelper.UpdateLaptopSpecs(&laptopSpecs)

	case "mobile":
		var mobileSpecs models.MobileSpecsRequest
		if err = json.Unmarshal(body.Specifications, &mobileSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid mobile specifications format")
			return
		}
		if err = utils.Validate(mobileSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}
		mobileSpecs.AssetID = assetID
		err = dbhelper.UpdateMobileSpecs(&mobileSpecs)

	case "monitor":
		var monitorSpecs models.MonitorSpecsRequest
		if err = json.Unmarshal(body.Specifications, &monitorSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid monitor specifications format")
			return
		}
		if err = utils.Validate(monitorSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}
		monitorSpecs.AssetID = assetID
		err = dbhelper.UpdateMonitorSpecs(&monitorSpecs)

	case "mouse":
		var mouseSpecs models.MouseSpecsRequest
		if err = json.Unmarshal(body.Specifications, &mouseSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid mouse specifications format")
			return
		}
		if err = utils.Validate(mouseSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}
		mouseSpecs.AssetID = assetID
		err = dbhelper.UpdateMouseSpecs(&mouseSpecs)

	case "hard_disk":
		var hdSpecs models.HardDiskSpecsRequest
		if err = json.Unmarshal(body.Specifications, &hdSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid hard disk specifications format")
			return
		}
		if err = utils.Validate(hdSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}
		hdSpecs.AssetID = assetID
		err = dbhelper.UpdateHardDiskSpecs(&hdSpecs)

	case "pen_drive":
		var pdSpecs models.PenDriveSpecsRequest
		if err = json.Unmarshal(body.Specifications, &pdSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid pen drive specifications format")
			return
		}
		if err = utils.Validate(pdSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}
		pdSpecs.AssetID = assetID
		err = dbhelper.UpdatePenDriveSpecs(&pdSpecs)

	case "sim":
		var simSpecs models.SimSpecsRequest
		if err = json.Unmarshal(body.Specifications, &simSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid sim specifications format")
			return
		}
		if err = utils.Validate(simSpecs); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, err.Error())
			return
		}
		simSpecs.AssetID = assetID
		err = dbhelper.UpdateSimSpecs(&simSpecs)

	default:
		utils.ResponseError(w, http.StatusBadRequest, "asset type is incorrect")
		return
	}

	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to update asset specs")
		return
	}

	utils.ResponseJSON(w, http.StatusOK, struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}{
		Status:  http.StatusOK,
		Message: "asset specs updated successfully",
	})

}
