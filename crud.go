package main

import (
	"database/sql"
	"net/http"

	"github.com/alexlovelltroy/chassis"
	"github.com/bikeshack/dcim/pkg/components"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type ComponentService struct {
	*chassis.Microservice
	CDB components.ComponentDatabase
}

type ComponentServiceConfig struct {
	chassis.MicroserviceConfig
	parallelism int
}

func NewComponentService(cfg *ComponentServiceConfig) *ComponentService {
	service := ComponentService{
		Microservice: chassis.NewMicroservice(cfg),
	}
	return &service
}

func (s *ComponentService) CreateComponent(c *gin.Context) {
	// Create a new empty component
	component := &components.Component{}
	// Bind the JSON data to the component
	err := c.BindJSON(component) // This aborts with 400 with any error
	log.Debug("Component: ", component)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not parse JSON: " + err.Error(),
		})
		return
	}
	// Validate the component
	if err = component.Validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"error": "Could not validate component: " + err.Error(),
		})
		return
	}
	uid, err := s.CDB.InsertComponent(component)
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
	if _, err := uuid.Parse(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid UUID: " + err.Error(),
		})
		return
	}
	component, err := s.CDB.GetComponent(c.Param("id"))
	switch err {

	case nil: // We found it
		c.JSON(200, gin.H{"component": component})

	case sql.ErrNoRows: // We couldn't find it
		c.JSON(404, gin.H{
			"error": "Could not find component: " + err.Error(),
		})

	default: // Something wicked this way comes
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Component Error: " + err.Error(),
		})
	}
}

func (s *ComponentService) ReplaceComponent(c *gin.Context) {
	updateComponent := &components.Component{}
	// Bind the JSON data to the component
	err := c.BindJSON(updateComponent) // This aborts with 400 with any error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
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

	err = s.CDB.UpdateComponent(updateComponent)
	switch err {
	case nil: // It worked!
		c.JSON(200, gin.H{
			"component": updateComponent,
		})
		return
	default:
		c.JSON(404, gin.H{
			"error": "Could not find component: " + err.Error(),
		})
	}
}

func (s *ComponentService) DeleteComponent(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid UUID: " + err.Error(),
		})
		return
	}
	err = s.CDB.DeleteComponent(id.String())
	switch err {
	case nil:
		c.JSON(http.StatusOK, gin.H{
			"message": c.Param("id") + " deleted",
		})
	case sql.ErrNoRows:
		c.JSON(http.StatusNotFound,
			gin.H{
				"error": c.Param("id") + " Not Found",
			})
	default:
		c.JSON(http.StatusInternalServerError,
			gin.H{
				"error": "Error deleting " + c.Param("id") + ": " + err.Error(),
			})
	}
}
