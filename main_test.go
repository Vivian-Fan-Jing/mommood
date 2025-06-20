package main

import (
	"testing"

	"github.com/Vivian-Fan-Jing/mommood/tests"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
)

func TestNewApp(t *testing.T) {
	var (
		testApp    *tests.TestApp
		testAppErr error
	)
	testApp, testAppErr = tests.NewTestApp()
	if testAppErr != nil {
		t.Fatalf("Failed to initialize the test app instance: %v", testAppErr)
	}

	defer testApp.Cleanup()

	baseRouter, err := apis.NewRouter(testApp)
	if err != nil {
		t.Fatal(err)
	}

	serveEvent := new(core.ServeEvent)
	serveEvent.App = testApp
	serveEvent.Router = baseRouter

	err = bindFn(serveEvent)
	if err != nil {
		t.Fatal(err)
	}

	err = newApp(bindFn)
	if err != nil {
		t.Fatal(err)
	}

	main()
}
