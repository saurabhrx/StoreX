package models

import "github.com/guregu/null"

type ServiceRequest struct {
	AssetID   string      `json:"assetID"`
	VendorID  string      `json:"vendorID"`
	Cost      null.Float  `json:"cost"`
	Remark    string      `json:"remark"`
	StartDate string      `json:"startDate"`
	EndDate   null.String `json:"endDate"`
}
