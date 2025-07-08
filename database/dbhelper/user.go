package dbhelper

import (
	"github.com/jmoiron/sqlx"
	"storeX/database"
	"storeX/models"
)

func Register(db sqlx.Ext, body *models.LoginUserRequest) (string, error) {
	query := `INSERT INTO employees(first_name, last_name, email) VALUES ($1,$2,$3) RETURNING id`
	var userID string
	err := db.QueryRowx(query, body.FirstName, body.LastName, body.Email).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}

func IsUserExists(email string) (string, error) {
	query := `SELECT id FROM employees WHERE email=$1`
	var userID string
	err := database.STOREX.Get(&userID, query, email)
	if err != nil {
		return "", err
	}
	return userID, nil
}
func CreateEmployeeRole(db sqlx.Ext, userID string) error {
	query := `INSERT INTO employee_role(employee_id) VALUES ($1)`
	_, err := db.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
func GetEmployeeRole(userID string) (string, error) {
	query := `SELECT role FROM employee_role WHERE employee_id=$1`
	var role string
	err := database.STOREX.Get(&role, query, userID)
	if err != nil {
		return "", err
	}
	return role, nil
}
func CreateEmployeeType(db sqlx.Ext, userID string) error {
	query := `INSERT INTO employee_type(employee_id) VALUES ($1)`
	_, err := db.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
