package controller

import (
	API "east-docker-ui/common"
	"east-docker-ui/model/dto"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// 容器集合
func ContainerList(ctx *gin.Context) {
	var req dto.ConSearchListDto
	err := ctx.ShouldBindJSON(req)
	if err != nil {
	}

	_, _, client := getDockerClient(ctx)

	containers, err := client.ListContainers(docker.ListContainersOptions{All: false})
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), API.FAIL.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	var resp API.ApiResponseObject
	resp.Success4data(containers)
	ctx.JSON(http.StatusOK, resp)
	return
}

// 容器详情
func ContainerInfo(ctx *gin.Context) {
	param := ctx.Param("id")

	if strings.TrimSpace(param) == "" {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	_, _, client := getDockerClient(ctx)

	opts := docker.InspectContainerOptions{
		ID: param,
	}
	options, err := client.InspectContainerWithOptions(opts)

	if err != nil || options == nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), "容器不存在")
		ctx.JSON(http.StatusOK, resp)
		return
	}

	var resp API.ApiResponseObject
	resp.Success4data(options)
	ctx.JSON(http.StatusOK, resp)
	return
}

// 容器删除
func ContainerDel(ctx *gin.Context) {
	var ids dto.CommonIds
	err2 := ctx.ShouldBindJSON(&ids)

	if err2 != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}
	_, _, client := getDockerClient(ctx)

	var res []string
	for _, id := range ids.Ids {
		err := client.RemoveImage(id)
		if err != nil {
			res = append(res, id)
		}
	}

	if len(res) > 0 {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), "一些容器ID不存在")
		ctx.JSON(http.StatusOK, resp)
		return
	}
	var resp API.ApiResponseObject
	resp.Success4data(nil)
	ctx.JSON(http.StatusOK, resp)
	return
}
