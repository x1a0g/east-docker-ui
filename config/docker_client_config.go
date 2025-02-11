package config

import (
	API "east-docker-ui/common"
	"github.com/BurntSushi/toml"
	docker "github.com/fsouza/go-dockerclient"
	"go.uber.org/zap"
	"strconv"
)

type DockerClientConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
	Tls  bool   `json:"tls"`
	Cert string `json:"cert"`
	Key  string `json:"key"`
	CA   string `json:"ca"`
}

var DockerClientConfigInstance *DockerClientConfig
var Tp bool

func init() {
	var config API.Config

	_, err := toml.DecodeFile("./conf.toml", &config)
	if err != nil {
		panic(err)
	}
	Tp = config.System.Type == "remote"
	if config.System.Type == "remote" {
		DockerClientConfigInstance = &DockerClientConfig{
			Host: config.Remote.Host,
			Port: strconv.Itoa(config.Remote.Port),
			Tls:  true,
			Cert: "",
			Key:  "",
			CA:   "",
		}
	} else {
		DockerClientConfigInstance = &DockerClientConfig{}
	}

}

func (d *DockerClientConfig) GetHost() string {
	if Tp {
		return "tcp://" + d.Host + ":" + d.Port
	}
	return "unix:///var/run/docker.sock"
}

func (d *DockerClientConfig) GetRemoteClient() *docker.Client {
	logger, err := zap.NewDevelopment()
	if err != nil {
		logger.Error("docker client init error", zap.Error(err))
		return nil
	}
	client, err := docker.NewClient(d.GetHost())
	if err != nil {
		logger.Error("docker client init error", zap.Error(err))
		return nil
	}

	err = client.Ping()
	if err != nil {
		logger.Error("docker client ping error", zap.Error(err))
		return nil
	}

	logger.Info("docker client init success")

	return client
}
