package core

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(data map[string]string) string

var InterfaceRegistry = make(map[string]bool)
var PluginRegistry = make(map[string]map[string]map[string]HandlerFunc) // interface -> plugin -> method -> handler

func RegisterInterface(interfaceName string) {
	InterfaceRegistry[interfaceName] = true
	if _, ok := PluginRegistry[interfaceName]; !ok {
		PluginRegistry[interfaceName] = make(map[string]map[string]HandlerFunc)
	}
}

func RegisterMethod(interfaceName, pluginName, methodName string, handler HandlerFunc) {
	if !InterfaceRegistry[interfaceName] {
		fmt.Printf("Interface '%s' not registered. Cannot register plugin '%s'\n", interfaceName, pluginName)
		return
	}
	if _, ok := PluginRegistry[interfaceName][pluginName]; !ok {
		PluginRegistry[interfaceName][pluginName] = make(map[string]HandlerFunc)
	}
	PluginRegistry[interfaceName][pluginName][methodName] = handler
	fmt.Printf("Registered method '%s' for plugin '%s' under interface '%s'\n", methodName, pluginName, interfaceName)
}

func Route(w http.ResponseWriter, r *http.Request) {
	interfaceName := r.URL.Query().Get("interface")
	pluginName := r.URL.Query().Get("plugin")
	methodName := r.URL.Query().Get("method")
	message := r.URL.Query().Get("message")

	// checks Interface Exists or not
	ifacePlugins, ok := PluginRegistry[interfaceName]
	if !ok {
		http.Error(w, "Interface not found", http.StatusNotFound)
		return
	}

	// checks Plugin Exists or not
	pluginMethods, ok := ifacePlugins[pluginName]
	if !ok {
		http.Error(w, "Plugin not found", http.StatusNotFound)
		return
	}

	// checks Method Exists or not
	handler, ok := pluginMethods[methodName]
	if !ok {
		http.Error(w, "Method not found", http.StatusNotFound)
		return
	}

	data := map[string]string{"message": message}
	// Calls the method
	result := handler(data)
	fmt.Fprintln(w, result)
}
