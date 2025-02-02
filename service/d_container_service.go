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
	r.GET(API.CON_START, controller.ContainerStart)
	r.GET(API.CON_STOP, controller.ContainerStop)
	r.GET(API.CON_RESTART, controller.ContainerReStart)
	r.GET(API.CON_PAUSE, controller.ContainerPauseOrUnpause)
	r.GET(API.CON_INFO, controller.ContainerInfo)

	return r
}
