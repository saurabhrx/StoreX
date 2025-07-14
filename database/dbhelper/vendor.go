package dbhelper

import (
	"storeX/database"
	"storeX/models"
)

func CreateVendor(body *models.CreateVendorRequest) error {
	query := `INSERT INTO vendors(name, phone_no, address)
             VALUES($1,$2,$3)`
	_, err := database.STOREX.Exec(query, body.Name, body.Phone, body.Address)
	if err != nil {
		return err
	}
	return nil
}
