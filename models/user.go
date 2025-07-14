package models

import (
	"github.com/guregu/null"
	"github.com/lib/pq"
)

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
	ID        string         `json:"id" db:"id"`
	Name      string         `json:"name" db:"name"`
	Email     string         `json:"email" db:"email"`
	Phone     null.String    `json:"phone_no" db:"phone_no"`
	Assets    pq.StringArray `json:"assets" db:"assets"`
	Role      string         `json:"role" db:"role"`
	Type      string         `json:"type" db:"type"`
	CreatedAt string         `json:"createdAt" db:"created_at"`
}

type UserTimeline struct {
	ID     string          `json:"id" db:"id"`
	Assets []AssignedAsset `json:"assets" db:"assets"`
}
type AssignedAsset struct {
	AssetID   string      `json:"assetID" db:"id"`
	Type      string      `json:"type" db:"asset_type"`
	StartDate string      `json:"startDate" db:"start_date"`
	EndDate   null.String `json:"endDate" db:"end_date"`
}

type DashboardResponse struct {
	ID     string                 `json:"id" db:"id"`
	Name   string                 `json:"name" db:"name"`
	Assets []AssignedAssetDetails `json:"assets" db:"assets"`
}

type AssignedAssetDetails struct {
	AssetID      string `json:"assetID" db:"id"`
	Brand        string `json:"brand" db:"brand"`
	Model        string `json:"model" db:"model"`
	Serial       string `json:"serial" db:"serial_no"`
	AssignedDate string `json:"assignedDate" db:"assigned_date"`
	Status       string `json:"status" db:"status"`
}

type UpdateUserDetails struct {
	UserID    string      `json:"userID"`
	FirstName string      `json:"firstName"`
	LastName  null.String `json:"lastName"`
	Phone     string      `json:"phone"`
}

type UserRoleChangeRequest struct {
	UserID  string `json:"userID"`
	NewRole string `json:"newRole"`
}
type UserTypeChangeRequest struct {
	UserID  string `json:"userID"`
	NewType string `json:"newType"`
}
