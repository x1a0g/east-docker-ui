package service

import (
	API "east-docker-ui/common"
	"east-docker-ui/controller"
	"github.com/gin-gonic/gin"
)

func ImageService(r *gin.Engine) *gin.Engine {

	r.POST(API.IMAGE_LIST, controller.ListDockerImage)
	r.POST(API.IMAGE_DEL, controller.DeleteImages)
	r.GET(API.IMAGE_INFO, controller.ImageInfo)
	r.GET(API.IMAGE_PULL, controller.CreateImage)
	r.GET(API.IMAGE_EXPORT, controller.ExportImages)
	r.GET(API.IMAGE_IMPORT, controller.ImportImages)
	r.GET(API.IMAGE_SEARCH, controller.SearchImages)
	return r
}
