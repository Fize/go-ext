package ginserver

import (
	"net/http"
	"testing"
	"time"

	"github.com/fize/go-ext/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestInitGinServer(t *testing.T) {
	cfg := config.NewConfig()
	r := InitGinServer(cfg)
	assert.NotNil(t, r)
}

func TestRun(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	cfg := &config.ServerConfig{
		BindAddr: ":8080",
	}
	go Run(router, cfg, false)
	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:8080")
	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
