package dto

import (
	docker "github.com/fsouza/go-dockerclient"
	"strconv"
	"time"
)

type CreateConRequest struct {
	Container ContainerConfig `json:"container"`
}

//// NetworkConfig 网络配置
//type NetworkConfig struct {
//	Name   string `json:"name"`   // 自定义网络名称（必填）
//	Driver string `json:"driver"` // 驱动类型（默认 bridge）
//}

// ContainerConfig 容器配置
type ContainerConfig struct {
	Name         string            `json:"name"` // 容器名称（必填）
	DockerConfig *DockerConfig     `json:"dockerConfig"`
	Networking   *NetworkingConfig `json:"networking"`
	Host         *HostConfig       `json:"host"`
}

type DockerConfig struct {
	Image        string            `json:"image"` // 镜像名称（必填）
	Environments map[string]string `json:"envs"`  // 环境变量
	Cmd          []string          `json:"cmd"`
	HealthCheck  *HealthCheck      `json:"healthCheck"` // 健康检查配置
}

// HostConfig 主机配置
type HostConfig struct {
	//Detach        *bool         `json:"detach"`        // 是否守护进程运行
	MemoryMB      *int           `json:"memoryMB"`      // 内存限制（MB）
	DisableSwap   *bool          `json:"disableSwap"`   // 是否禁用Swap
	CPUPercent    *int           `json:"cpuPercent"`    // CPU配额（百分比）               // 资源限制
	Volumes       []VolumeMount  `json:"binds"`         // 数据卷挂载
	PortMapping   []PortConfig   `json:"portMapping"`   // 端口映射
	RestartPolicy *RestartPolicy `json:"restartPolicy"` // 重启策略
}

// NetworkingConfig 网络配置
type NetworkingConfig struct {
	Aliases []string `json:"aliases"` // 网络别名
}

type HealthCheck struct {
	Test     []string      `json:"test"`     // 检测命令  []string{"CMD", "curl", "--fail", "http://localhost:3306 || exit 1"},
	Interval time.Duration `json:"interval"` // 检测间隔
	Timeout  time.Duration `json:"timeout"`  // 超时时间
	Retries  *int          `json:"retries"`  // 重试次数
}

type VolumeMount struct {
	HostPath      string `json:"hostPath"`      // 宿主机路径
	ContainerPath string `json:"containerPath"` // 容器路径
	Mode          string `json:"mode"`          // 挂载模式 ro/rw
}

type PortConfig struct {
	ContainerPort string `json:"containerPort"` // 容器端口（如：3306/tcp）
	HostPort      []int  `json:"hostPort"`      // 宿主机端口
}

type RestartPolicy struct {
	Policy     *string `json:"policy"`     // 策略类型
	MaxRetries *int    `json:"maxRetries"` // 最大重试次数
}

func (c *CreateConRequest) Validate() bool {
	name := c.Container.Name
	image := c.Container.DockerConfig.Image

	return name != "" && image != ""
}

// ResolveEnvs 解析环境变量
func (c *CreateConRequest) ResolveEnvs() []string {

	var res []string
	envs := c.Container.DockerConfig.Environments

	if envs != nil {
		//map转切片
		for key, value := range envs {
			res = append(res, key+"="+value)
		}
	}

	return res
}

// bind
func (c *CreateConRequest) ResolveBind() []string {
	var res []string
	var bind = c.Container.Host.Volumes

	if bind != nil {
		for _, value := range bind {
			res = append(res, value.HostPath+":"+value.ContainerPath+":"+value.Mode)
		}
	}
	return res
}

// port
func (c *CreateConRequest) ResolvePort() map[docker.Port][]docker.PortBinding {
	var res map[docker.Port][]docker.PortBinding
	var ports = c.Container.Host.PortMapping

	if ports != nil {
		for _, port := range ports {
			res[docker.Port(port.ContainerPort)] = func() []docker.PortBinding {
				var re []docker.PortBinding
				for _, hostPort := range port.HostPort {
					re = append(re, docker.PortBinding{
						HostIP:   "0.0.0.0",
						HostPort: strconv.Itoa(hostPort),
					})
				}
				return re
			}()
		}
	}
	return res
}

func (c *CreateConRequest) SetMem() int64 {
	return int64(*c.Container.Host.MemoryMB) * 1024 * 1024
}

func (c *CreateConRequest) SetCpu() int64 {
	return int64(*c.Container.Host.CPUPercent * 1000)
}
