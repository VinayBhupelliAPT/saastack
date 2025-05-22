package main

import (
	"net/http"
	"sample/core"
	_ "sample/plugins"
)

func main() {
	http.HandleFunc("/api", core.Route)
	http.ListenAndServe(":8080", nil)
}
