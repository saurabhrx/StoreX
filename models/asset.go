package models

import (
	"encoding/json"
	"github.com/guregu/null"
)

type CreateAssetRequest struct {
	Brand             string          `json:"brand"`
	Model             string          `json:"model"`
	Serial            string          `json:"serial"`
	AssetType         string          `json:"assetType"`
	OwnedBy           string          `json:"ownedBy"`
	WarrantyStartDate string          `json:"warrantyStartDate"`
	WarrantyEndDate   string          `json:"warrantyEndDate"`
	PurchasedAt       string          `json:"purchasedAt"`
	Price             float64         `json:"price"`
	CreatedBy         string          `json:"createdBy"`
	Specifications    json.RawMessage `json:"specifications"`
}
type LaptopSpecsRequest struct {
	AssetID   string `json:"assetID"`
	Ram       int    `json:"ram"`
	Storage   int    `json:"storage"`
	Processor string `json:"processor"`
	OS        string `json:"os"`
}
type MobileSpecsRequest struct {
	AssetID string `json:"assetID"`
	Ram     int    `json:"ram"`
	Storage int    `json:"storage"`
	OS      string `json:"os"`
	IMEI1   string `json:"imei_1"`
	IMEI2   string `json:"imei_2"`
}
type MonitorSpecsRequest struct {
	AssetID    string  `json:"assetID"`
	ScreenSize float64 `json:"screenSize"`
	Resolution string  `json:"resolution"`
}
type MouseSpecsRequest struct {
	AssetID        string `json:"assetID"`
	ConnectionType string `json:"connectionType"`
	DPI            int    `json:"dpi"`
}

type HardDiskSpecsRequest struct {
	AssetID   string `json:"assetID"`
	Type      string `json:"type"`
	Capacity  int    `json:"capacity"`
	Interface string `json:"interface"`
	RPM       int    `json:"rpm"`
}

type PenDriveSpecsRequest struct {
	AssetID   string `json:"assetID"`
	Capacity  int    `json:"capacity"`
	Interface string `json:"interface"`
}
type SimSpecsRequest struct {
	AssetID        string `json:"assetID"`
	SimNumber      string `json:"simNumber"`
	Career         string `json:"career"`
	PlanType       string `json:"planType"`
	ActivationDate string `json:"activationDate"`
}
type AccessoriesSpecsRequest struct {
	AssetID string `json:"assetID"`
	Type    string `json:"type"`
}

type AssignAssetRequest struct {
	EmployeeID string `json:"employeeID"`
	AssetID    string `json:"assetID"`
}

type AssetResponse struct {
	ID             string      `json:"id" db:"id"`
	Brand          string      `json:"brand" db:"brand"`
	Model          string      `json:"model" db:"model"`
	AssetType      string      `json:"assetType" db:"asset_type"`
	Serial         string      `json:"serial" db:"serial_no"`
	AssetStatus    string      `json:"assetStatus" db:"status"`
	AssignedTo     null.String `json:"assignedTo" db:"assigned_to"`
	AssignedDate   null.String `json:"assignedDate" db:"assigned_date"`
	OwnedBy        string      `json:"ownedBy" db:"owned_by"`
	PurchasedAt    string      `json:"purchasedAt" db:"purchased_at"`
	Specifications interface{} `json:"specifications" db:"specifications"`
}

type AssetTimeline struct {
	AssetID  string            `json:"assetId" db:"asset_id"`
	Employee []AssetAssignedTo `json:"employee" db:"employee"`
}

type AssetAssignedTo struct {
	EmpID     string      `json:"empID" db:"id"`
	Name      string      `json:"name" db:"name"`
	Email     string      `json:"email" db:"email"`
	StartDate string      `json:"startDate" db:"start_date"`
	EndDate   null.String `json:"endDate" db:"end_date"`
}

type LaptopSpecsResponse struct {
	Ram       int    `json:"ram" db:"ram"`
	Storage   int    `json:"storage" db:"storage"`
	Processor string `json:"processor" db:"processor"`
	OS        string `json:"os" db:"os"`
}
type MobileSpecsResponse struct {
	Ram     int    `json:"ram" db:"ram"`
	Storage int    `json:"storage" db:"storage"`
	OS      string `json:"os" db:"os"`
	IMEI1   string `json:"imei1" db:"imei_1"`
	IMEI2   string `json:"imei2" db:"imei_2"`
}
type MonitorSpecsResponse struct {
	ScreenSize float64 `json:"screenSize" db:"screen_size"`
	Resolution string  `json:"resolution" db:"resolution"`
}
type MouseSpecsResponse struct {
	ConnectionType string `json:"connectionType" db:"connection_type"`
	DPI            int    `json:"dpi" db:"dpi"`
}

type HardDiskSpecsResponse struct {
	Type      string `json:"type" db:"type"`
	Capacity  int    `json:"capacity" db:"capacity"`
	Interface string `json:"interface" db:"interface"`
	RPM       int    `json:"rpm" db:"rpm"`
}

type PenDriveSpecsResponse struct {
	Capacity  int    `json:"capacity" db:"capacity"`
	Interface string `json:"interface" db:"interface"`
}
type SimSpecsResponse struct {
	SimNumber      string `json:"simNumber" db:"sim_number"`
	Career         string `json:"career" db:"career"`
	PlanType       string `json:"planType" db:"plan_type"`
	ActivationDate string `json:"activationDate" db:"activation_date"`
}
type AccessoriesSpecsResponse struct {
	Type string `json:"type" db:"type"`
}

type AssignAssetResponse struct {
	EmployeeID string `json:"employeeID" db:"employee_id"`
	AssetID    string `json:"assetID" db:"asset_id"`
}

type AssetStatsResponse struct {
	Total            int `json:"total" db:"total"`
	Available        int `json:"available" db:"available"`
	Assigned         int `json:"assigned" db:"assigned"`
	WaitingForRepair int `json:"waitingForRepair" db:"waiting_for_repair"`
	Service          int `json:"service" db:"service"`
	Damaged          int `json:"damaged" db:"damaged"`
}
