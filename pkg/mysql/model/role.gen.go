// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameRole = "role"

// Role mapped from table <role>
type Role struct {
	RoleID    int32     `gorm:"column:role_id;type:int(11);primaryKey;autoIncrement:true" json:"role_id"`
	Name      string    `gorm:"column:name;type:varchar(40);not null" json:"name"`                                    // 角色名称
	Status    int32     `gorm:"column:status;type:tinyint(4);not null" json:"status"`                                 // 角色状态,1|正常,2|禁用
	CreatedAt time.Time `gorm:"column:created_at;type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	AdminList []Admin   `gorm:"many2many:admin_role;foreignKey:RoleID;joinForeignKey:RoleID;joinReferences:AdminID" json:"admin_list"`
	NodeList  []Node    `gorm:"many2many:role_node;foreignKey:RoleID;joinForeignKey:RoleID;joinReferences:NodeID" json:"node_list"`
}

// TableName Role's table name
func (*Role) TableName() string {
	return TableNameRole
}