package dbhelper

import (
	"storeX/database"
	"storeX/models"
)

func CreateService(body *models.ServiceRequest) error {
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
	_, err := database.STOREX.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil

}
