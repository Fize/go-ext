package ginserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockController struct {
	DefaultServer
}

func (m *MockController) Name() string {
	return "mock"
}

func TestRestfulAPI_handleParameter(t *testing.T) {
	api := &RestfulAPI{
		PreParameter:  "pre",
		PostParameter: "post",
	}
	controller := &MockController{}
	api.handleParameter(controller)

	assert.Equal(t, "/pre/mock", api.path)
	assert.Equal(t, "/pre/mock/post", api.longpath)
}

func TestRestfulAPI_handleMiddlewares(t *testing.T) {
	api := &RestfulAPI{}
	controller := &MockController{}
	middlewares := api.handleMiddlewares(controller)

	assert.Nil(t, middlewares)
}
