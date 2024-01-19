package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tea.chunkbyte.com/kato/drive-health/lib/hardware"
)

func setupFrontend(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Set up a route for the root URL
	r.GET("/", func(c *gin.Context) {

		hardDrives, err := hardware.GetSystemHardDrives()
		if err != nil {
			c.AbortWithStatus(500)
		}

		for _, hdd := range hardDrives {
			hdd.GetTemperature(true)
		}

		// Render the HTML template
		c.HTML(http.StatusOK, "index.html", gin.H{
			"drives": hardDrives,
		})
	})
}
