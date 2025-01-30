package route

import (
	"east-docker-ui/service"
	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {
	r := gin.Default()

	service.UserService(r)
	service.ImageService(r)
	service.BaseService(r)
	service.LogService(r)
	service.ConService(r)
	service.RepoService(r)

	r.LoadHTMLGlob("ui/*")
	r.GET("/sse", func(ctx *gin.Context) {
		ctx.HTML(200, "b.html", gin.H{})
	})
	return r
}
