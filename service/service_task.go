package service

import (
	"context"
	"errors"
	"fmt"
	"gin-api/pkg/cron"
	"gin-api/pkg/mysql"
	"gin-api/pkg/mysql/model"
	"gin-api/request"
	"gin-api/response"
	"os/exec"
	"sync"
	"time"

	"gorm.io/gorm"
)

type task struct {
	db             *gorm.DB
	taskServer     cron.Server
	logChan        chan *model.TaskLog
	processingTask map[int64]context.CancelFunc
	mu             sync.Mutex
}

var TaskService = &task{
	db:             mysql.GetDB(),
	taskServer:     cron.GetCron(),
	logChan:        make(chan *model.TaskLog),
	processingTask: make(map[int64]context.CancelFunc),
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
		t.Remove(ct.TaskId)
	} else {
		t.Update(task)
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
		t.Update(model)
	} else {
		t.Remove(model.TaskID)
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

func (t *task) Execute(taskId int32) (err error) {
	task := &model.Task{}
	if err = t.db.First(&task, taskId).Error; err != nil {
		err = errors.New("任务不存在")
		return
	}

	// 创建方法并执行
	ta := t.makeTask(*task)
	ta.Func()
	return
}

func (t *task) StopTask(logId int64) (err error) {
	log := model.TaskLog{}
	if err = t.db.Where("task_log_id = ?", logId).First(&log, logId).Error; err != nil {
		err = errors.New("任务不存在或已结束")
		return
	}

	if log.Status != 1 {
		err = errors.New("任务已结束")
		return
	}

	t.processingTask[logId]()
	t.mu.Lock()
	delete(t.processingTask, logId)
	t.mu.Unlock()

	log.Status = 4
	log.EndTime = model.Time(time.Now())
	t.logChan <- &log

	return
}

// --------------------

func (t *task) StartTask() {
	t.taskServer.Start()

	taskList := t.GetAllTask()
	for _, val := range taskList {
		t.Add(val)
	}

	go t.taskLogListener()
}

func (t *task) Update(cronTask model.Task) {
	ta := t.makeTask(cronTask)
	t.taskServer.Update(ta)
}

func (t *task) Remove(taskId int32) {
	t.taskServer.Remove(taskId)
}

func (t *task) Add(cronTask model.Task) {
	ta := t.makeTask(cronTask)
	t.taskServer.Add(ta)
}

// 监听日志执行结果
func (t *task) taskLogListener() {
	for {
		logModel := <-t.logChan
		TaskLogService.UpdateLog(logModel)
	}
}

func (t *task) makeTask(cronTask model.Task) (ta cron.Task) {
	ta = cron.Task{
		TaskId:     cronTask.TaskID,
		TntryId:    0,
		Spec:       cronTask.Spec,
		ProcessNum: cronTask.ProcessNum,
		Func: func() {
			// 任务执行开始时写入日志
			taskLog := TaskLogService.SaveLog(cronTask.TaskID)

			var ctx context.Context
			var cancel context.CancelFunc
			if cronTask.Timeout > 0 {
				ctx, cancel = context.WithTimeout(context.Background(), time.Duration(cronTask.Timeout)*time.Second)
			} else {
				ctx, cancel = context.WithCancel(context.Background())
			}

			t.mu.Lock()
			t.processingTask[taskLog.TaskLogID] = cancel
			t.mu.Unlock()

			forever := make(chan struct{})

			done := make(chan struct{})

			f := func(doneCh chan struct{}) {
				fmt.Println(cronTask.Command)

				c := exec.Command("bash", "-c", cronTask.Command)
				output, err := c.CombinedOutput()

				if err != nil {
					taskLog.Status = 3
				} else {
					taskLog.Status = 2
				}

				result := string(output)

				taskLog.Log = result
				taskLog.EndTime = model.Time(time.Now())

				done <- struct{}{}
				t.logChan <- taskLog
			}

			go func(ctx context.Context) {
				go f(done)

				select {
				case <-ctx.Done(): // 调用cancel方法
					forever <- struct{}{}
					return
				case <-done: // 任务执行完成
					forever <- struct{}{}
					cancel()
					return
				}
			}(ctx)

			<-forever

			t.mu.Lock()
			delete(t.processingTask, taskLog.TaskLogID)
			t.mu.Unlock()
		},
	}
	return
}
