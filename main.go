package main

import (
	"log"

	"github.com/Vivian-Fan-Jing/mommood/internal/routers"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func main() {
	app := pocketbase.New()

	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		registry := template.NewRegistry()

		se.Router.GET("/login/{name}", routers.LoginRouter(registry))

		return se.Next()
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
