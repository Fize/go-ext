package ginserver

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var unimplementedErr = errors.New("unimplemented")

// DefaultServer is the default interface for restful api.
// You can use it to composite your own interface.
type DefaultServer struct{}

// Create is the method for router.POST
func (d *DefaultServer) Create() (gin.HandlerFunc, error) {
	return nil, unimplementedErr
}

// Delete is the method for router.DELETE
func (d *DefaultServer) Delete() (gin.HandlerFunc, error) {
	return nil, unimplementedErr
}

// Update is the method for router.PUT
func (d *DefaultServer) Update() (gin.HandlerFunc, error) {
	return nil, unimplementedErr
}

// Patch is the method for router.PATCH
func (d *DefaultServer) Patch() (gin.HandlerFunc, error) {
	return nil, unimplementedErr
}

// Get is the method for router.GET
func (d *DefaultServer) Get() (gin.HandlerFunc, error) {
	return nil, unimplementedErr
}

// List is the method for router.GET with query parameters
func (d *DefaultServer) List() (gin.HandlerFunc, error) {
	return nil, unimplementedErr
}

// Version return the restful API version, default is v1
func (d *DefaultServer) Version() string {
	return "v1"
}

// Name return the restful API name, default is empty
func (d *DefaultServer) Name() string {
	return ""
}

func (d *DefaultServer) Middlewares() []MiddlewaresObject {
	return nil
}
