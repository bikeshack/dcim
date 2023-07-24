package main

import (
	"database/sql"
	"net/http"

	"github.com/bikeshack/dcim/internal/postgres"
	"github.com/bikeshack/dcim/pkg/components"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (s *ComponentService) CreateComponent(c *gin.Context) {
	// Create a new empty component
	component := &components.Component{}
	// Bind the JSON data to the component
	err := c.BindJSON(component)
	log.Debug("Component: ", component)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Could not parse JSON: " + err.Error(),
		})
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
	switch err {
	case nil: // It worked!
		c.JSON(200, gin.H{
			"component": updateComponent,
		})
	default:
		c.JSON(404, gin.H{
			"error": "Could not find component: " + err.Error(),
		})

	}

}

func (s *ComponentService) DeleteComponent(c *gin.Context) {
	err := postgres.DeleteComponent(s.DB, c.Param("id"))
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
