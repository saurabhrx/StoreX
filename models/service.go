package models

import "github.com/guregu/null"

type ServiceRequest struct {
	AssetID   string      `json:"assetID"`
	VendorID  string      `json:"vendorID" validate:"required"`
	Cost      null.Float  `json:"cost"`
	Remark    null.String `json:"remark"`
	StartDate string      `json:"startDate" validate:"required"`
	EndDate   null.String `json:"endDate"`
}
