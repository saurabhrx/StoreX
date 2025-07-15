package models

import (
	"encoding/json"
	"github.com/guregu/null"
)

type CreateAssetRequest struct {
	Brand             string          `json:"brand" validate:"required"`
	Model             string          `json:"model" validate:"required"`
	Serial            string          `json:"serial" validate:"required"`
	AssetType         string          `json:"assetType" validate:"required"`
	OwnedBy           string          `json:"ownedBy" validate:"required"`
	WarrantyStartDate string          `json:"warrantyStartDate" validate:"required"`
	WarrantyEndDate   string          `json:"warrantyEndDate" validate:"required"`
	PurchasedAt       string          `json:"purchasedAt" validate:"required"`
	Price             float64         `json:"price" validate:"required,gt=0"`
	CreatedBy         string          `json:"createdBy"`
	Specifications    json.RawMessage `json:"specifications" validate:"required"`
}

type LaptopSpecsRequest struct {
	AssetID   string `json:"assetID" validate:"required"`
	Ram       int    `json:"ram" validate:"required,gt=0"`
	Storage   int    `json:"storage" validate:"required,gt=0"`
	Processor string `json:"processor" validate:"required"`
	OS        string `json:"os" validate:"required"`
}

type MobileSpecsRequest struct {
	AssetID string `json:"assetID" validate:"required"`
	Ram     int    `json:"ram" validate:"required,gt=0"`
	Storage int    `json:"storage" validate:"required,gt=0"`
	OS      string `json:"os" validate:"required"`
	IMEI1   string `json:"imei_1" validate:"required,len=15"`
	IMEI2   string `json:"imei_2"`
}

type MonitorSpecsRequest struct {
	AssetID    string  `json:"assetID" validate:"required"`
	ScreenSize float64 `json:"screenSize" validate:"required,gt=0"`
	Resolution string  `json:"resolution" validate:"required"`
}

type MouseSpecsRequest struct {
	AssetID        string `json:"assetID" validate:"required"`
	ConnectionType string `json:"connectionType" validate:"required"`
	DPI            int    `json:"dpi" validate:"required,gt=0"`
}

type HardDiskSpecsRequest struct {
	AssetID   string `json:"assetID" validate:"required"`
	Type      string `json:"type" validate:"required"`
	Capacity  int    `json:"capacity" validate:"required,gt=0"`
	Interface string `json:"interface" validate:"required"`
	RPM       int    `json:"rpm" validate:"required,min=0"`
}

type PenDriveSpecsRequest struct {
	AssetID   string `json:"assetID" validate:"required"`
	Capacity  int    `json:"capacity" validate:"required,gt=0"`
	Interface string `json:"interface" validate:"required"`
}

type SimSpecsRequest struct {
	AssetID        string `json:"assetID" validate:"required"`
	SimNumber      string `json:"simNumber" validate:"required"`
	Career         string `json:"career" validate:"required"`
	PlanType       string `json:"planType" validate:"required"`
	ActivationDate string `json:"activationDate" validate:"required"`
}

type AccessoriesSpecsRequest struct {
	AssetID string `json:"assetID" validate:"required"`
	Type    string `json:"type" validate:"required"`
}

type AssignAssetRequest struct {
	EmployeeID string `json:"employeeID" validate:"required"`
	AssetID    string `json:"assetID" validate:"required"`
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
	Ram     int         `json:"ram" db:"ram"`
	Storage int         `json:"storage" db:"storage"`
	OS      string      `json:"os" db:"os"`
	IMEI1   string      `json:"imei1" db:"imei_1"`
	IMEI2   null.String `json:"imei2" db:"imei_2"`
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
