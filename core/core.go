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
	Interfaces []struct {
		Name    string `yaml:"name"`
		Plugins []struct {
			Name    string   `yaml:"name"`
			Methods []string `yaml:"methods"`
		} `yaml:"plugins"`
	} `yaml:"interfaces"`
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

	pluginValue := reflect.ValueOf(plugin)
	pluginType := pluginValue.Type()

	if _, ok := PluginRegistry[interfaceName][pluginName]; !ok {
		PluginRegistry[interfaceName][pluginName] = make(map[string]HandlerFunc)
	}

	// Get all methods from the plugin
	for i := 0; i < pluginType.NumMethod(); i++ {
		method := pluginType.Method(i)
		methodName := method.Name

		// Create a handler function for this method
		handler := func(data map[string]string) string {
			args := []reflect.Value{pluginValue, reflect.ValueOf(data)}
			results := method.Func.Call(args)
			return results[0].String()
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

	for _, iface := range config.Interfaces {
		InterfaceRegistry[iface.Name] = true
		PluginRegistry[iface.Name] = make(map[string]map[string]HandlerFunc)

		for _, plugin := range iface.Plugins {
			pluginInstance, exists := pluginMap[plugin.Name]
			if !exists {
				return fmt.Errorf("plugin '%s' not found in plugin map", plugin.Name)
			}

			if err := RegisterPlugin(iface.Name, plugin.Name, pluginInstance); err != nil {
				return err
			}
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