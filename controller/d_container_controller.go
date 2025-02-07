package controller

import (
	API "east-docker-ui/common"
	"east-docker-ui/model/dto"
	"errors"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	containers, err := client.ListContainers(docker.ListContainersOptions{All: true})
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
		err := client.RemoveContainer(docker.RemoveContainerOptions{ID: id})
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

// 创建容器
func ContainerCreate(ctx *gin.Context) {
	var req dto.CreateConRequest
	err := ctx.ShouldBindJSON(&req)
	if err != nil || !req.Validate() {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}
	_, _, client := getDockerClient(ctx)
	//默认创建网络模式为bridge
	// 创建自定义网络（隔离环境）
	networkName := uuid.NewString() + "-network"
	if _, err := client.CreateNetwork(docker.CreateNetworkOptions{
		Name:   networkName,
		Driver: "bridge",
	}); err != nil && errors.Is(err, docker.ErrContainerAlreadyExists) {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), err.Error())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	//创建容器
	containerOpts := docker.CreateContainerOptions{Name: req.Container.Name}

	//1、配置容器的配置信息==========================
	configOpts := docker.Config{
		Image: req.Container.DockerConfig.Image,
	}

	if len(req.Container.DockerConfig.Environments) > 0 {
		//设置环境变量
		envs := req.ResolveEnvs()
		configOpts.Env = envs
	}

	if len(req.Container.DockerConfig.Cmd) > 1 {
		configOpts.Cmd = req.Container.DockerConfig.Cmd
	}

	//设置健康检查
	if req.Container.DockerConfig.HealthCheck != nil {
		configOpts.Healthcheck = &docker.HealthConfig{}
		if req.Container.DockerConfig.HealthCheck.Interval > 0 {
			configOpts.Healthcheck.Interval = req.Container.DockerConfig.HealthCheck.Interval
		}
		if req.Container.DockerConfig.HealthCheck.Retries != nil {
			configOpts.Healthcheck.Retries = *req.Container.DockerConfig.HealthCheck.Retries
		}
		if len(req.Container.DockerConfig.HealthCheck.Test) > 0 {
			configOpts.Healthcheck.Test = req.Container.DockerConfig.HealthCheck.Test
		}
		if req.Container.DockerConfig.HealthCheck.Timeout > 0 {
			configOpts.Healthcheck.Timeout = req.Container.DockerConfig.HealthCheck.Timeout
		}
		configOpts.Healthcheck = &docker.HealthConfig{}
	}
	containerOpts.Config = &configOpts

	//2、设置hostconfig==============================

	hostConfig := docker.HostConfig{}

	if req.Container.Host != nil {
		if req.Container.Host.MemoryMB != nil {
			hostConfig.Memory = req.SetMem()
		}
		if req.Container.Host.CPUPercent != nil {
			hostConfig.CPUPercent = req.SetCpu()
		}
		if req.Container.Host.DisableSwap != nil {
			hostConfig.MemorySwap = map[bool]int64{true: -1, false: -1}[*req.Container.Host.DisableSwap] //200mb交换内存
		}

		if len(req.Container.Host.Volumes) > 0 {
			hostConfig.Binds = req.ResolveBind()
		}

		if len(req.Container.Host.PortMapping) > 0 {
			hostConfig.PortBindings = req.ResolvePort()
		}

		if req.Container.Host.RestartPolicy != nil {
			restartOpts := docker.RestartPolicy{}

			if req.Container.Host.RestartPolicy.Policy != nil {
				restartOpts.Name = *req.Container.Host.RestartPolicy.Policy
			}

			if req.Container.Host.RestartPolicy.MaxRetries != nil {
				restartOpts.MaximumRetryCount = *req.Container.Host.RestartPolicy.MaxRetries
			}

			hostConfig.RestartPolicy = restartOpts
		}
	}
	//hostConfig.NetworkMode = networkName

	containerOpts.HostConfig = &hostConfig

	//3、网络别名===================================>
	if (req.Container.Networking) != nil {
		var networkOpts = docker.NetworkingConfig{}

		if len(req.Container.Networking.Aliases) > 0 {
			networkOpts.EndpointsConfig = map[string]*docker.EndpointConfig{
				networkName: {
					Aliases: req.Container.Networking.Aliases,
				},
			}
		}
		containerOpts.NetworkingConfig = &networkOpts
	}

	//4、构建容器=====================================>
	container, err := client.CreateContainer(containerOpts)
	switch {
	case errors.Is(err, docker.ErrContainerAlreadyExists):
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), err.Error())
		ctx.JSON(http.StatusOK, resp)
		return
	case err != nil:
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), err.Error())
		ctx.JSON(http.StatusOK, resp)
		return
	default:
		var resp API.ApiResponseObject
		resp.Success4data(container.ID)
		ctx.JSON(http.StatusOK, resp)
		return
	}

}

// 开始容器
func ContainerStart(ctx *gin.Context) {

	cid := ctx.Query("id")

	if cid == "" {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	_, _, client := getDockerClient(ctx)
	err := client.StartContainer(cid, nil)
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), err.Error())
		ctx.JSON(http.StatusOK, resp)
		return
	}
	var resp API.ApiResponseObject
	resp.Success4data(nil)
	ctx.JSON(http.StatusOK, resp)
}

// 停止容器
func ContainerStop(ctx *gin.Context) {

	cid := ctx.Query("id")

	if cid == "" {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	_, _, client := getDockerClient(ctx)
	err := client.StopContainer(cid, 10)
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), err.Error())
		ctx.JSON(http.StatusOK, resp)
		return
	}
	var resp API.ApiResponseObject
	resp.Success4data(nil)
	ctx.JSON(http.StatusOK, resp)
}

// 重启容器
func ContainerReStart(ctx *gin.Context) {

	cid := ctx.Query("id")

	if cid == "" {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	_, _, client := getDockerClient(ctx)
	err := client.RestartContainer(cid, 10)
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), err.Error())
		ctx.JSON(http.StatusOK, resp)
		return
	}
	var resp API.ApiResponseObject
	resp.Success4data(nil)
	ctx.JSON(http.StatusOK, resp)
}

// 暂停
func ContainerPauseOrUnpause(ctx *gin.Context) {

	cid := ctx.Query("id")
	status := ctx.Query("status")

	if cid == "" || status == "" {
		var resp API.ApiResponseObject
		resp.Fail(API.ERROR_PARAM.GetCode(), API.ERROR_PARAM.GetName())
		ctx.JSON(http.StatusOK, resp)
		return
	}

	_, _, client := getDockerClient(ctx)
	var err error
	if status == "pause" {
		err = client.PauseContainer(cid)
	} else {
		err = client.UnpauseContainer(cid)
	}
	if err != nil {
		var resp API.ApiResponseObject
		resp.Fail(API.FAIL.GetCode(), err.Error())
		ctx.JSON(http.StatusOK, resp)
		return
	}
	var resp API.ApiResponseObject
	resp.Success4data(nil)
	ctx.JSON(http.StatusOK, resp)
}

//取消暂停
