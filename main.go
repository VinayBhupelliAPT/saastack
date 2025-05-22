package main

import (
	"net/http"
	"sample/core"
	"sample/plugins"
)

func main() {

	pluginMap := map[string]interface{}{
		"email":  &plugins.EmailPlugin{},
		"stripe": &plugins.StripePlugin{},
	}

	if err := core.InitializeFromConfig("config/plugins.yaml", pluginMap); err != nil {
		panic(err)
	}

	http.HandleFunc("/api", core.Route)
	http.ListenAndServe(":8080", nil)
}

// http://localhost:8080/api?interface=notification&plugin=email&method=Send&message=Hello
