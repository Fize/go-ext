package ginserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExceptResponse(t *testing.T) {
	resp := ExceptResponse(1, "error")
	assert.Equal(t, 1, resp.State.Code)
	assert.Equal(t, "error", resp.State.Msg)
}

func TestDataResponse(t *testing.T) {
	data := "test data"
	resp := DataResponse(data)
	assert.Equal(t, 0, resp.State.Code)
	assert.Equal(t, success, resp.State.Msg)
	assert.Equal(t, data, resp.Data)
}

func TestListResponse(t *testing.T) {
	data := []string{"item1", "item2"}
	resp := ListResponse(2, data)
	assert.Equal(t, 0, resp.State.Code)
	assert.Equal(t, success, resp.State.Msg)
	assert.Equal(t, data, resp.Data.(ListData).Items)
	assert.Equal(t, 2, resp.Data.(ListData).Total)
}

func TestOkResponse(t *testing.T) {
	resp := OkResponse()
	assert.Equal(t, 0, resp.State.Code)
	assert.Equal(t, success, resp.State.Msg)
	assert.Nil(t, resp.Data)
}
