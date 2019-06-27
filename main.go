package main

import (
	"fmt"

	api "./app/api"
)

func main() {

	forum, err := api.Initialize()

	if err != nil {
		panic(err)
	}

	fmt.Println("\nServer started, listening and serve on: ", forum.RouterConfig.Host+":"+forum.RouterConfig.Port, "\n")

	err = forum.Run()

	if err != nil {
		panic(err)
	}
}
