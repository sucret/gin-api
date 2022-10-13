package request

type SaveRole struct {
	RoleID int32  `form:"role_id" json:"role_id"`
	Name   string `form:"name" json:"name" binding:"required"`
}

func (SaveRole) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"name.required": "角色名称不能为空",
	}
}

type SaveRoleNode struct {
	RoleID  int32   `form:"role_id" json:"role_id" binding:"required"`
	NodeIDs []int32 `form:"node_ids" json:"node_ids"`
}

func (SaveRoleNode) GetMessages() ValidatorMessages {
	return ValidatorMessages{
		"role_id.required": "角色ID不能为空",
	}
}
