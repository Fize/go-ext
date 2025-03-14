package ginserver

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequest_Default(t *testing.T) {
	req := &Request{}
	req.Default()
	assert.Equal(t, -1, req.Limit)
	assert.Equal(t, 1, req.Page)
}

func TestRequest_HandleQueryParam(t *testing.T) {
	req := &Request{Limit: 10, Page: 1, Order: "invalid"}
	totalPages := req.HandleQueryParam(100)
	assert.Equal(t, 10, totalPages)
	assert.Equal(t, Desc, req.Order)
}
