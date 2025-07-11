package dbhelper

import (
	"github.com/jmoiron/sqlx"
	"storeX/database"
	"storeX/models"
)

func CreateAsset(db sqlx.Ext, body *models.CreateAssetRequest) (string, error) {
	query := `INSERT INTO assets(brand, model, serial_no, asset_type, owned_by, purchased_at, price, created_by) 
               VALUES($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id`
	var assetID string
	err := db.QueryRowx(query, body.Brand, body.Model, body.Serial, body.AssetType, body.OwnedBy, body.PurchasedAt, body.Price, body.CreatedBy).Scan(&assetID)
	if err != nil {
		return "", err
	}
	return assetID, nil

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
func CreateLaptopSpecs(db sqlx.Ext, specs *models.LaptopSpecs) error {
	query := `INSERT INTO laptop_specs(asset_id, ram, storage_capacity, processor, os) 
              VALUES($1,$2,$3,$4,$5) `
	_, err := db.Exec(query, specs.AssetID, specs.Ram, specs.Storage, specs.Processor, specs.OS)
	if err != nil {
		return err
	}
	return nil
}
func CreateMobileSpecs(db sqlx.Ext, specs *models.MobileSpecs) error {
	query := `INSERT INTO mobile_specs(asset_id, ram, storage_capacity, os, imei_1, imei_2) 
              VALUES($1,$2,$3,$4,$5,$6) `
	_, err := db.Exec(query, specs.AssetID, specs.Ram, specs.Storage, specs.OS, specs.IMEI1, specs.IMEI2)
	if err != nil {
		return err
	}
	return nil
}
func CreateMouseSpecs(db sqlx.Ext, specs *models.MouseSpecs) error {
	query := `INSERT INTO mouse_specs(asset_id, connection_type, dpi) 
              VALUES($1,$2,$3) `
	_, err := db.Exec(query, specs.AssetID, specs.ConnectionType, specs.DPI)
	if err != nil {
		return err
	}
	return nil
}
func CreateMonitorSpecs(db sqlx.Ext, specs *models.MonitorSpecs) error {
	query := `INSERT INTO monitor_specs(asset_id, screen_size, resolution) 
              VALUES($1,$2,$3) `
	_, err := db.Exec(query, specs.AssetID, specs.AssetID, specs.ScreenSize, specs.Resolution)
	if err != nil {
		return err
	}
	return nil
}
func CreateHardDiskSpecs(db sqlx.Ext, specs *models.HardDiskSpecs) error {
	query := `INSERT INTO hard_disk_specs(asset_id, type, capacity, interface, rpm) 
              VALUES($1,$2,$3,$4,$5) `
	_, err := db.Exec(query, specs.AssetID, specs.Type, specs.Capacity, specs.Interface, specs.RPM)
	if err != nil {
		return err
	}
	return nil
}
func CreatePenDriveSpecs(db sqlx.Ext, specs *models.PenDriveSpecs) error {
	query := `INSERT INTO pen_drive_specs(asset_id, capacity, interface) 
              VALUES($1,$2,$3) `
	_, err := db.Exec(query, specs.AssetID, specs.Capacity, specs.Interface)
	if err != nil {
		return err
	}
	return nil
}
func CreateSimSpecs(db sqlx.Ext, specs *models.SimSpecs) error {
	query := `INSERT INTO sim_specs(asset_id, sim_number, career, plan_type, activation_date) 
              VALUES($1,$2,$3,$4,$5) `
	_, err := db.Exec(query, specs.AssetID, specs.SimNumber, specs.Career, specs.PlanType, specs.ActivationDate)
	if err != nil {
		return err
	}
	return nil
}
func CreateAccessoriesSpecs(db sqlx.Ext, specs *models.AccessoriesSpecs) error {
	query := `INSERT INTO accessories_specs(asset_id, type) 
              VALUES($1,$2) `
	_, err := db.Exec(query, specs.AssetID, specs.Type)
	if err != nil {
		return err
	}
	return nil
}
