package controller

import (
	API "east-docker-ui/common"
	"east-docker-ui/model"
	"east-docker-ui/model/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LogTop5(ctx *gin.Context) {
	top5 := (&model.UserLog{}).Top5(database.DB)
	if len(top5) > 0 {
		var resp API.ApiResponseObject
		resp.Success4data(top5)
		ctx.JSON(http.StatusOK, resp)
		return
	}

	var resp API.ApiResponseObject
	resp.Fail(API.FAIL.GetCode(), "获取失败")
	ctx.JSON(http.StatusOK, resp)
}
