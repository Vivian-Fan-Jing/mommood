package routers

import (
	"net/http"
	"testing"

	"github.com/Vivian-Fan-Jing/mommood/tests"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func TestLoginRouter(t *testing.T) {
	userName := "user1"
	loginRouter := func(e *core.RequestEvent) error {
		name := e.Request.PathValue("name")
		return e.HTML(http.StatusOK, "welcome "+name)
	}

	scenario := tests.ApiScenario{
		Name:           "test login router",
		Method:         http.MethodGet,
		URL:            "/login/" + userName,
		ExpectedStatus: 200,
		ExpectedContent: []string{
			userName,
		},
		RouterFunc: func(r *router.Router[*core.RequestEvent]) {
			r.GET("/login/{name}", loginRouter)
		},
	}

	scenario.Test(t)
}
