package service

import (
	"errors"
	"fmt"
	"gin-api/internal/request"
	"gin-api/internal/response"
	"gin-api/pkg/cron"
	"gin-api/pkg/mysql"
	"gin-api/pkg/mysql/model"
	"os/exec"
	"time"

	"gorm.io/gorm"
)

type task struct {
	db *gorm.DB
}

var TaskService = &task{
	db: mysql.GetDB(),
}

func (t *task) GetAllTask() (taskList []model.Task) {
	t.db.Where("status = ?", 1).Find(&taskList)
	return
}

func (t *task) List() (taskList []model.Task) {
	t.db.Order("task_id DESC").Find(&taskList)
	return
}

func (t *task) ChangeStatus(ct request.ChangeTaskStatus) (task model.Task, err error) {
	if err = t.db.Where("task_id = ?", ct.TaskId).First(&task).Error; err != nil {
		err = errors.New("任务不存在")
		return
	}

	if task.Status == ct.Status {
		err = errors.New("状态有误")
		return
	}

	t.db.Model(&task).Update("status", ct.Status)

	if ct.Status == 2 {
		Remove(ct.TaskId)
	} else {
		Update(task)
	}
	return
}

func (t *task) Save(task request.SaveTask) (model model.Task, err error) {
	if task.TaskID != 0 {
		if err = t.db.First(&model, task.TaskID).Error; err != nil {
			err = errors.New("任务不存在")
			return
		}
	}

	model.Name = task.Name
	model.Command = task.Command
	model.Spec = task.Spec
	model.ProcessNum = task.ProcessNum
	model.Status = task.Status

	if task.TaskID != 0 {
		err = t.db.Save(&model).Error
	} else {
		err = t.db.Create(&model).Error
	}

	if err != nil {
		return
	}

	if model.Status == 1 {
		Update(model)
	} else {
		Remove(model.TaskID)
	}

	return
}

func (t *task) Detail(taskId int32) (task model.Task, err error) {
	err = t.db.First(&task, taskId).Error
	return
}

func (t *task) Log(form request.TaskLogList) (resp response.TaskLogResponse, err error) {
	t.db.Model(&model.TaskLog{}).Where("task_id = ?", form.TaskID).Count(&resp.Total)

	t.db.Where("task_id = ?", form.TaskID).Order("task_log_id desc").Offset((form.Page - 1) * form.PageSize).Limit(form.PageSize).Find(&resp.LogList)
	return
}

// ---

var (
	logChan    = make(chan *model.TaskLog)
	taskServer cron.Server
)

func StartTask() {
	taskServer = cron.GetCron()
	taskServer.Start()

	fmt.Println(logChan)

	taskList := TaskService.GetAllTask()
	for _, val := range taskList {
		Add(val)
	}

	go taskLogListener()
}

func Update(cronTask model.Task) {
	t := makeTask(cronTask)
	taskServer.Update(t)
}

func Remove(taskId int32) {
	taskServer.Remove(taskId)
}

func Add(cronTask model.Task) {
	t := makeTask(cronTask)
	taskServer.Add(t)
}

// 监听日志执行结果
func taskLogListener() {
	for {
		logModel := <-logChan
		TaskLogService.UpdateLog(logModel)
	}
}

func makeTask(cronTask model.Task) (t cron.Task) {
	t = cron.Task{
		TaskId:     cronTask.TaskID,
		TntryId:    0,
		Spec:       cronTask.Spec,
		ProcessNum: cronTask.ProcessNum,
		Func: func() {
			// 任务执行开始时写入日志
			taskLog := TaskLogService.SaveLog(cronTask.TaskID)

			fmt.Println(cronTask.Command)
			c := exec.Command("bash", "-c", cronTask.Command)

			output, err := c.CombinedOutput()

			if err != nil {
				taskLog.Status = 3
			} else {
				taskLog.Status = 2
			}

			result := string(output)

			taskLog.Log = &result
			taskLog.EndTime = time.Now()

			logChan <- taskLog
		},
	}
	return
}
