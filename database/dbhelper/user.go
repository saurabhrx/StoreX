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
func CreateUserRole(db sqlx.Ext, userID, empRole string) error {
	query := `INSERT INTO employee_role(employee_id,role) VALUES ($1,$2)`
	_, err := db.Exec(query, userID, empRole)
	if err != nil {
		return err
	}
	return nil
}
func GetUserRole(userID string) (string, error) {
	query := `SELECT role FROM employee_role WHERE employee_id=$1`
	var role string
	err := database.STOREX.Get(&role, query, userID)
	if err != nil {
		return "", err
	}
	return role, nil
}
func CreateUserType(db sqlx.Ext, userID, empType string) error {
	query := `INSERT INTO employee_type(employee_id,type) VALUES ($1,$2)`
	_, err := db.Exec(query, userID, empType)
	if err != nil {
		return err
	}
	return nil
}
func CreateUser(db sqlx.Ext, body *models.CreateUserRequest) (string, error) {
	query := `INSERT INTO employees(first_name, last_name, email, phone_no,created_by) VALUES ($1,$2,$3,$4,$5) RETURNING id`
	var userID string
	err := db.QueryRowx(query, body.FirstName, body.LastName, body.Email, body.Phone, body.CreatedBy).Scan(&userID)
	if err != nil {
		return "", err
	}
	return userID, nil
}
