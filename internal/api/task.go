package api

import (
	"gin-api/internal/response"
	"gin-api/internal/service"

	"github.com/gin-gonic/gin"
)

type taskApi struct {
}

var TaskApi = new(taskApi)

func (*taskApi) List(c *gin.Context) {
	list := service.TaskService.List()

	response.Success(c, list)
}
