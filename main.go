package main

import (
	"spataro/api"
	"spataro/config"
)

func main() {
	config := config.GetConfig()
	app := &api.App{}
	app.Initialize(config)
	app.Run(":3000")
}
