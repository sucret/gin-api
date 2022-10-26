package api

import (
	"gin-api/request"
	"gin-api/response"
	"gin-api/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type roleApi struct{}

var RoleApi = new(roleApi)

func (*roleApi) RoleList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))

	roleList := service.RoleService.List(page)

	response.Success(c, roleList)
}

func (*roleApi) SaveRole(c *gin.Context) {
	var form request.SaveRole
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	role, err := service.RoleService.Save(form)

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, role)
}

func (*roleApi) RoleDetail(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Query("role_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}
	role, err := service.RoleService.Detail(uint(roleId))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, role)
}

func (*roleApi) RoleNode(c *gin.Context) {
	roleId, err := strconv.Atoi(c.Query("role_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	nodeList, err := service.RoleService.GetRoleNode(int32(roleId))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, nodeList)
}

func (*roleApi) SaveRoleNode(c *gin.Context) {
	var form request.SaveRoleNode

	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	err := service.RoleService.SaveNode(form.RoleID, form.NodeIDs)

	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}
	response.Success(c, "")
}
