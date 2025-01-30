package service

import (
	API "east-docker-ui/common"
	"east-docker-ui/controller"
	"github.com/gin-gonic/gin"
)

func RepoService(r *gin.Engine) *gin.Engine {

	r.POST(API.REPO_LIST, controller.ListRepo)
	r.POST(API.REPO_INFO, controller.RepoInfo)
	r.POST(API.REPO_DEL, controller.DeleteRepo)
	r.POST(API.REPO_CREATE, controller.CreateRepo)
	r.POST(API.REPO_UPDATE, controller.UpdateRepo)
	r.POST(API.REPO_DOWN, controller.DownList)
	return r
}
