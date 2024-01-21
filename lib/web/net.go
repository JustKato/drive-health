package web

import (
	"github.com/gin-gonic/gin"
	"tea.chunkbyte.com/kato/drive-health/lib/config"
)

func SetupRouter() *gin.Engine {
	cfg := config.GetConfiguration()

	if !cfg.DebugMode {
		// Set gin to release
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize the Gin engine
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
