package main

import (
	"github.com/Kong/go-pdk/server"
	"github.com/keremdokumaci/goql-plugin/plugin"
)

func main() {
	_ = server.StartServer(plugin.New, plugin.VERSION, plugin.PRIORITY)
}
