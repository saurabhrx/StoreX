package dbhelper

import (
	"storeX/database"
	"storeX/models"
)

func IsVendorAlreadyExists(body *models.CreateVendorRequest) (string, error) {
	query := `SELECT id FROM vendors WHERE name=$1 AND phone_no=$2`
	var vendorID string
	err := database.STOREX.Get(&vendorID, query, body.Name, body.Phone)
	if err != nil {
		return "", err
	}
	return vendorID, nil
}
func CreateVendor(body *models.CreateVendorRequest) error {
	query := `INSERT INTO vendors(name, phone_no, address)
             VALUES($1,$2,$3)`
	_, err := database.STOREX.Exec(query, body.Name, body.Phone, body.Address)
	if err != nil {
		return err
	}
	return nil
}
