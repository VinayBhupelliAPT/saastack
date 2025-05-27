package core

import (
	"fmt"
	"net/http"
	"os"
	"reflect"

	"gopkg.in/yaml.v3"
)

type HandlerFunc func(data map[string]string) string

var InterfaceRegistry = make(map[string]bool)
var PluginRegistry = make(map[string]map[string]map[string]HandlerFunc)

type Config struct {
	Plugins []struct {
		Name      string `yaml:"name"`
		Interface string `yaml:"interface"`
	} `yaml:"plugins"`
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %v", err)
	}

	return &config, nil
}

func RegisterPlugin(interfaceName, pluginName string, plugin interface{}) error {
	if !InterfaceRegistry[interfaceName] {
		return fmt.Errorf("interface '%s' not registered", interfaceName)
	}

	if _, ok := PluginRegistry[interfaceName][pluginName]; !ok {
		PluginRegistry[interfaceName][pluginName] = make(map[string]HandlerFunc)
	}

	// Get method values using reflection
	pluginValue := reflect.ValueOf(plugin)
	pluginType := pluginValue.Type()

	// Get all methods from the plugin
	for i := 0; i < pluginType.NumMethod(); i++ {
		method := pluginType.Method(i)
		methodName := method.Name

		// Get the method value
		methodValue := pluginValue.Method(i)

		// Create a handler function for this method
		handler := func(data map[string]string) string {
			// Convert method value to HandlerFunc type
			if methodFunc, ok := methodValue.Interface().(func(map[string]string) string); ok {
				return methodFunc(data)
			}
			return fmt.Sprintf("Error: method %s is not of type func(map[string]string) string", methodName)
		}

		PluginRegistry[interfaceName][pluginName][methodName] = handler
		fmt.Printf("Registered method '%s' for plugin '%s' under interface '%s'\n", methodName, pluginName, interfaceName)
	}

	return nil
}

func InitializeFromConfig(configPath string, pluginMap map[string]interface{}) error {
	config, err := LoadConfig(configPath)
	if err != nil {
		return err
	}

	// First pass: register all interfaces
	for _, plugin := range config.Plugins {
		InterfaceRegistry[plugin.Interface] = true
		if _, ok := PluginRegistry[plugin.Interface]; !ok {
			PluginRegistry[plugin.Interface] = make(map[string]map[string]HandlerFunc)
		}
	}

	// Second pass: register all plugins
	for _, plugin := range config.Plugins {
		pluginInstance, exists := pluginMap[plugin.Name]
		if !exists {
			return fmt.Errorf("plugin '%s' not found in plugin map", plugin.Name)
		}

		if err := RegisterPlugin(plugin.Interface, plugin.Name, pluginInstance); err != nil {
			return err
		}
	}

	return nil
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

// we do not want user to know interfaces and plugins
// email/send
// payment/charge

// profobuff (they will know this)
// try grpc client to access
// service
// it should work from grpc client
// create more endpoints
