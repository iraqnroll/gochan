package main

import (
	"github.com/iraqnroll/gochan/cmd/api"
	"github.com/iraqnroll/gochan/cmd/front"
	"github.com/iraqnroll/gochan/config"
)

func main() {
	config.InitConfig()

	if !config.ApiEnabled() && !config.FrontendEnabled() {
		panic("All features disabled, nothing to run !")
	}

	if config.ApiEnabled() {
		gochanApi := api.Api{}
		gochanApi.Init()
		gochanApi.Run(":4000")
	}

	if config.FrontendEnabled() {
		gochanFront := front.Frontend{}
		gochanFront.Init()
		gochanFront.Run(":3000")
	}
}
