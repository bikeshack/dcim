package main

import (
	"net/http"

	"github.com/alexlovelltroy/chassis"
	"github.com/bikeshack/dcim/internal/postgres"
	"github.com/bikeshack/dcim/pkg/components"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ComponentService struct {
	*chassis.Microservice
}

type ComponentServiceConfig struct {
	chassis.MicroserviceConfig
	parallelism int
}

func (s *ComponentService) CreateComponent(c *gin.Context) {
	// Create a new empty component
	component := &components.Component{}
	// Bind the JSON data to the component
	err := c.BindJSON(component)
	log.Debug("Component: ", component)
	if err != nil {
		c.JSON(400, gin.H{
			"message":    "Could not parse JSON: " + err.Error(),
			"statusCode": "400",
		})
		return
	}
	// Validate the component
	if err = component.Validate(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not validate component: " + err.Error(),
		})
		return
	}
	uid, err := postgres.InsertComponent(s.DB, component)
	log.Debug("UID: ", uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not add component: " + err.Error(),
		})
		return
	}
	component.Uid = uuid.MustParse(uid)
	c.JSON(http.StatusCreated, gin.H{"component": component})
}

func (s *ComponentService) ReadComponent(c *gin.Context) {
	component, err := postgres.GetComponent(s.DB, c.Param("id"))
	if err != nil {
		c.JSON(404, gin.H{
			"error": "Could not find component: " + err.Error(),
		})
		return
	}
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Could not serialize component: " + err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{"component": component})
}

func (s *ComponentService) ReplaceComponent(c *gin.Context) {
	updateComponent := &components.Component{}
	// Bind the JSON data to the component
	err := c.BindJSON(updateComponent)
	if err != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{
			"error": "Could not parse JSON: " + err.Error(),
		})
		return
	}
	// Validate the component
	if err = updateComponent.Validate(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could not validate component: " + err.Error(),
		})
		return
	}
	err = postgres.UpdateComponent(s.DB, updateComponent)
	if err != nil {
		c.JSON(404, gin.H{
			"error": "Could not find component: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"component": updateComponent,
	})
}

func (s *ComponentService) DeleteComponent(c *gin.Context) {
	c.JSON(200, gin.H{
		"message":    "value",
		"statusCode": "statusCode",
	})
}

func main() {
	cfg := &ComponentServiceConfig{
		MicroserviceConfig: chassis.DefaultMicroserviceConfig(),
		parallelism:        1,
	}
	service := ComponentService{
		Microservice: chassis.NewMicroservice(cfg),
	}
	chassis.ServeCmd.Run = func(cmd *cobra.Command, args []string) {
		service.Init() // Establish connection(s) to external services and configure the gin router
		service.AddRoute("POST", "/components", service.CreateComponent)
		service.AddRoute("GET", "/components/:id", service.ReadComponent)
		service.AddRoute("PUT", "/components/:id", service.ReplaceComponent)
		service.Serve() // Start the gin router
	}
	chassis.Execute()
}
