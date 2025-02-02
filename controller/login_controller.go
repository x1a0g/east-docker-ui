package controller

import (
	API "east-docker-ui/common"
	"east-docker-ui/model"
	"east-docker-ui/model/database"
	"east-docker-ui/model/dto"
	"east-docker-ui/model/vo"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func UserLogin(ctx *gin.Context) {
	var req dto.LoginRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), API.FAIL.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	//用户名密码校验
	var user model.UserInfo
	err = (&user).CheckUserNameAndPass(database.DB, req.Username, req.Password)
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.LOGIN_FAIL.GetCode(), API.LOGIN_FAIL.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	var res vo.LoginVo
	res.Userinfo = &user
	res.Token = uuid.New().String()

	var resp API.ApiResponseObject
	resp.Success4data(res)
	ctx.JSON(http.StatusOK, resp)

}
