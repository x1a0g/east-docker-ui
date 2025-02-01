package service

import (
	API "east-docker-ui/common"
	"east-docker-ui/controller"
	"github.com/gin-gonic/gin"
)

func ConService(r *gin.Engine) *gin.Engine {

	r.POST(API.CON_LIST, controller.ContainerList)
	r.POST(API.CON_DEL, controller.ContainerDel)
	r.POST(API.CON_CREATE, controller.ContainerCreate)
	r.GET(API.CON_INFO, controller.ContainerInfo)

	return r
}
