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

func IsValidRole(r string) bool {
	return r == RoleAdmin || r == RoleEmployee || r == RoleEmployeeManager || r == RoleAssetManager
}
func IsValidType(t string) bool {
	return t == TypeIntern || t == TypeFreelancer || t == TypeFullTime
}

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
	EmployeeRole string `json:"employeeRole"`
	EmployeeType string `json:"employeeType"`
	CreatedBy    string `json:"createdBy"`
}
