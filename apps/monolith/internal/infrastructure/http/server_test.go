package http

import (
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewServer_Creates(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	logger := zerolog.Nop()

	server := NewServer("0", engine, logger)
	assert.NotNil(t, server)
	assert.NotNil(t, server.httpServer)
}

func TestServer_StartAndShutdown(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/test", func(c *gin.Context) {
		c.String(200, "ok")
	})
	logger := zerolog.Nop()

	server := NewServer("0", engine, logger)

	errCh := make(chan error, 1)
	go func() {
		errCh <- server.Start()
	}()

	time.Sleep(100 * time.Millisecond)

	// Get the actual port
	addr := server.httpServer.Addr
	resp, err := http.Get("http://localhost" + addr + "/test")
	if err == nil {
		resp.Body.Close()
	}

	server.Shutdown()

	err = <-errCh
	assert.NoError(t, err)
}
