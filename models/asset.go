package models

import "encoding/json"

type CreateAssetRequest struct {
	Brand          string          `json:"brand"`
	Model          string          `json:"model"`
	Serial         string          `json:"serial"`
	AssetType      string          `json:"assetType"`
	OwnedBy        string          `json:"ownedBy"`
	PurchasedAt    string          `json:"purchasedAt"`
	Price          int             `json:"price"`
	CreatedBy      string          `json:"createdBy"`
	Type           string          `json:"type"`
	Specifications json.RawMessage `json:"specifications"`
}
type LaptopSpecs struct {
	AssetID   string `json:"assetID"`
	Ram       int    `json:"ram"`
	Storage   int    `json:"storage"`
	Processor string `json:"processor"`
	OS        string `json:"os"`
}
type MobileSpecs struct {
	AssetID string `json:"assetID"`
	Ram     int    `json:"ram"`
	Storage int    `json:"storage"`
	OS      string `json:"os"`
	IMEI1   string `json:"imei_1"`
	IMEI2   string `json:"imei_2"`
}
type MonitorSpecs struct {
	AssetID    string  `json:"assetID"`
	ScreenSize float64 `json:"screenSize"`
	Resolution string  `json:"resolution"`
}
type MouseSpecs struct {
	AssetID        string `json:"assetID"`
	ConnectionType string `json:"connectionType"`
	DPI            int    `json:"dpi"`
}

type HardDiskSpecs struct {
	AssetID   string `json:"assetID"`
	Type      string `json:"type"`
	Capacity  int    `json:"capacity"`
	Interface string `json:"interface"`
	RPM       int    `json:"rpm"`
}

type PenDriveSpecs struct {
	AssetID   string `json:"assetID"`
	Capacity  int    `json:"capacity"`
	Interface string `json:"interface"`
}
type SimSpecs struct {
	AssetID        string `json:"assetID"`
	SimNumber      string `json:"simNumber"`
	Career         string `json:"career"`
	PlanType       string `json:"planType"`
	ActivationDate string `json:"activationDate"`
}
type AccessoriesSpecs struct {
	AssetID string `json:"assetID"`
	Type    string `json:"type"`
}

type AssignAssetRequest struct {
	EmployeeID string `json:"employeeID"`
	AssetID    string `json:"assetID"`
}
