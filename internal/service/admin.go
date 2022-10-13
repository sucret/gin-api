package service

import (
	"errors"
	"fmt"
	"gin-api/internal/request"
	"gin-api/internal/response"
	"gin-api/pkg/mysql"
	"gin-api/pkg/mysql/model"
	redis_ "gin-api/pkg/redis"
	"math/rand"
	"time"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

type adminService struct {
	db    *gorm.DB
	redis *redis.Client
}

var AdminService = &adminService{
	db:    mysql.GetDB(),
	redis: redis_.GetRedis(),
}

func (a *adminService) Save(data request.SaveAdmin) (admin model.Admin, err error) {
	admin.AdminID = data.AdminId
	admin.Username = data.Username
	admin.Mobile = data.Mobile

	if data.Password != "" {
		admin.Password = data.Password
	}

	tx := a.db.Begin()

	if data.AdminId > 0 {
		if err = a.db.Where("admin_id = ?", data.AdminId).First(&model.Admin{}).Error; err != nil {
			err = errors.New("用户不存在")
			return
		}

		a.db.Model(&model.Admin{}).Where("admin_id = ?", data.AdminId).Updates(&admin)

		// 删除用户对应的角色
		err = a.db.Where("admin_id = ?", data.AdminId).Delete(model.AdminRole{}).Error
		if err != nil {
			tx.Rollback()
			return
		}
	} else {
		err = a.db.Create(&admin).Error
	}

	// 设置角色信息
	for _, val := range data.Role {
		adminRole := model.AdminRole{
			AdminID: admin.AdminID,
			RoleID:  val,
		}

		err = a.db.Create(&adminRole).Error

		if err != nil {
			err = errors.New("保存失败")
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	return
}

func (a *adminService) Detail(adminId uint) (admin response.AdminDetail, err error) {
	a.db.Where("admin_id = ?", adminId).Find(&admin.Admin)
	a.db.Table("admin_role").Select("admin_role.role_id, role.name").Joins("left join role on admin_role.role_id = role.role_id").Where("admin_role.admin_id = ?", adminId).Scan(&admin.RoleInfo)

	// 获取所有角色
	a.db.Model(&model.Role{}).Scan(&admin.RoleList)
	return
}

func (a *adminService) List() (adminList []map[string]interface{}, err error) {
	// 获取所有用户
	a.db.Model(&model.Admin{}).Select("admin_id", "username", "mobile", "created_at").Order("admin_id DESC").Find(&adminList)

	var adminIds []int32
	for _, val := range adminList {
		id, _ := val["admin_id"].(int32)
		adminIds = append(adminIds, id)
	}

	// 获取用户角色关联信息
	var roleIds []int32
	var adminRole []model.AdminRole
	a.db.Where("admin_id IN ? ", adminIds).Find(&adminRole)
	for _, val := range adminRole {
		roleIds = append(roleIds, val.RoleID)
	}

	// 获取角色列表
	var roleList []model.Role

	a.db.Where("role_id IN ? ", roleIds).Find(&roleList)

	adminRoleList := make(map[int32][]model.Role)
	for _, val := range adminRole {
		for _, v := range roleList {
			if val.RoleID == v.RoleID {
				adminRoleList[val.AdminID] = append(adminRoleList[val.AdminID], v)
			}
		}
	}

	for key, val := range adminList {
		id, _ := val["admin_id"].(int32)
		adminList[key]["role_list"] = adminRoleList[id]
	}

	return
}

// 登陆
func (a *adminService) AdminLogin(params request.AdminLogin) (admin *model.Admin, err error) {
	if params.Captcha != "" {
		// 判断验证码是否正确
		if code, _ := a.redis.Get("smscode_" + params.Mobile).Result(); code != params.Captcha {
			err = errors.New("验证码不正确")
			return
		}
	}

	if err = a.db.Where("mobile = ?", params.Mobile).First(&admin).Error; err != nil {
		err = errors.New("用户不存在")
	}

	return
}

func (a *adminService) GetAdminById(adminId int) (admin *model.Admin, err error) {
	if err = a.db.Where("admin_id = ?", adminId).First(&admin).Error; err != nil {
		err = errors.New("用户不存在")
	}

	return
}

// 获取登陆验证码
func (a *adminService) GetLoginSms(mobile string) (code string, err error) {
	var admin model.Admin

	if err = a.db.Where("mobile = ?", mobile).First(&admin).Error; err != nil {
		err = errors.New("用户不存在")
		return
	}

	rand.Seed(time.Now().Unix())

	for i := 0; i < 4; i++ {
		code = fmt.Sprintf("%s%d", code, rand.Intn(10))
	}

	// TODO 发送验证码

	// 验证码写入redis 120秒后过期
	if err = a.redis.Set("smscode_"+mobile, code, 120*time.Second).Err(); err != nil {
		err = errors.New("验证码写入失败")
	}

	return
}

// 验证后端管理员权限
func (a *adminService) CheckAdminPermission(userId int, path string) bool {
	// 超级管理员可以访问所有接口
	if userId == 1 {
		return true
	}

	node := model.Node{}
	roleList := []model.AdminRole{}
	roleIdList := []int32{}

	a.db.Where("path = ?", path).First(&node)
	if node.NodeID == 0 {
		return false
	}

	a.db.Where("admin_id = ?", userId).Find(&roleList)
	if len(roleList) == 0 {
		return false
	}
	for _, val := range roleList {
		roleIdList = append(roleIdList, val.RoleID)
	}

	var count int64
	a.db.Model(&model.RoleNode{}).Where("role_id IN ? AND node_id = ?", roleIdList, node.NodeID).Count(&count)

	return count != 0
}
