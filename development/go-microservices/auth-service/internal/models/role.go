package models

import (
	"fmt"

	"gorm.io/gorm"
)

// Permission represents a specific action that can be performed
type Permission string

const (
	// Order Management Permissions
	PermissionCreateOrder Permission = "create_order"
	PermissionViewOrder   Permission = "view_order"
	PermissionUpdateOrder Permission = "update_order"
	PermissionDeleteOrder Permission = "delete_order"
	PermissionVoidOrder   Permission = "void_order"
	PermissionRefundOrder Permission = "refund_order"

	// Payment Management Permissions
	PermissionProcessPayment Permission = "process_payment"
	PermissionViewPayment    Permission = "view_payment"
	PermissionRefundPayment  Permission = "refund_payment"

	// Menu Management Permissions
	PermissionCreateItem Permission = "create_item"
	PermissionUpdateItem Permission = "update_item"
	PermissionDeleteItem Permission = "delete_item"
	PermissionViewItem   Permission = "view_item"

	// Inventory Management Permissions
	PermissionManageStock Permission = "manage_stock"
	PermissionViewStock   Permission = "view_stock"
	PermissionAdjustStock Permission = "adjust_stock"

	// Table Management Permissions
	PermissionManageTable Permission = "manage_table"
	PermissionViewTable   Permission = "view_table"
	PermissionAssignTable Permission = "assign_table"

	// Reporting Permissions
	PermissionViewReports   Permission = "view_reports"
	PermissionExportReports Permission = "export_reports"

	// User Management Permissions
	PermissionManageUsers Permission = "manage_users"
	PermissionViewUsers   Permission = "view_users"
	PermissionCreateUser  Permission = "create_user"
	PermissionUpdateUser  Permission = "update_user"
	PermissionDeleteUser  Permission = "delete_user"

	// Settings Management Permissions
	PermissionManageSettings Permission = "manage_settings"
	PermissionViewSettings   Permission = "view_settings"
)

// Role represents a collection of permissions
type Role struct {
	ID          string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string `gorm:"not null;unique"`
	Description string
	Permissions []Permission `gorm:"type:text[]"`
	IsActive    bool         `gorm:"default:true"`
}

// PredefinedRoles contains all the standard roles with their permissions
var PredefinedRoles = map[string][]Permission{
	"super_admin": {
		PermissionCreateOrder, PermissionViewOrder, PermissionUpdateOrder, PermissionDeleteOrder,
		PermissionVoidOrder, PermissionRefundOrder,
		PermissionProcessPayment, PermissionViewPayment, PermissionRefundPayment,
		PermissionCreateItem, PermissionUpdateItem, PermissionDeleteItem, PermissionViewItem,
		PermissionManageStock, PermissionViewStock, PermissionAdjustStock,
		PermissionManageTable, PermissionViewTable, PermissionAssignTable,
		PermissionViewReports, PermissionExportReports,
		PermissionManageUsers, PermissionViewUsers, PermissionCreateUser,
		PermissionUpdateUser, PermissionDeleteUser,
		PermissionManageSettings, PermissionViewSettings,
	},
	"manager": {
		PermissionCreateOrder, PermissionViewOrder, PermissionUpdateOrder,
		PermissionProcessPayment, PermissionViewPayment,
		PermissionCreateItem, PermissionUpdateItem, PermissionViewItem,
		PermissionManageStock, PermissionViewStock,
		PermissionManageTable, PermissionViewTable,
		PermissionViewReports, PermissionExportReports,
		PermissionViewUsers,
		PermissionViewSettings,
	},
	"clerk": {
		PermissionCreateOrder, PermissionViewOrder,
		PermissionProcessPayment, PermissionViewPayment,
		PermissionViewItem,
		PermissionViewStock,
		PermissionViewTable,
	},
	"kitchen_staff": {
		PermissionViewOrder, PermissionUpdateOrder,
		PermissionViewItem,
	},
	"waiter": {
		PermissionCreateOrder, PermissionViewOrder,
		PermissionViewItem,
		PermissionViewTable, PermissionAssignTable,
	},
	"cashier": {
		PermissionCreateOrder, PermissionViewOrder,
		PermissionProcessPayment, PermissionViewPayment,
		PermissionViewItem,
		PermissionViewTable,
	},
	"inventory_clerk": {
		PermissionViewItem,
		PermissionManageStock, PermissionViewStock, PermissionAdjustStock,
	},
	"viewer": {
		PermissionViewOrder,
		PermissionViewItem,
		PermissionViewTable,
		PermissionViewReports,
	},
}

// HasPermission checks if a role has a specific permission
func (r *Role) HasPermission(permission Permission) bool {
	for _, p := range r.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// AddPermission adds a permission to a role
func (r *Role) AddPermission(permission Permission) {
	if !r.HasPermission(permission) {
		r.Permissions = append(r.Permissions, permission)
	}
}

// RemovePermission removes a permission from a role
func (r *Role) RemovePermission(permission Permission) {
	for i, p := range r.Permissions {
		if p == permission {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			break
		}
	}
}

// GetRoleByName returns a predefined role by name
func GetRoleByName(name string) (*Role, error) {
	permissions, exists := PredefinedRoles[name]
	if !exists {
		return nil, fmt.Errorf("role %s does not exist", name)
	}
	return &Role{
		Name:        name,
		Permissions: permissions,
		IsActive:    true,
	}, nil
}

// UserRoleAssignment represents the role assignment for a user
type UserRoleAssignment struct {
	ID         string `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID     string `gorm:"type:uuid;not null"`
	RoleID     string `gorm:"type:uuid;not null"`
	BusinessID string `gorm:"type:uuid;not null"`
	BranchID   string `gorm:"type:uuid"`
	IsActive   bool   `gorm:"default:true"`
}

var db *gorm.DB

// SetDB sets the database connection for the role package
func SetDB(database *gorm.DB) {
	db = database
}

// GetUserRoles returns all roles for a user in a specific business
func GetUserRoles(userID, businessID string) ([]*Role, error) {
	var userRoles []UserRoleAssignment
	if err := db.Where("user_id = ? AND business_id = ?", userID, businessID).Find(&userRoles).Error; err != nil {
		return nil, err
	}

	var roles []*Role
	for _, ur := range userRoles {
		role, err := GetRoleByName(ur.RoleID)
		if err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// HasPermission checks if a user has a specific permission in a business
func HasPermission(userID, businessID string, permission Permission) (bool, error) {
	roles, err := GetUserRoles(userID, businessID)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role.HasPermission(permission) {
			return true, nil
		}
	}
	return false, nil
}
