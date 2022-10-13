package router

import (
	"gin-api/pkg/mysql"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewHttpServer() {
	db := mysql.GetDB()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "a")
	})

	setAdminRouter(r, db)

	r.Run(":8082")
}
