package models

const (
	RoleAdmin           = "admin"
	RoleEmployee        = "employee"
	RoleAssetManager    = "asset_manager"
	RoleEmployeeManager = "employee_manager"
)

const (
	TypeFullTime   = "full_time"
	TypeIntern     = "intern"
	TypeFreelancer = "freelancer"
)

type LoginUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type CreateUserRequest struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	EmployeeRole string `json:"employee_role"`
	EmployeeType string `json:"employee_type"`
	CreatedBy    string `json:"created_by"`
}
