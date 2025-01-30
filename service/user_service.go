package service

import (
	API "east-docker-ui/common"
	"east-docker-ui/controller"
	"github.com/gin-gonic/gin"
)

func UserService(r *gin.Engine) *gin.Engine {

	r.POST(API.LOGIN, controller.UserLogin)

	return r
}
