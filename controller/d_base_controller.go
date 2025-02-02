package controller

import (
	API "east-docker-ui/common"
	"east-docker-ui/config"
	"east-docker-ui/model"
	"east-docker-ui/model/database"
	"east-docker-ui/model/vo"
	"east-docker-ui/utils"
	"encoding/json"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func IndexBaseInfo(ctx *gin.Context) {

	info, done, _ := getDockerClient(ctx)
	if done {
		return
	}

	var res vo.IndexInfoVo
	res.PtVersion = "v0.0.1"
	res.OsVersion = info.OperatingSystem
	res.Who = "eastasia@gamil.com"
	res.DockerVersion = "v" + info.ServerVersion
	res.HostName = info.Name

	//utils.PrintDockerInfo(info)
	var resp API.ApiResponseObject
	resp.Success4data(res)
	ctx.JSON(http.StatusOK, resp)

}

func getDockerClient(ctx *gin.Context) (*docker.DockerInfo, bool, *docker.Client) {
	instance := config.DockerClientConfigInstance

	if instance == nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), API.FAIL.GetName())
		ctx.JSON(http.StatusOK, resp)
		return nil, true, nil
	}

	client := instance.GetRemoteClient()

	info, err := client.Info()
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), API.FAIL.GetName())
		ctx.JSON(http.StatusOK, resp)
		return nil, true, nil
	}
	return info, false, client
}

func StaticCount(ctx *gin.Context) {
	info, done, _ := getDockerClient(ctx)
	if done {
		return
	}

	//var res vo.StaticCount
	//res.TotalContainers = info.Containers
	//res.ContainersRunning = info.ContainersRunning
	//res.ContainersPaused = info.ContainersPaused
	//res.ContainersStopped = info.ContainersStopped

	var res []vo.StaticCountChat
	res = append(res, vo.StaticCountChat{
		Count: info.ContainersRunning,
		Type:  "运行中的容器",
	})
	res = append(res, vo.StaticCountChat{
		Count: info.ContainersPaused,
		Type:  "暂停的容器",
	})
	res = append(res, vo.StaticCountChat{
		Count: info.ContainersStopped,
		Type:  "停止的容器",
	})
	//utils.PrintDockerInfo(info)
	var resp API.ApiResponseObject
	resp.Success4data(res)
	ctx.JSON(http.StatusOK, resp)

}

func ResourceUsing(ctx *gin.Context) {
	//info, done, client := getDockerClient(ctx)
	//if done {
	//	return
	//}
	//_ = info.NCPU
	//_ = info.MemTotal
	//
	//_, err := client.ListContainers(docker.ListContainersOptions{All: false})
	//if err != nil {
	//	var resp API.ApiResponseObject
	//	resp.Fail(API.FAIL.GetCode(), API.FAIL.GetName())
	//	ctx.JSON(http.StatusOK, resp)
	//	return
	//}
	//for _, item := range containers {
	//
	//}

	return

}

func StreamUpload(ctx *gin.Context) {
	reader, err := ctx.Request.MultipartReader()
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), API.FAIL.GetName())
		marshal, _ := json.Marshal(resp)
		fmt.Fprintf(ctx.Writer, string(marshal))
		return
	}

	var fileLog model.FileLog
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break // end of multipart-form data
		}
		if err != nil {
			var resp API.ApiResponseObject
			resp.Fail(API.FAIL.GetCode(), API.FAIL.GetName())
			marshal, _ := json.Marshal(resp)
			fmt.Fprintf(ctx.Writer, string(marshal))

			return
		}

		// 忽略非文件部分
		if part.FormName() != "file" {
			continue
		}

		fileName := part.FileName()
		split := strings.Split(fileName, ".")
		if len(split) < 2 {
			var resp API.ApiResponseObject
			resp.Fail(API.FAIL.GetCode(), "文件格式必须为xxxx.tar")
			marshal, _ := json.Marshal(resp)
			fmt.Fprintf(ctx.Writer, string(marshal))

			return
		}
		if split[len(split)-1] != "tar" {
			var resp API.ApiResponseObject
			resp.Fail(API.FAIL.GetCode(), "文件格式必须为xxxx.tar")
			marshal, _ := json.Marshal(resp)
			fmt.Fprintf(ctx.Writer, string(marshal))
			return
		}

		fileLog.FileName = part.FileName()

		dst := "uploads/" + part.FileName() // 保存文件到本地
		// 提取目录部分
		dirPath := filepath.Dir(dst)
		errx := os.MkdirAll(dirPath, os.ModePerm)
		if errx != nil {
			fmt.Println("Failed to create directories:", err)
			return
		}
		outStream, err := os.Create(dst)
		if err != nil {
			var resp API.ApiResponseObject
			resp.Fail(API.FAIL.GetCode(), API.FAIL.GetName())
			marshal, _ := json.Marshal(resp)
			fmt.Fprintf(ctx.Writer, string(marshal))
			return
		}
		defer outStream.Close()

		//io.copy赋值文件
		_, err = io.Copy(outStream, part)
		if err != nil {
			var resp API.ApiResponseObject
			resp.Fail(API.FAIL.GetCode(), API.FAIL.GetName())
			marshal, _ := json.Marshal(resp)
			fmt.Fprintf(ctx.Writer, string(marshal))
			return
		}

		logger, _ := zap.NewDevelopment()
		logger.Info("文件上传成功")
	}
	fileLog.FilePath = "uploads/" + fileLog.FileName
	fileLog.Id = uuid.New().String()
	fileLog.CreateAt = strconv.FormatInt(utils.GetTimestamp(), 10)
	err = fileLog.SaveFileLog(database.DB)
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), "文件数据保存失败")
		marshal, _ := json.Marshal(resp)
		fmt.Fprintf(ctx.Writer, string(marshal))
		return
	}
	var resp API.ApiResponseObject
	resp.Success4data(fileLog)
	marshal, _ := json.Marshal(resp)
	fmt.Fprintf(ctx.Writer, string(marshal))
}
