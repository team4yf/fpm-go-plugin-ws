package main

import (
	_ "github.com/team4yf/fpm-go-plugin-ws/plugin"
	"github.com/team4yf/yf-fpm-server-go/fpm"
)

func main() {

	fpmApp := fpm.New()
	fpmApp.Init()
	fpmApp.Run()
}
