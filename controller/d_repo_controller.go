package controller

import (
	API "east-docker-ui/common"
	"east-docker-ui/model"
	"east-docker-ui/model/database"
	"east-docker-ui/model/dto"
	"east-docker-ui/model/vo"
	"east-docker-ui/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

// 创建
func CreateRepo(ctx *gin.Context) {

	var req dto.CreateRepoDto

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}
	//如果地址的最后一位不是/就添加
	if req.RepoAddr[len(req.RepoAddr)-1] != '/' {
		req.RepoAddr = req.RepoAddr + "/"
	}

	md := model.DockerRepo{
		RepoAddr:     req.RepoAddr,
		RepoDesc:     req.RepoDesc,
		RepoName:     req.RepoName,
		RepoPassword: req.RepoPassword,
		RepoUsername: req.RepoUsername,
		CreateAt:     utils.GetTimestamp(),
		Id:           uuid.NewString(),
	}
	addr, err2 := md.FindByAddr(database.DB, req.RepoAddr)
	if err2 != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), err2.Error())
		ctx.JSON(http.StatusOK, resp)
		return
	}
	if addr != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), "仓库地址已存在")
		ctx.JSON(http.StatusOK, resp)
		return
	}
	err := (&md).Create(database.DB)
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), err.Error())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	var resp API.ApiResponseObject
	resp.Success(API.SUCCESS.GetCode(), API.SUCCESS.GetName())
	ctx.JSON(http.StatusOK, resp)
	return
}

// 删除
func DeleteRepo(ctx *gin.Context) {

	var req dto.CommonIds

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	var res []string
	for _, id := range req.Ids {
		err := (&model.DockerRepo{}).DeleteByID(database.DB, id)
		if err != nil {
			res = append(res, id)
		}
	}

	if len(res) > 0 {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), fmt.Sprintf("部分仓库ID有误:%s", res))
		ctx.JSON(http.StatusOK, resp)
		return
	}

	var resp API.ApiResponseObject
	resp.Success(API.SUCCESS.GetCode(), API.SUCCESS.GetName())
	ctx.JSON(http.StatusOK, resp)

}

// 列表查询
func ListRepo(ctx *gin.Context) {

	var req dto.SearchRepoDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}
	total, list, err := (&model.DockerRepo{}).GetList(database.DB, req.Keyword, req.Page, req.PageSize)
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), err.Error())
		ctx.JSON(http.StatusOK, resp)
		return
	}
	var resp API.ApiResponseObject
	resp.Success4data(vo.RepoListVo{
		Total:    total,
		RepoList: list,
	})
	ctx.JSON(http.StatusOK, resp)

}

// 下拉框查询
func DownList(ctx *gin.Context) {
	list := (&model.DockerRepo{}).GetRepoLs(database.DB)

	var res []vo.RepoDownVo
	for _, repo := range list {
		res = append(res, vo.RepoDownVo{
			Id: repo.Id,
		})
	}
	var resp API.ApiResponseObject
	resp.Success4data(res)
	ctx.JSON(http.StatusOK, resp)
}

// 详情
func RepoInfo(ctx *gin.Context) {

	var req dto.CommonIds

	if err := ctx.ShouldBindJSON(&req); err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	if len(req.Ids) != 1 {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
	}

	repo, err := (&model.DockerRepo{}).GetByID(database.DB, req.Ids[0])

	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), err.Error())
		ctx.JSON(http.StatusOK, resp)
		return
	}
	var resp API.ApiResponseObject
	resp.Success4data(repo)
	ctx.JSON(http.StatusOK, resp)
	return
}

// 编辑
func UpdateRepo(ctx *gin.Context) {
	var req dto.UpdateRepoDto
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	//如果地址的最后一位不是/就添加
	if req.RepoAddr[len(req.RepoAddr)-1] != '/' {
		req.RepoAddr = req.RepoAddr + "/"
	}
	//UpdateRepoDto 转换为 DockerRepo
	repo := model.DockerRepo{
		Id:           req.Id,
		RepoAddr:     req.RepoAddr,
		RepoDesc:     req.RepoDesc,
		RepoName:     req.RepoName,
		RepoPassword: req.RepoPassword,
		RepoUsername: req.RepoUsername,
	}

	err := (&repo).Update(database.DB)
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), err.Error())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	var resp API.ApiResponseObject
	resp.Success(API.SUCCESS.GetCode(), API.SUCCESS.GetName())
	ctx.JSON(http.StatusOK, resp)
	return
}
