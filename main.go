package main

import (
	"github.com/hermesespinola/proxy-app/api/handlers"
	"github.com/hermesespinola/proxy-app/api/server"
	"github.com/hermesespinola/proxy-app/utils"
)

func init() {
	utils.LoadEnv()
}

func main() {
	app := server.Setup()
	handlers.HandlerRedirection(app)
	server.RunServer(app)
}
