package route

import (
	"east-docker-ui/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {
	r := gin.Default()

	// 配置 CORS 中间件
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true // 允许所有来源的跨域请求
	// 或者指定允许的源
	// config.AllowOrigins = []string{"https://yourdomain.com"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"}

	// 如果需要支持凭证（credentials），请设置 AllowCredentials 为 true 并明确指定 AllowOrigins
	// config.AllowCredentials = true
	// 注意：当设置了 AllowCredentials 为 true 时，AllowOrigins 不能设置为 "*"，必须明确列出允许的源

	// 添加 CORS 中间件到路由器
	r.Use(cors.New(config))

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
