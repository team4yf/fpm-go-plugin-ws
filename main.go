package main

import (
	"fmt"

	_ "github.com/team4yf/fpm-go-plugin-ws/plugin"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

func main() {

	fpmApp := fpm.New()
	fpmApp.Init()
	handler := func(topic string, data interface{}) {
		fmt.Printf("%s %v\n", topic, data)
	}
	fpmApp.Subscribe("#ws/connect", handler)
	fpmApp.Subscribe("#ws/close", handler)
	fpmApp.Subscribe("#ws/error", handler)
	fpmApp.Subscribe("#ws/receive", handler)

	fpmApp.Run()
}
