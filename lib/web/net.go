package web

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	// Initialize the Gin engine
	r := gin.Default()

	// Setup Health Pings
	setupHealth(r)
	// Setup Api
	setupApi(r)
	// Setup Frontend
	setupFrontend(r)

	return r
}
