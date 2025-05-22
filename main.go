package main

import (
	"net/http"
	"sample/core"
	"sample/plugins"
)

func main() {
	// Create a map of available plugins
	pluginMap := map[string]interface{}{
		"email":  &plugins.EmailPlugin{},
		"stripe": &plugins.StripePlugin{},
	}

	// Initialize plugins from config
	if err := core.InitializeFromConfig("config/plugins.yaml", pluginMap); err != nil {
		panic(err)
	}

	http.HandleFunc("/api", core.Route)
	http.ListenAndServe(":8080", nil)
}
