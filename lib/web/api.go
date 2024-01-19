package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"tea.chunkbyte.com/kato/drive-health/lib/hardware"
)

func setupApi(r *gin.Engine) {
	api := r.Group("/v1/api")

	api.GET("/disks", func(ctx *gin.Context) {
		// Fetch the disk list
		disks, err := hardware.GetSystemHardDrives()
		if err != nil {
			ctx.Error(err)
		}

		if ctx.Request.URL.Query().Get("temp") != "" {
			for _, d := range disks {
				temp := d.GetTemperature(true)
				fmt.Printf("Disk Temp: %v", temp)
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Disk List",
			"disks":   disks,
		})
	})

}
