// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameAdmin = "admin"

// Admin mapped from table <admin>
type Admin struct {
	AdminID   int32  `gorm:"column:admin_id;type:int(11);primaryKey;autoIncrement:true" json:"admin_id"`
	Username  string `gorm:"column:username;type:varchar(191);not null" json:"username"`                            // 用户名称
	Mobile    string `gorm:"column:mobile;type:varchar(191);not null" json:"mobile"`                                // 用户手机号
	Password  string `gorm:"column:password;type:varchar(191);not null" json:"password"`                            // 用户密码
	Status    int64  `gorm:"column:status;type:bigint(20) unsigned;not null" json:"status"`                         // 状态：1启用|2禁用
	CreatedAt Time   `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"` // 创建时间
	RoleList  []Role `gorm:"many2many:admin_role;foreignKey:AdminID;joinForeignKey:AdminID;joinReferences:RoleID" json:"role_list"`
}

// TableName Admin's table name
func (*Admin) TableName() string {
	return TableNameAdmin
}
