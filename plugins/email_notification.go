package plugins

import (
	"fmt"
)

type EmailPlugin struct{}

func (e *EmailPlugin) Send(data map[string]string) string {
	msg := data["message"]
	fmt.Println("EmailPlugin sending:", msg)
	return "EmailPlugin sent: " + msg
}

func (e *EmailPlugin) Delete(data map[string]string) string {
	msg := data["message"]
	fmt.Println("EmailPlugin deleting:", msg)
	return "EmailPlugin deleted: " + msg
}

func (e *EmailPlugin) Update(data map[string]string) string {
	msg := data["message"]
	fmt.Println("EmailPlugin updating:", msg)
	return "EmailPlugin updated: " + msg
}
