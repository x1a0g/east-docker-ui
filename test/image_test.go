package test

import (
	"context"
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"os"
	"testing"
)

func TestImageMcPull(t *testing.T) {
	// 创建 Docker 客户端
	client, err := docker.NewClient("tcp://192.168.70.129:2375")
	if err != nil {
		panic(err)
	}

	//proxyAddr := "docker.1ms.run/library/"
	// 定义镜像名称和标签（例如：ubuntu:latest）
	imageName := "docker.1ms.run/library/alpine"
	imageTag := "latest"
	fullImageName := fmt.Sprintf("%s:%s", imageName, imageTag)

	// 配置拉取选项
	opts := docker.PullImageOptions{
		Repository:   imageName,
		Tag:          imageTag,
		Context:      context.Background(),
		OutputStream: os.Stdout, // 输出拉取日志到标准输出
	}

	// 认证配置（如果是私有仓库）
	auth := docker.AuthConfiguration{
		//Username:      "your_username", // 例如 Docker Hub 用户名
		//Password:      "your_password",
		//ServerAddress: "https://docker.1ms.run/", // Docker Hub 地址
	}

	// 拉取镜像
	err = client.PullImage(opts, auth)
	if err != nil {
		panic(err)
	}

	fmt.Printf("镜像 %s 拉取成功！\n", fullImageName)
}

func TestBaseMcRepo(t *testing.T) {
	// 创建 Docker 客户端
	client, err := docker.NewClient("tcp://192.168.70.129:2375")
	if err != nil {
		panic(err)
	}
	client.LoadImage(docker.LoadImageOptions{})
	ex, err := client.SearchImagesEx("docker.1ms.run/library/redis", docker.AuthConfiguration{})
	if err != nil {
		panic(err)
	}
	for _, search := range ex {
		fmt.Println(search)
	}
}
