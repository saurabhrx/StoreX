package dbhelper

import (
	"github.com/jmoiron/sqlx"
	"storeX/database"
	"storeX/models"
)

func CreateService(db sqlx.Ext, body *models.ServiceRequest) error {
	args := []interface{}{
		body.AssetID,
		body.VendorID,
		body.StartDate,
		body.EndDate,
		body.Cost,
		body.Remark,
	}
	query := `INSERT INTO services(asset_id, vendor_id,start_date, end_date, cost, remark)
              VALUES ($1,$2,$3,$4,$5,$6)`
	_, err := db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil

}
func UpdateService(body *models.ServiceRequest) error {
	args := []interface{}{
		body.AssetID,
		body.VendorID,
		body.StartDate,
		body.EndDate,
		body.Cost,
		body.Remark,
	}
	query := `UPDATE services SET vendor_id=$2, cost=$5 , remark=$6 , start_date=$3, end_date=$4 
                WHERE asset_id=$1 `
	_, err := database.STOREX.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil

}
