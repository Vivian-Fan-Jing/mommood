package routers

import (
	"net/http"
	"path"
	"path/filepath"
	"runtime"

	"github.com/pocketbase/pocketbase/tools/template"

	"github.com/pocketbase/pocketbase/core"
)

func LoginRouter(registry *template.Registry) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		name := e.Request.PathValue("name")

		_, currentFile, _, _ := runtime.Caller(0)
		html, err := registry.LoadFiles(
			filepath.Join(path.Dir(currentFile), "/../../views/layout.html"),
			filepath.Join(path.Dir(currentFile), "/../../views/login.html"),
		).Render(map[string]any{
			"name": name,
		})

		if err != nil {
			return e.NotFoundError("", err)
		}

		return e.HTML(http.StatusOK, html)
	}
}
