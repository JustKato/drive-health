package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"tea.chunkbyte.com/kato/drive-health/lib/hardware"
	"tea.chunkbyte.com/kato/drive-health/lib/svc"
)

func setupFrontend(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Set up a route for the root URL
	r.GET("/", func(c *gin.Context) {
		olderThan := time.Now().Add(time.Minute * time.Duration(10) * -1)
		newerThan := time.Now()

		hardDrives, err := hardware.GetSystemHardDrives(svc.GetDatabaseRef(), &olderThan, &newerThan)
		if err != nil {
			c.AbortWithStatus(500)
		}

		for _, hdd := range hardDrives {
			hdd.GetTemperature()
		}

		// Render the HTML template
		c.HTML(http.StatusOK, "index.html", gin.H{
			"drives": hardDrives,
		})
	})
}
