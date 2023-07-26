package main

import (
	"github.com/alexlovelltroy/chassis"
	"github.com/bikeshack/dcim/internal/postgres"

	"github.com/spf13/cobra"
)

func main() {
	cfg := &ComponentServiceConfig{
		MicroserviceConfig: chassis.DefaultMicroserviceConfig(),
		parallelism:        1,
	}
	service := NewComponentService(cfg)
	chassis.ServeCmd.Run = func(cmd *cobra.Command, args []string) {
		service.Init() // Establish connection(s) to external services and configure the gin router
		service.CDB = &postgres.PostgresComponentDatabase{DB: service.DB}
		service.AddRoute("POST", "/components", service.CreateComponent)
		service.AddRoute("GET", "/components/:id", service.ReadComponent)
		service.AddRoute("PUT", "/components/:id", service.ReplaceComponent)
		service.Serve() // Start the gin router
	}
	chassis.Execute()
}
