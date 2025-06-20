package main

import (
	"log"

	"github.com/Vivian-Fan-Jing/mommood/internal/routers"
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/template"
)

func main() {
	if err := newApp(bindFn); err != nil {
		log.Fatal(err)
	}
}

func newApp(bindFn func(se *core.ServeEvent) error) error {
	app := pocketbase.New()
	app.OnServe().BindFunc(bindFn)
	return app.Start()
}

func bindFn(se *core.ServeEvent) error {
	registry := template.NewRegistry()
	se.Router.GET("/login/{name}", routers.LoginRouter(registry))
	return se.Next()
}
