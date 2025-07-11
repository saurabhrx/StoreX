package models

import "github.com/guregu/null"

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

type UserResponse struct {
	Name   string      `json:"name" db:"name"`
	Email  null.String `json:"email" db:"email"`
	Phone  null.String `json:"phone_no" db:"phone_no"`
	Assets []string    `json:"assets" db:"assets"`
	Role   string      `json:"role" db:"role"`
	Type   string      `json:"type" db:"type"`
}
