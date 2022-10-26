package response

import "gin-api/pkg/mysql/model"

type AdminList struct {
	AdminId   uint        `json:"admin_id"`
	Username  string      `json:"username"`
	Mobile    string      `json:"mobile"`
	CreatedAt string      `json:"created_at"`
	RoleList  interface{} `json:"role_list"`
}

type AdminDetail struct {
	model.Admin
	RoleInfo []struct {
		RoleId uint   `json:"role_id"`
		Name   string `json:"name"`
	} `json:"role_info"`
	RoleList []struct {
		RoleId uint   `json:"role_id"`
		Name   string `json:"name"`
	} `json:"role_list"`
}

type NodeMeta struct {
	Icon         string `json:"icon"`
	Show         bool   `json:"show"`
	Title        string `json:"title"`
	HideChildren bool   `json:"hideChildren"`
}

type AdminNode struct {
	Component string   `json:"component"`
	Id        int32    `json:"id"`
	Meta      NodeMeta `json:"meta"`
	Name      string   `json:"name"`
	ParentId  int32    `json:"parentId"`
	Redirect  string   `json:"redirect"`
	Path      string   `json:"path"`
}
