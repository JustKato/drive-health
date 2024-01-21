package web

import (
	"github.com/gin-gonic/gin"
	"tea.chunkbyte.com/kato/drive-health/lib/config"
)

func SetupRouter() *gin.Engine {
	// Initialize the Gin engine
	cfg := config.GetConfiguration()
	r := gin.Default()

	r.Use(BasicAuthMiddleware(cfg.IdentityUsername, cfg.IdentityPassword))

	// Setup Health Pings
	setupHealth(r)
	// Setup Api
	setupApi(r)
	// Setup Frontend
	setupFrontend(r)

	return r
}
