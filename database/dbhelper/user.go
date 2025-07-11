package dbhelper

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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

func GetUsers(search string, empRole, empType []string) ([]models.UserResponse, error) {
	fmt.Println(empRole, empType)
	query := `SELECT CONCAT(e.first_name,' ',e.last_name) AS name, 
                     e.email, 
                     e.phone_no,
                     t.type, 
                     r.role, 
                     ARRAY_AGG(DISTINCT a.id) AS assets
              FROM employees e
              JOIN employee_role r ON e.id = r.employee_id
              JOIN employee_type t ON e.id = t.employee_id
              JOIN assigned_asset aa ON e.id = aa.employee_id
              JOIN assets a ON aa.asset_id = a.id
              WHERE ($1 = '' OR 
                     CONCAT(e.first_name, ' ', e.last_name) ILIKE $2 OR 
                     e.email ILIKE $2 OR 
                     e.phone_no ILIKE $2)
                AND (r.role = ANY($3::role_type[]))
                AND ( t.type = ANY($4::empl_type[]))
              GROUP BY e.first_name, e.last_name, e.email, e.phone_no, t.type, r.role`

	finalSearch := "%" + search + "%"
	var users []models.UserResponse
	err := database.STOREX.Select(&users, query, search, finalSearch, pq.Array(empRole), pq.Array(empType))
	if err != nil {
		return nil, err
	}
	return users, nil
}
