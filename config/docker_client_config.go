package config

import (
	docker "github.com/fsouza/go-dockerclient"
	"go.uber.org/zap"
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

func init() {
	DockerClientConfigInstance = &DockerClientConfig{
		Host: "192.168.70.129",
		Port: "2375",
		Tls:  true,
		Cert: "",
		Key:  "",
		CA:   "",
	}
}

func (d *DockerClientConfig) GetHost() string {
	if d.Tls {
		return "tcp://" + d.Host + ":" + d.Port
	}
	return "unix://" + d.Host
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
