package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bikeshack/dcim/pkg/components"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type mockPostgresComponentDatabase struct{}

func (mpcd *mockPostgresComponentDatabase) InsertComponent(component *components.Component) (string, error) {
	return uuid.Must(uuid.NewRandom()).String(), nil
}

func (mpcd *mockPostgresComponentDatabase) GetComponent(id string) (*components.Component, error) {
	return &components.Component{}, nil
}

func (mpcd *mockPostgresComponentDatabase) UpdateComponent(component *components.Component) error {
	return nil
}

func (mpcd *mockPostgresComponentDatabase) DeleteComponent(id string) error {
	return nil
}

func TestCreateComponent(t *testing.T) {
	// Create a new Gin router and set the handler
	r := gin.Default()
	db := &mockPostgresComponentDatabase{}
	componentService := &ComponentService{
		CDB: db,
	} // Assuming you have initialized ComponentService correctly
	r.POST("/create", componentService.CreateComponent)

	// Test case: Successful creation
	payload := `{"xname": "x3000b7n3", "role": "compute", "class": "river", "arch": "x86_64", "net_type": "ethernet", "flag": "ok"}`
	req, _ := http.NewRequest("POST", "/create", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	if !assert.Equal(t, http.StatusCreated, resp.Code) {
		fmt.Println(resp.Body)
	}

	// Test case: Invalid JSON payload
	req, _ = http.NewRequest("POST", "/create", strings.NewReader("invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)

	// Test case: Invalid component
	payload = `{"xname": "x3000b7n3", "role": "compute", "class": "river", "arch": "x86_64", "net_type": "token ring", "flag": "ok"}`
	req, _ = http.NewRequest("POST", "/create", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusUnprocessableEntity, resp.Code)
}

func TestReadComponentHandler(t *testing.T) {
	// Create a new Gin router and set the handler
	r := gin.Default()
	componentService := &ComponentService{} // Assuming you have initialized ComponentService correctly
	r.GET("/read/:id", componentService.ReadComponent)

	// Test case: Component found
	req, _ := http.NewRequest("GET", "/read/1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	// Additional assertions based on the response body or any other expected behavior

	// Test case: Component not found
	req, _ = http.NewRequest("GET", "/read/non_existent_id", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	// Additional assertions based on the response body or any other expected behavior
}

func TestReplaceComponentHandler(t *testing.T) {
	// Create a new Gin router and set the handler
	r := gin.Default()
	componentService := &ComponentService{} // Assuming you have initialized ComponentService correctly
	r.PUT("/replace", componentService.ReplaceComponent)

	// Test case: Successful update
	payload := `{"id": 1, "name": "Updated Component A", "description": "This is the updated component A"}`
	req, _ := http.NewRequest("PUT", "/replace", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	// Additional assertions based on the response body or any other expected behavior

	// Test case: Invalid JSON payload
	req, _ = http.NewRequest("PUT", "/replace", strings.NewReader("invalid-json"))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotAcceptable, resp.Code)
	// Additional assertions based on the response body or any other expected behavior
}

func TestDeleteComponentHandler(t *testing.T) {
	// Create a new Gin router and set the handler
	r := gin.Default()
	componentService := &ComponentService{} // Assuming you have initialized ComponentService correctly
	r.DELETE("/delete/:id", componentService.DeleteComponent)

	// Test case: Successful deletion
	req, _ := http.NewRequest("DELETE", "/delete/1", nil)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	// Additional assertions based on the response body or any other expected behavior

	// Test case: Component not found
	req, _ = http.NewRequest("DELETE", "/delete/non_existent_id", nil)
	resp = httptest.NewRecorder()
	r.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	// Additional assertions based on the response body or any other expected behavior
}
