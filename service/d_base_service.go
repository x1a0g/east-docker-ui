package service

import (
	API "east-docker-ui/common"
	"east-docker-ui/controller"
	"github.com/gin-gonic/gin"
)

func BaseService(r *gin.Engine) *gin.Engine {

	r.GET(API.BASE_INDEX, controller.IndexBaseInfo)
	r.GET(API.BASE_STATIC, controller.StaticCount)
	r.GET(API.BASE_RESOURCE, controller.ResourceUsing)
	r.POST(API.BASE_UPLOAD, controller.StreamUpload)
	return r
}
