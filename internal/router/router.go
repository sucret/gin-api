package router

import (
	"gin-api/pkg/mysql"

	"github.com/gin-gonic/gin"
)

func NewHttpServer() {
	db := mysql.GetDB()

	r := gin.New()

	setAdminRouter(r, db)

	r.Run(":8082")
}
