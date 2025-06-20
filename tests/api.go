package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/router"
)

// ApiScenario defines a single api request test case/scenario.
type ApiScenario struct {
	// Name is the test name.
	Name string

	// Method is the HTTP method of the test request to use.
	Method string

	// URL is the url/path of the endpoint you want to test.
	URL string

	// Body specifies the body to send with the request.
	//
	// For example:
	//
	//	strings.NewReader(`{"title":"abc"}`)
	Body io.Reader

	// Headers specifies the headers to send with the request (e.g. "Authorization": "abc")
	Headers map[string]string

	// Delay adds a delay before checking the expectations usually
	// to ensure that all fired non-awaited go routines have finished
	Delay time.Duration

	// Timeout specifies how long to wait before cancelling the request context.
	//
	// A zero or negative value means that there will be no timeout.
	Timeout time.Duration

	// ExpectedStatus specifies the expected response HTTP status code.
	ExpectedStatus int

	// List of keywords that MUST exist in the response body.
	ExpectedContent []string

	// RouterFunc to register extra routers.
	RouterFunc func(*router.Router[*core.RequestEvent])
}

func (scenario *ApiScenario) Test(t *testing.T) {
	t.Run(scenario.normalizedName(), func(t *testing.T) {
		scenario.test(t)
	})
}

func (scenario *ApiScenario) normalizedName() string {
	var name = scenario.Name

	if name == "" {
		name = fmt.Sprintf("%s:%s", scenario.Method, scenario.URL)
	}

	return name
}

func (scenario *ApiScenario) test(t testing.TB) {
	var (
		testApp    *TestApp
		testAppErr error
	)
	testApp, testAppErr = NewTestApp()
	if testAppErr != nil {
		t.Fatalf("Failed to initialize the test app instance: %v", testAppErr)
	}

	defer testApp.Cleanup()

	baseRouter, err := apis.NewRouter(testApp)
	if err != nil {
		t.Fatal(err)
	}

	if scenario.RouterFunc != nil {
		scenario.RouterFunc(baseRouter)
	}

	// manually trigger the serve event to ensure that custom app routes and middlewares are registered
	serveEvent := new(core.ServeEvent)
	serveEvent.App = testApp
	serveEvent.Router = baseRouter

	serveErr := testApp.OnServe().Trigger(serveEvent, func(e *core.ServeEvent) error {
		// reset the event counters in case a hook was triggered from a before func (eg. db save)
		testApp.ResetEventCalls()

		// add middleware to timeout long-running requests (eg. keep-alive routes)
		e.Router.Bind(&hook.Handler[*core.RequestEvent]{
			Func: func(re *core.RequestEvent) error {
				slowTimer := time.AfterFunc(3*time.Second, func() {
					t.Logf("[WARN] Long running test %q", scenario.Name)
				})
				defer slowTimer.Stop()

				if scenario.Timeout > 0 {
					ctx, cancelFunc := context.WithTimeout(re.Request.Context(), scenario.Timeout)
					defer cancelFunc()
					re.Request = re.Request.Clone(ctx)
				}

				return re.Next()
			},
			Priority: -9999,
		})

		recorder := httptest.NewRecorder()

		req := httptest.NewRequest(scenario.Method, scenario.URL, scenario.Body)

		// set default header
		req.Header.Set("content-type", "application/json")

		// set scenario headers
		for k, v := range scenario.Headers {
			req.Header.Set(k, v)
		}

		// execute request
		mux, err := e.Router.BuildMux()
		if err != nil {
			t.Fatalf("Failed to build router mux: %v", err)
		}
		mux.ServeHTTP(recorder, req)

		res := recorder.Result()

		if res.StatusCode != scenario.ExpectedStatus {
			t.Errorf("Expected status code %d, got %d", scenario.ExpectedStatus, res.StatusCode)
		}

		if scenario.Delay > 0 {
			time.Sleep(scenario.Delay)
		}

		if len(scenario.ExpectedContent) == 0 {
			if len(recorder.Body.Bytes()) != 0 {
				t.Errorf("Expected empty body, got \n%v", recorder.Body.String())
			}
		} else {
			// normalize json response format
			buffer := new(bytes.Buffer)
			err := json.Compact(buffer, recorder.Body.Bytes())
			var normalizedBody string
			if err != nil {
				// not a json...
				normalizedBody = recorder.Body.String()
			} else {
				normalizedBody = buffer.String()
			}

			for _, item := range scenario.ExpectedContent {
				if !strings.Contains(normalizedBody, item) {
					t.Errorf("Cannot find %v in response body \n%v", item, normalizedBody)
					break
				}
			}
		}

		remainingEvents := maps.Clone(testApp.EventCalls)

		var noOtherEventsShouldRemain bool

		if noOtherEventsShouldRemain && len(remainingEvents) > 0 {
			t.Errorf("Missing expected remaining events:\n%#v\nAll triggered app events are:\n%#v", remainingEvents, testApp.EventCalls)
		}

		return nil
	})
	if serveErr != nil {
		t.Fatalf("Failed to trigger app serve hook: %v", serveErr)
	}
}
