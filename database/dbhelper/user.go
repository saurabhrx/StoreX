package dbhelper

import (
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

func GetUsers(filters *models.UserFilter) ([]models.UserResponse, error) {
	args := []interface{}{
		!filters.IsSearchText,
		filters.Search,
		pq.Array(filters.EmpRole),
		pq.Array(filters.EmpType),
		filters.Limit,
		filters.Offset,
	}

	query := `SELECT e.id,
                     CONCAT(e.first_name,' ',COALESCE(e.last_name, '')) AS name, 
                     e.email, 
                     e.phone_no,
                     t.type, 
                     r.role, 
                     COALESCE( ARRAY_AGG(DISTINCT a.id) FILTER (WHERE a.id IS NOT NULL),
                      '{}') AS assets,
                     e.created_at
              FROM employees e
              JOIN employee_role r ON e.id = r.employee_id
              JOIN employee_type t ON e.id = t.employee_id
              LEFT JOIN assigned_asset aa ON e.id = aa.employee_id
              LEFT JOIN assets a ON aa.asset_id = a.id
              WHERE ($1 OR 
                     CONCAT(e.first_name, ' ', e.last_name) ILIKE '%' || $2::TEXT || '%' OR 
                     e.email  ILIKE '%' || $2::TEXT || '%'  OR 
                     e.phone_no  ILIKE '%' || $2::TEXT || '%' )
                AND (CARDINALITY($3::role_type[])=0 OR r.role = ANY($3::role_type[]))
                AND (CARDINALITY($4::empl_type[])=0 OR t.type = ANY($4::empl_type[]))
                AND e.archived_at IS NULL
              GROUP BY e.id, t.type, r.role,e.created_at
              ORDER BY e.created_at DESC 
              LIMIT $5 OFFSET $6`

	var users []models.UserResponse
	err := database.STOREX.Select(&users, query, args...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func UserTimeline(userID string) (models.UserTimeline, error) {
	query := `SELECT a.id , a.asset_type,aa.start_date,aa.end_date 
               FROM employees e
               JOIN assigned_asset aa ON aa.employee_id=e.id AND e.id=$1
               JOIN assets a ON a.id=aa.asset_id ORDER BY aa.start_date DESC`
	var body models.UserTimeline
	body.ID = userID
	err := database.STOREX.Select(&body.Assets, query, userID)
	if err != nil {
		return models.UserTimeline{}, err
	}
	return body, nil
}
func Dashboard(userID string, body *models.DashboardResponse) error {
	query := `SELECT assets.id, assets.brand , assets.model , assets.serial_no,
                assets.status , aa.start_date as assigned_date FROM assets
               JOIN assigned_asset aa ON assets.id = aa.asset_id AND aa.employee_id=$1`
	err := database.STOREX.Select(&body.Assets, query, userID)
	if err != nil {
		return err
	}
	return nil
}

func UpdateUserDetails(body *models.UpdateUserDetails) error {
	args := []interface{}{
		body.FirstName,
		body.LastName,
		body.Phone,
		body.UserID,
	}
	query := `UPDATE employees SET first_name=$1 , last_name=$2 , phone_no=$3
             WHERE id=$4`
	_, err := database.STOREX.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
func DeleteUser(userID string) error {
	query := `UPDATE employees SET archived_at=NOW()
               WHERE id=$1`
	_, err := database.STOREX.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
