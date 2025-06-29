package routers

import (
	"net/http"

	"github.com/pocketbase/pocketbase/tools/template"

	"github.com/pocketbase/pocketbase/core"
)

func LoginRouter(registry *template.Registry) func(*core.RequestEvent) error {
	return func(e *core.RequestEvent) error {
		name := e.Request.PathValue("name")
		html, err := registry.LoadString(loginTemplate).Render(map[string]any{
			"name": name,
		})

		if err != nil {
			return e.NotFoundError("", err)
		}

		return e.HTML(http.StatusOK, html)
	}
}

const loginTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <title>{{block "title" .}}MomMoon{{end}}</title>
</head>
<body>
    {{block "body" .}}
        <p>Hello {{.name}}</p>
        <p>Welcome back</p>
    {{end}}
    MomMoon
</body>
</html>
`
