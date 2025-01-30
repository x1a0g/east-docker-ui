package controller

import (
	"east-docker-ui/model"
	"east-docker-ui/model/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LogTop5(ctx *gin.Context) {
	top5 := (&model.UserLog{}).Top5(database.DB)
	if len(top5) > 0 {
		ctx.JSON(http.StatusOK, top5)
		return
	}

	ctx.JSON(http.StatusOK, []model.UserLog{})
}
