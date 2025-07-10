package models

type CreateAssetRequest struct {
	Brand       string `json:"brand"`
	Model       string `json:"model"`
	Serial      string `json:"serial"`
	AssetType   string `json:"assetType"`
	OwnedBy     string `json:"ownedBy"`
	PurchasedAt string `json:"purchasedAt"`
	Price       int    `json:"price"`
	CreatedBy   string `json:"createdBy"`
}
type AssignAssetRequest struct {
	EmployeeID string `json:"employeeID"`
	AssetID    string `json:"assetID"`
}
