package service

import (
	API "east-docker-ui/common"
	"east-docker-ui/controller"
	"github.com/gin-gonic/gin"
)

func LogService(r *gin.Engine) *gin.Engine {

	r.GET(API.LOG_TOP5, controller.LogTop5)
	return r
}
