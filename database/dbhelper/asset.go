package dbhelper

import (
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
	query := `SELECT employee_id FROM assigned_asset WHERE asset_id=$1 AND end_date IS NULL`
	var emplID string
	err := database.STOREX.Get(&emplID, query, assetID)
	if err != nil {
		return "", err
	}
	return emplID, nil

}
func IsAssetAvailable(assetID string) (string, error) {
	query := `SELECT status FROM assets WHERE id=$1`
	var status string
	err := database.STOREX.Get(&status, query, assetID)
	if err != nil {
		return "", err
	}
	return status, nil

}

func AssignAsset(db sqlx.Ext, body *models.AssignAssetRequest) error {
	query := `INSERT INTO assigned_asset(asset_id, employee_id) VALUES ($1,$2)`
	_, err := db.Exec(query, body.AssetID, body.EmployeeID)
	if err != nil {
		return err
	}
	return nil
}
func ChangeAssetStatus(db sqlx.Ext, assetID, status string) error {
	query := `UPDATE assets SET status=$2 WHERE id=$1`
	_, err := db.Exec(query, assetID, status)
	if err != nil {
		return err
	}
	return nil
}
func CreateLaptopSpecs(db sqlx.Ext, specs *models.LaptopSpecsRequest) error {
	query := `INSERT INTO laptop_specs(asset_id, ram, storage_capacity, processor, os) 
              VALUES($1,$2,$3,$4,$5)`
	_, err := db.Exec(query, specs.AssetID, specs.Ram, specs.Storage, specs.Processor, specs.OS)
	if err != nil {
		return err
	}
	return nil
}
func CreateMobileSpecs(db sqlx.Ext, specs *models.MobileSpecsRequest) error {
	query := `INSERT INTO mobile_specs(asset_id, ram, storage_capacity, os, imei_1, imei_2) 
              VALUES($1,$2,$3,$4,$5,$6) `
	_, err := db.Exec(query, specs.AssetID, specs.Ram, specs.Storage, specs.OS, specs.IMEI1, specs.IMEI2)
	if err != nil {
		return err
	}
	return nil
}
func CreateMouseSpecs(db sqlx.Ext, specs *models.MouseSpecsRequest) error {
	query := `INSERT INTO mouse_specs(asset_id, connection_type, dpi) 
              VALUES($1,$2,$3) `
	_, err := db.Exec(query, specs.AssetID, specs.ConnectionType, specs.DPI)
	if err != nil {
		return err
	}
	return nil
}
func CreateMonitorSpecs(db sqlx.Ext, specs *models.MonitorSpecsRequest) error {
	query := `INSERT INTO monitor_specs(asset_id, screen_size, resolution) 
              VALUES($1,$2,$3) `
	_, err := db.Exec(query, specs.AssetID, specs.ScreenSize, specs.Resolution)
	if err != nil {
		return err
	}
	return nil
}
func CreateHardDiskSpecs(db sqlx.Ext, specs *models.HardDiskSpecsRequest) error {
	query := `INSERT INTO hard_disk_specs(asset_id, type, capacity, interface, rpm) 
              VALUES($1,$2,$3,$4,$5) `
	_, err := db.Exec(query, specs.AssetID, specs.Type, specs.Capacity, specs.Interface, specs.RPM)
	if err != nil {
		return err
	}
	return nil
}
func CreatePenDriveSpecs(db sqlx.Ext, specs *models.PenDriveSpecsRequest) error {
	query := `INSERT INTO pen_drive_specs(asset_id, capacity, interface) 
              VALUES($1,$2,$3) `
	_, err := db.Exec(query, specs.AssetID, specs.Capacity, specs.Interface)
	if err != nil {
		return err
	}
	return nil
}
func CreateSimSpecs(db sqlx.Ext, specs *models.SimSpecsRequest) error {
	query := `INSERT INTO sim_specs(asset_id, sim_number, career, plan_type, activation_date) 
              VALUES($1,$2,$3,$4,$5) `
	_, err := db.Exec(query, specs.AssetID, specs.SimNumber, specs.Career, specs.PlanType, specs.ActivationDate)
	if err != nil {
		return err
	}
	return nil
}
func CreateAccessoriesSpecs(db sqlx.Ext, specs *models.AccessoriesSpecsRequest) error {
	query := `INSERT INTO accessories_specs(asset_id, type) 
              VALUES($1,$2) `
	_, err := db.Exec(query, specs.AssetID, specs.Type)
	if err != nil {
		return err
	}
	return nil
}
func CreateWarranty(db sqlx.Ext, assetID, warrantyStart, warrantyEnd string) error {
	query := `INSERT INTO warranty(asset_id, start_date,end_date) 
              VALUES($1,$2,$3) `
	_, err := db.Exec(query, assetID, warrantyStart, warrantyEnd)
	if err != nil {
		return err
	}
	return nil
}

func GetAllAssets(filters *models.AssetFilter) ([]models.AssetResponse, error) {
	args := []interface{}{
		!filters.IsSearchText,
		filters.Search,
		pq.Array(filters.AssetStatus),
		pq.Array(filters.AssetType),
		pq.Array(filters.OwnedType),
		filters.Limit,
		filters.Offset,
	}

	query := `SELECT assets.id,
       assets.brand,
       assets.model,
       assets.asset_type,
       assets.serial_no,
       assets.status,
       assets.owned_by,
       assets.purchased_at,
       CONCAT(first_name,' ',COALESCE(e.last_name, '')) as assigned_to,
       aa.start_date  as assigned_date
FROM assets
         LEFT JOIN assigned_asset aa ON assets.id = aa.asset_id
         LEFT JOIN employees e ON aa.employee_id=e.id
WHERE ($1 OR assets.brand ILIKE '%' || $2::TEXT || '%'
    OR assets.model ILIKE '%' || $2::TEXT || '%' 
    OR assets.serial_no ILIKE '%' || $2::TEXT || '%'
    OR CONCAT(first_name,' ',e.last_name) ILIKE '%' || $2::TEXT || '%') 
    AND (CARDINALITY($3::status_type[])=0 OR assets.status=ANY($3::status_type[]))
    AND (CARDINALITY($4::assets_type[])=0 OR assets.asset_type=ANY($4::assets_type[]))
    AND (CARDINALITY($5::owned_type[])=0 OR assets.owned_by=ANY($5::owned_type[]))
    ORDER BY aa.start_date DESC 
    LIMIT $6 OFFSET $7`

	var users []models.AssetResponse
	err := database.STOREX.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}
	return users, nil

}

func AssetTimeline(assetID string) (models.AssetTimeline, error) {
	query := `SELECT
		NULL::UUID AS employee_id,
		NULL::TEXT AS name,
		NULL::TEXT AS email,
		s.start_date,
		s.end_date,
		s.remark AS remark,
		'service' AS record_type
	FROM services s
	WHERE s.asset_id = $1

	UNION ALL

	SELECT
		aa.employee_id,
		CONCAT(e.first_name, ' ', COALESCE(e.last_name, '')) AS name,
		e.email,
		aa.start_date,
		aa.end_date,
		aa.reason_to_return AS remark,
		'assigned' AS record_type
	FROM assigned_asset aa
	JOIN employees e ON aa.employee_id = e.id
	WHERE aa.asset_id = $1
	ORDER BY start_date DESC`

	var body models.AssetTimeline
	body.AssetID = assetID
	err := database.STOREX.Select(&body.Assigned, query, assetID)
	if err != nil {
		return models.AssetTimeline{}, err
	}
	return body, nil
}

func GetLaptopSpec(assetID string) (models.LaptopSpecsResponse, error) {
	query := `SELECT ram,storage_capacity as storage , processor,os 
              FROM laptop_specs WHERE asset_id=$1`
	var body models.LaptopSpecsResponse
	err := database.STOREX.Get(&body, query, assetID)
	if err != nil {
		return models.LaptopSpecsResponse{}, err
	}
	return body, nil

}

func GetMobileSpec(assetID string) (models.MobileSpecsResponse, error) {
	query := `SELECT ram,storage_capacity as storage,os,imei_1,imei_2
              FROM mobile_specs WHERE asset_id=$1`
	var body models.MobileSpecsResponse
	err := database.STOREX.Get(&body, query, assetID)
	if err != nil {
		return models.MobileSpecsResponse{}, err
	}
	return body, nil

}

func GetMouseSpec(assetID string) (models.MouseSpecsResponse, error) {
	query := `SELECT connection_type,dpi
              FROM mouse_specs WHERE asset_id=$1`
	var body models.MouseSpecsResponse
	err := database.STOREX.Get(&body, query, assetID)
	if err != nil {
		return models.MouseSpecsResponse{}, err
	}
	return body, nil

}

func GetMonitorSpec(assetID string) (models.MonitorSpecsResponse, error) {
	query := `SELECT  screen_size,resolution
              FROM monitor_specs WHERE asset_id=$1`
	var body models.MonitorSpecsResponse
	err := database.STOREX.Get(&body, query, assetID)
	if err != nil {
		return models.MonitorSpecsResponse{}, err
	}
	return body, nil

}

func GetHardDiskSpec(assetID string) (models.HardDiskSpecsResponse, error) {
	query := `SELECT  type,capacity,interface,rpm
              FROM hard_disk_specs WHERE asset_id=$1`
	var body models.HardDiskSpecsResponse
	err := database.STOREX.Get(&body, query, assetID)
	if err != nil {
		return models.HardDiskSpecsResponse{}, err
	}
	return body, nil

}

func GetPenDriveSpec(assetID string) (models.PenDriveSpecsResponse, error) {
	query := `SELECT capacity,interface
              FROM pen_drive_specs WHERE asset_id=$1`
	var body models.PenDriveSpecsResponse
	err := database.STOREX.Get(&body, query, assetID)
	if err != nil {
		return models.PenDriveSpecsResponse{}, err
	}
	return body, nil

}

func GetSimSpec(assetID string) (models.SimSpecsResponse, error) {
	query := `SELECT sim_number,career,plan_type,activation_date
              FROM sim_specs WHERE asset_id=$1`
	var body models.SimSpecsResponse
	err := database.STOREX.Get(&body, query, assetID)
	if err != nil {
		return models.SimSpecsResponse{}, err
	}
	return body, nil

}
func GetAccessoriesSpec(assetID string) (models.AccessoriesSpecsResponse, error) {
	query := `SELECT type
              FROM accessories_specs WHERE asset_id=$1`
	var body models.AccessoriesSpecsResponse
	err := database.STOREX.Get(&body, query, assetID)
	if err != nil {
		return models.AccessoriesSpecsResponse{}, err
	}
	return body, nil

}

func UnassignAsset(db sqlx.Ext, assetID, reason string) error {
	query := `UPDATE assigned_asset SET end_date=NOW(), reason_to_return=$2 WHERE asset_id=$1`
	_, err := db.Exec(query, assetID, reason)
	if err != nil {
		return err
	}
	return nil

}

func AssetStats() (models.AssetStatsResponse, error) {
	query := `SELECT COUNT(*) as total,
                COUNT(*) FILTER (WHERE status = 'available') as available,
                COUNT(*) FILTER (WHERE status = 'assigned') as assigned,
                COUNT(*) FILTER (WHERE status = 'waiting_for_repair') as waiting_for_repair,
                COUNT(*) FILTER (WHERE status = 'service') as service,
                COUNT(*) FILTER (WHERE status = 'damaged') as damaged
              FROM assets
              WHERE archived_at IS NULL`

	var body models.AssetStatsResponse
	err := database.STOREX.Get(&body, query)
	if err != nil {
		return models.AssetStatsResponse{}, err
	}
	return body, nil

}

func DeleteAsset(assetID string) error {
	query := `UPDATE assets SET archived_at=NOW() , status='deleted'
               WHERE id=$1`
	_, err := database.STOREX.Exec(query, assetID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateLaptopSpecs(body *models.LaptopSpecsRequest) error {
	args := []interface{}{
		body.Ram,
		body.Storage,
		body.Processor,
		body.OS,
		body.AssetID,
	}
	query := `UPDATE laptop_specs SET ram=$1, storage_capacity=$2, processor=$3, os=$4 , updated_at=NOW() WHERE asset_id=$5`
	_, err := database.STOREX.Exec(query, args...)
	return err
}
func UpdateMobileSpecs(body *models.MobileSpecsRequest) error {
	args := []interface{}{
		body.Ram,
		body.Storage,
		body.OS,
		body.IMEI1,
		body.IMEI2,
		body.AssetID,
	}
	query := `UPDATE mobile_specs SET ram=$1, storage_capacity=$2, os=$3, imei_1=$4, imei_2=$5 , updated_at=NOW() WHERE asset_id=$6`
	_, err := database.STOREX.Exec(query, args...)
	return err
}

func UpdateMouseSpecs(body *models.MouseSpecsRequest) error {
	args := []interface{}{
		body.ConnectionType,
		body.DPI,
		body.AssetID,
	}
	query := `UPDATE mouse_specs SET connection_type=$1, dpi=$2 , updated_at=NOW() WHERE asset_id=$3`
	_, err := database.STOREX.Exec(query, args...)
	return err
}

func UpdateMonitorSpecs(body *models.MonitorSpecsRequest) error {
	args := []interface{}{
		body.ScreenSize,
		body.Resolution,
		body.AssetID,
	}
	query := `UPDATE monitor_specs SET screen_size=$1, resolution=$2 , updated_at=NOW() WHERE asset_id=$3`
	_, err := database.STOREX.Exec(query, args...)
	return err
}

func UpdateHardDiskSpecs(body *models.HardDiskSpecsRequest) error {
	args := []interface{}{
		body.Type,
		body.Capacity,
		body.Interface,
		body.RPM,
		body.AssetID,
	}
	query := `UPDATE hard_disk_specs SET type=$1, capacity=$2, interface=$3, rpm=$4, updated_at=NOW() WHERE asset_id=$5`
	_, err := database.STOREX.Exec(query, args...)
	return err
}

func UpdatePenDriveSpecs(body *models.PenDriveSpecsRequest) error {
	args := []interface{}{
		body.Capacity,
		body.Interface,
		body.AssetID,
	}
	query := `UPDATE pen_drive_specs SET capacity=$1, interface=$2 , updated_at=NOW() WHERE asset_id=$3`
	_, err := database.STOREX.Exec(query, args...)
	return err
}

func UpdateSimSpecs(body *models.SimSpecsRequest) error {
	args := []interface{}{
		body.SimNumber,
		body.Career,
		body.PlanType,
		body.ActivationDate,
		body.AssetID,
	}
	query := `UPDATE sim_specs SET sim_number=$1, career=$2, plan_type=$3, activation_date=$4 , updated_at=NOW() WHERE asset_id=$5`
	_, err := database.STOREX.Exec(query, args...)
	return err
}
