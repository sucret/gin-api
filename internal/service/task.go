package service

import (
	"gin-api/pkg/mysql"
	"gin-api/pkg/mysql/model"

	"gorm.io/gorm"
)

type task struct {
	db *gorm.DB
}

var TaskService = &task{
	db: mysql.GetDB(),
}

func (t *task) GetAllTask() (taskList []model.CronTask) {
	t.db.Find(&taskList)
	return
}

func (t *task) List() (taskList []model.CronTask) {
	t.db.Find(&taskList)
	return
}
