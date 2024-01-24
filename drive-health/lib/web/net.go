package web

import (
	"github.com/JustKato/drive-health/lib/config"
	"github.com/gin-gonic/gin"
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
