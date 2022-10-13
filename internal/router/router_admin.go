package router

import (
	"gin-api/internal/api"
	"gin-api/internal/middleware"
	"gin-api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func setAdminRouter(r *gin.Engine, db *gorm.DB) {
	r.Use(middleware.Cors())

	// 登陆接口不需要走jwt验证
	r.POST("/admin/login", api.AdminApi.AdminLogin)
	r.POST("/admin/get-login-sms", api.AdminApi.AdminSendLoginSms)

	// 菜单接口不走权限验证
	r.GET("/admin/admin/menu", api.AdminApi.MenuList).Use(middleware.JWTAuth(service.AppGuardName), gin.Logger(), middleware.CustomRecovery())

	adminRouters := r.Group("/admin").Use(middleware.JWTAuth(service.AppGuardName), middleware.CheckAdminPermission(), gin.Logger(), middleware.CustomRecovery())
	{
		adminRouters.GET("/profile", api.AdminApi.AdminProfile)
		adminRouters.GET("/admin/list", api.AdminApi.AdminList)
		adminRouters.GET("/admin/detail", api.AdminApi.AdminDetail)
		adminRouters.POST("/admin/save", api.AdminApi.SaveAdmin)

		// 角色
		adminRouters.GET("/role/list", api.RoleApi.RoleList)
		adminRouters.POST("/role/save", api.RoleApi.SaveRole)
		adminRouters.GET("/role/detail", api.RoleApi.RoleDetail)
		adminRouters.GET("/role/node", api.RoleApi.RoleNode)
		adminRouters.POST("/role/node/save", api.RoleApi.SaveRoleNode)

		// 权限
		adminRouters.GET("/node/list", api.NodeApi.NodeList)
		adminRouters.POST("/node/save", api.NodeApi.SaveNode)
		adminRouters.GET("/node/tree", api.NodeApi.NodeTree)
		adminRouters.GET("/node/detail", api.NodeApi.NodeDetail)
		adminRouters.GET("/node/delete", api.NodeApi.DeleteNode)

		adminRouters.GET("/task/list", api.TaskApi.List)
	}

	// group := r.Group("/admin")
	// {
	// 	group.GET("/profile", api.AdminApi.Profile())
	// }
}
