package main

import "east-docker-ui/route"

func main() {

	r := route.Route()

	err := r.Run("0.0.0.0:8081")
	if err != nil {
		panic(err)
	}
}
