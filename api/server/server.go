package server

import (
	"os"

	"github.com/kataras/iris"
)

// Setup create server application
func Setup() *iris.Application {
	app := iris.New()
	app.Logger().SetLevel("debug")
	return app
}

// RunServer should start server
func RunServer(app *iris.Application) {
	app.Run(iris.Addr(os.Getenv("PORT")))
}
