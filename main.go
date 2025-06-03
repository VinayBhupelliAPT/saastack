package main

import (
	"saastack/core"
	_ "saastack/interfaces/bookstore"
	_ "saastack/interfaces/notification"
	_ "saastack/interfaces/payment"
	_ "saastack/plugins"
)

func main() {
	core.InitializefromConfig("config/plugins.yaml")
	core.Start()
}
