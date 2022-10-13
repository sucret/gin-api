// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameAdminRole = "admin_role"

// AdminRole mapped from table <admin_role>
type AdminRole struct {
	AdminRoleID int32     `gorm:"column:admin_role_id;type:int(11);primaryKey;autoIncrement:true" json:"admin_role_id"`
	AdminID     int32     `gorm:"column:admin_id;type:int(11);not null" json:"admin_id"`                                // 用户ID
	RoleID      int32     `gorm:"column:role_id;type:int(11);not null" json:"role_id"`                                  // 角色ID
	CreatedAt   time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
}

// TableName AdminRole's table name
func (*AdminRole) TableName() string {
	return TableNameAdminRole
}
