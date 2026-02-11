package main

import (
	"github.com/iraqnroll/gochan/cmd/api"
	"github.com/iraqnroll/gochan/cmd/front"
	"github.com/iraqnroll/gochan/config"
)

func main() {
	cfg := config.InitConfig()

	if !cfg.Api.Enabled && !cfg.Frontend.Enabled {
		panic("All features disabled, nothing to run !")
	}

	if cfg.Api.Enabled {
		gochanApi := api.Api{}
		gochanApi.Init(cfg)
		gochanApi.Run(":4000")
	}

	if cfg.Frontend.Enabled {
		gochanFront := front.Frontend{}
		gochanFront.Init(cfg)
		gochanFront.Run(":3000")
	}
}
