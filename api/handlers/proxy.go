package handlers

import (
	"github.com/hermesespinola/proxy-app/api/middlewares"
	"github.com/kataras/iris"
)

// HandlerRedirection should redirect traffic
func HandlerRedirection(app *iris.Application) {
	app.Get("/push", middlewares.PushNode)
	app.Get("/pop", middlewares.PopNode)
	app.Get("/read", middlewares.Read)
}
