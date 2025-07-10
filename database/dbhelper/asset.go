package dbhelper

import (
	"storeX/database"
	"storeX/models"
)

func CreateAsset(body *models.CreateAssetRequest) error {
	query := `INSERT INTO assets(brand, model, serial_no, asset_type, owned_by, purchased_at, price, created_by) 
               VALUES($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := database.STOREX.Exec(query, body.Brand, body.Model, body.Serial, body.AssetType, body.OwnedBy, body.PurchasedAt, body.Price, body.CreatedBy)
	if err != nil {
		return err
	}
	return nil

}
func IsAssetExists(serial string) (string, error) {
	query := `SELECT id FROM assets WHERE serial_no=$1`
	var assetID string
	err := database.STOREX.Get(&assetID, query, serial)
	if err != nil {
		return "", err
	}
	return assetID, nil
}
func IsAssetAssign(assetID string) (string, error) {
	query := `SELECT employee_id FROM assigned_asset WHERE asset_id=$1 AND (end_date IS NULL OR end_date>NOW())`
	var emplID string
	err := database.STOREX.Get(&emplID, query, assetID)
	if err != nil {
		return "", err
	}
	return emplID, nil

}
func AssignAsset(body *models.AssignAssetRequest) error {
	query := `INSERT INTO assigned_asset(asset_id, employee_id) VALUES ($1,$2)`
	_, err := database.STOREX.Exec(query, body.AssetID, body.EmployeeID)
	if err != nil {
		return err
	}
	return nil
}
