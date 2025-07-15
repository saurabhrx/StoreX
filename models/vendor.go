package models

type CreateVendorRequest struct {
	Name    string `json:"name" validate:"required"`
	Phone   string `json:"phone" validate:"required"`
	Address string `json:"address" validate:"required"`
}
