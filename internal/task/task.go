package task

import (
	"bufio"
	"fmt"
	"gin-api/internal/service"
	"gin-api/pkg/cron"
	"gin-api/pkg/mysql/model"
	"io"
	"os/exec"
	"sync"
)

var task cron.Server

func StartTask() {
	task = cron.GetCron()
	task.Start()

	taskList := service.TaskService.GetAllTask()
	for _, val := range taskList {
		Add(val)
	}
}

func Update(cronTask model.CronTask) {
	t := makeTask(cronTask)
	task.Update(t)
}

func Add(cronTask model.CronTask) {
	t := makeTask(cronTask)
	task.Add(t)
}

func makeTask(cronTask model.CronTask) (t cron.Task) {
	t = cron.Task{
		TaskId:     cronTask.CronTaskID,
		TntryId:    0,
		Spec:       cronTask.Spec,
		ProcessNum: cronTask.ProcessNum,
		Func: func() {
			c := exec.Command("bash", "-c", cronTask.Command)
			stdout, err := c.StdoutPipe()
			if err != nil {
				return
			}
			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				defer wg.Done()
				reader := bufio.NewReader(stdout)
				for {
					readString, err := reader.ReadString('\n')
					if err != nil || err == io.EOF {
						return
					}
					fmt.Print(readString)
				}
			}()
			c.Start()
			wg.Wait()
		},
	}
	return
}
