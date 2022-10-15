package task

import (
	"fmt"
	"gin-api/internal/service"
	"gin-api/pkg/cron"
	"gin-api/pkg/mysql/model"
	"os/exec"
	"time"
)

var (
	logChan = make(chan *model.TaskLog)
	task    cron.Server
)

func StartTask() {
	task = cron.GetCron()
	task.Start()

	fmt.Println(logChan)

	taskList := service.TaskService.GetAllTask()
	for _, val := range taskList {
		Add(val)
	}

	go taskLogListener()
}

func Update(cronTask model.Task) {
	t := makeTask(cronTask)
	task.Update(t)
}

func Add(cronTask model.Task) {
	t := makeTask(cronTask)
	task.Add(t)
}

// 监听日志执行结果
func taskLogListener() {
	for {
		logModel := <-logChan
		service.TaskLogService.UpdateLog(logModel)
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
			taskLog := service.TaskLogService.SaveLog(cronTask.TaskID)

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
