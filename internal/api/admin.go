package api

import (
	"fmt"
	"gin-api/internal/request"
	"gin-api/internal/response"
	"gin-api/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type adminApi struct{}

var AdminApi = new(adminApi)

func (*adminApi) Profile() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("aaaefawe")
	}
}

func (*adminApi) SaveAdmin(c *gin.Context) {
	var form request.SaveAdmin
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	admin, err := service.AdminService.Save(form)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, admin)
}

func (*adminApi) AdminDetail(c *gin.Context) {
	adminId, err := strconv.Atoi(c.Query("admin_id"))
	if err != nil {
		response.BusinessFail(c, "参数错误")
		return
	}

	admin, err := service.AdminService.Detail(uint(adminId))
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, admin)
}

// 登陆
func (*adminApi) AdminLogin(c *gin.Context) {
	var form request.AdminLogin
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	if admin, err := service.AdminService.AdminLogin(form); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		tokenData, _, err := service.JwtService.CreateToken(service.AppGuardName, admin)

		if err != nil {
			response.BusinessFail(c, err.Error())
			return
		}

		response.Success(c, tokenData)
	}
}

// 用户首页
func (*adminApi) AdminProfile(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))

	if user, err := service.AdminService.GetAdminById(userId); err != nil {
		response.BusinessFail(c, err.Error())
	} else {
		fmt.Println(user)

		response.Success(c, user)
	}
}

// 获取验证码
func (*adminApi) AdminSendLoginSms(c *gin.Context) {
	var form request.AdminSendLoginSms
	if err := c.ShouldBindJSON(&form); err != nil {
		response.ValidateFail(c, request.GetErrorMsg(form, err))
		return
	}

	code, err := service.AdminService.GetLoginSms(form.Mobile)
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	fmt.Println(code)
	// global.Log.Info(fmt.Sprintf("手机号:%s 验证码:%s", form.Mobile, code))

	response.Success(c, "")
}

func (*adminApi) AdminList(c *gin.Context) {
	adminList, err := service.AdminService.List()
	if err != nil {
		response.BusinessFail(c, err.Error())
		return
	}

	response.Success(c, adminList)
}

// 这个接口只走登陆验证，不走权限验证
func (*adminApi) MenuList(c *gin.Context) {
	userId, _ := strconv.Atoi(c.GetString("userId"))

	nodeList := service.NodeService.UserNode(uint(userId))

	response.Success(c, nodeList)
}
