package controller

import (
	API "east-docker-ui/common"
	"east-docker-ui/config"
	"east-docker-ui/model"
	"east-docker-ui/model/database"
	"east-docker-ui/model/dto"
	"east-docker-ui/utils"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// 镜像列表
func ListDockerImage(ctx *gin.Context) {
	var req dto.SearchImageDto

	if err := ctx.ShouldBindJSON(&req); err != nil {
	}

	var instance = config.DockerClientConfigInstance

	client := instance.GetRemoteClient()

	if client == nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), API.FAIL.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	//获取镜像
	images, err := client.ListImages(docker.ListImagesOptions{All: false})
	if err != nil {
		var req API.ApiResponseObject
		req.Fail(API.FAIL.GetCode(), API.FAIL.GetName())
		ctx.JSON(http.StatusOK, req)
	}

	//封装数据
	var res []model.DockerImageInfo
	for _, item := range images {
		var image model.DockerImageInfo
		image.Id = item.ID
		image.Created = utils.Timestamp2ymd(item.Created)
		image.Size = strconv.FormatInt(item.Size/1024/1024, 10) + "MB"
		if len(item.RepoTags) > 0 {
			image.Name = item.RepoTags[0]
			image.RepoTags = item.RepoTags
		}
		res = append(res, image)
	}

	var resp API.ApiResponseObject
	resp.Success4data(res)
	ctx.JSON(http.StatusOK, resp)

}

// 拉取镜像
func CreateImage(ctx *gin.Context) {
	// 1. 初始化 SSE 流
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	var req dto.CreateImageDto

	name := ctx.Query("name")
	version := ctx.Query("version")
	repoId := ctx.Query("repoId")

	if name == "" || version == "" || repoId == "" {
		sendSSEError(ctx, "参数错误")
		return
	}

	req.Name = name
	req.RepoId = repoId
	req.Version = version

	//获取仓库信息
	repo, err2 := (&model.DockerRepo{}).GetByID(database.DB, req.RepoId)
	if err2 != nil || repo == nil {
		sendSSEError(ctx, "仓库不存在")
		return
	}

	_, _, client := getDockerClient(ctx)
	var proxy = ""
	var oldAddr = repo.RepoAddr
	if strings.Contains(repo.RepoAddr, "http://") {
		proxy = strings.Replace(repo.RepoAddr, "http://", "", -1)
	} else {
		proxy = strings.Replace(repo.RepoAddr, "https://", "", -1)
	}

	opts := docker.PullImageOptions{
		Repository:    req.Name,
		Tag:           req.Version,
		Context:       ctx.Request.Context(),
		OutputStream:  newSSEWriter(ctx), // 输出拉取日志到ctx输出
		RawJSONStream: true,              // 输出为 json 格式
	}
	auth := docker.AuthConfiguration{}

	// 如果仓库有账号密码，则添加账号密码
	if repo.RepoUsername != "" && repo.RepoPassword != "" {
		auth.Username = repo.RepoUsername
		auth.Password = repo.RepoPassword
		auth.ServerAddress = oldAddr
	} else {
		opts.Repository = proxy + opts.Repository
	}

	go func() {
		err := client.PullImage(opts, auth)
		if err != nil {
			sendSSEError(ctx, fmt.Sprintf("拉取失败: %v", err))
		} else {
			sendSSEComplete(ctx)
		}
	}()

	//发送完后，等待客户端关闭连接
	<-ctx.Done()

}

// 删除镜像
func DeleteImages(ctx *gin.Context) {
	var req *dto.CommonIds
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	if len(req.Ids) == 0 {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}
	_, _, client := getDockerClient(ctx)
	var res []string
	for _, id := range req.Ids {
		err := client.RemoveImage(id)
		if err != nil {
			res = append(res, id)
		}
	}

	if len(res) > 0 {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), fmt.Sprintf("部分镜像被引用或不存在无法删除"))
		ctx.JSON(http.StatusOK, resp)
		return
	}
	var resp API.ApiResponseObject
	resp.Success(API.SUCCESS.GetCode(), fmt.Sprintf("删除成功"))
	ctx.JSON(http.StatusOK, resp)
	return

}

// 镜像详情
func ImageInfo(ctx *gin.Context) {
	param := ctx.Param("id")
	if strings.TrimSpace(param) == "" {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	_, _, client := getDockerClient(ctx)

	image, err := client.InspectImage(param)
	if err != nil || image == nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), "镜像不存在")
		ctx.JSON(http.StatusOK, resp)
		return
	}
	var resp API.ApiResponseObject
	resp.Success4data(image)
	ctx.JSON(http.StatusOK, resp)
	return
}

// 导入镜像
func ImportImages(ctx *gin.Context) {
	// 1. 初始化 SSE 流
	ctx.Header("Content-Type", "text/event-stream")
	ctx.Header("Cache-Control", "no-cache")
	ctx.Header("Connection", "keep-alive")

	fileId := ctx.Query("fileId")
	name := ctx.Query("name")
	version := ctx.Query("version")

	if fileId == "" || name == "" || version == "" {
		sendSSEError(ctx, "参数错误")
		return
	}
	var req dto.ImportImageDto

	req.Name = name
	req.Version = version
	req.FileId = fileId

	_, _, cli := getDockerClient(ctx)

	log, err := (&model.FileLog{}).GetFileLog(database.DB, req.FileId)
	if err != nil {
		sendSSEError(ctx, fmt.Sprintf("镜像文件不存在: %v", err))
		return
	}

	//文件路径
	path := log.FilePath

	file, err := os.Open(path)

	opts := docker.ImportImageOptions{
		Repository:   fmt.Sprintf("%s:%s", req.Name, req.Version),
		Source:       "-",
		InputStream:  file,
		OutputStream: newSSEWriter(ctx),
	}
	if err := cli.ImportImage(opts); err != nil {
		sendSSEError(ctx, fmt.Sprintf("导入失败: %v", err))
	} else {
		sendSSEComplete(ctx)
	}
}

// 导出镜像
func ExportImages(ctx *gin.Context) {
	// 设置响应头
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Transfer-Encoding", "binary")

	imageId := ctx.Query("id")

	_, _, cli := getDockerClient(ctx)

	image, err := cli.InspectImage(imageId)
	if err != nil || image == nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), fmt.Sprintf("镜像不存在: %v", err))
		ctx.JSON(http.StatusOK, resp)
		return
	}
	var fileName = "file"
	if len(image.RepoTags) > 0 {
		fileName = image.RepoTags[0]
	} else {
		fileName = image.ID
	}
	ctx.Header("Content-Disposition", "attachment; filename="+fileName)

	writer := ctx.Writer
	opts := docker.ExportImageOptions{
		Name:         imageId,
		OutputStream: writer,
	}

	if err := cli.ExportImage(opts); err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), fmt.Sprintf("导出失败: %v", err))
		ctx.JSON(http.StatusOK, resp)
		return
	} else {
		writer.Flush()
	}

}

// 搜索镜像
func SearchImages(ctx *gin.Context) {

	keyword := ctx.Query("keyword")
	repoId := ctx.Query("repoId")
	if keyword == "" || repoId == "" {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	// 获取仓库信息
	repo, err := (&model.DockerRepo{}).GetByID(database.DB, repoId)
	if err != nil || repo == nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), fmt.Sprintf("仓库不存在: %v", err))
		ctx.JSON(http.StatusOK, resp)
		return
	}
	opts := docker.AuthConfiguration{}

	var proxy = ""
	if repo.RepoUsername != "" && repo.RepoPassword != "" {
		opts.Username = repo.RepoUsername
		opts.Password = repo.RepoPassword
		opts.ServerAddress = repo.RepoAddr
	}

	if strings.HasPrefix(repo.RepoAddr, "https://") {
		proxy = strings.Replace(repo.RepoAddr, "https://", "", -1)
	} else {
		proxy = strings.Replace(repo.RepoAddr, "http://", "", -1)
	}

	_, _, client := getDockerClient(ctx)
	ex, err := client.SearchImagesEx(proxy+keyword, opts)
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), fmt.Sprintf("搜索失败: %v", err))
		ctx.JSON(http.StatusOK, resp)
		return
	}
	var resp API.ApiResponseObject
	resp.Success4data(ex)
	ctx.JSON(http.StatusOK, resp)
}

// tag image
// push image
// create a image from con

// dockerfile构建
type sseWriter struct {
	ctx *gin.Context
}

// 新建一个ssewriter
func newSSEWriter(ctx *gin.Context) *sseWriter {
	return &sseWriter{ctx: ctx}
}

func (w *sseWriter) Write(p []byte) (int, error) {
	// 发送进度事件
	fmt.Fprintf(w.ctx.Writer, "data: %s\n\n", p)
	w.ctx.Writer.Flush()
	return len(p), nil
}

// 发送错误事件
func sendSSEError(c *gin.Context, msg string) {
	fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", msg)
	c.Writer.Flush()
}

// 发送完成事件
func sendSSEComplete(c *gin.Context) {
	fmt.Fprintf(c.Writer, "event: complete\ndata: ok\n\n")
	c.Writer.Flush()
}
