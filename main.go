package main

import (
	"log"
	"spataro/api"
	"spataro/config"
)

func main() {
	config := config.GetConfig()
	app := &api.App{}
	app.Initialize(config)
	log.Println(" - Server listen on port 4000")
	app.Run(":4000")
}
