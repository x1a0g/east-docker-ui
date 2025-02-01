package test

import (
	"fmt"
	docker "github.com/fsouza/go-dockerclient"
	"testing"
)

func CreateCon(t *testing.T) {
	// 创建 Docker 客户端
	client, err := docker.NewClient("tcp://192.168.70.129:2375")
	if err != nil {
		panic(err)
	}

	fmt.Println("创建容器")
	opts := docker.CreateContainerOptions{
		Name: "test",
		Config: &docker.Config{
			Image: "nginx",
			Cmd:   []string{"nginx"},
		},
	}
	container, err := client.CreateContainer(opts)
	if err != nil {
		panic(err)
	}
	fmt.Println("容器创建成功")
	fmt.Println(container.ID)

}
