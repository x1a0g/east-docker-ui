package main

import "east-docker-ui/route"

func main() {

	r := route.Route()

	err := r.Run("0.0.0.0:8081")
	if err != nil {
		panic(err)
	}

	//dockerAddr := "tcp://192.168.70.129:2375"
	//
	//client, err := docker.NewClient(dockerAddr)
	//if err != nil {
	//	panic(err)
	//}

	//images, err := client.ListImages(docker.ListImagesOptions{All: false})
	//if err != nil {
	//	println(err.Error())
	//}
	//
	//for _, img := range images {
	//	fmt.Println("ID: ", img.ID)
	//	fmt.Println("RepoTags: ", img.RepoTags)
	//	fmt.Println("Created: ", img.Created)
	//	fmt.Println("Size: ", img.Size)
	//	fmt.Println("VirtualSize: ", img.VirtualSize)
	//	fmt.Println("ParentId: ", img.ParentID)
	//}
	//
	//fmt.Println("=========================================================")
	//
	//containers, err := client.ListContainers(docker.ListContainersOptions{All: false})
	//
	//if err != nil {
	//	println(err.Error())
	//}
	//for _, container := range containers {
	//	fmt.Println("ID: ", container.ID)
	//	fmt.Println("Names: ", container.Names)
	//	fmt.Println("Image: ", container.Image)
	//	fmt.Println("Command: ", container.Command)
	//	fmt.Println("Created: ", container.Created)
	//	fmt.Println("Ports: ", container.Ports)
	//	fmt.Println("Status: ", container.Status)
	//	fmt.Println("Labels: ", container.Labels)
	//	fmt.Println("SizeRw: ", container.SizeRw)
	//}
}
