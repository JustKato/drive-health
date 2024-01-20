package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"tea.chunkbyte.com/kato/drive-health/lib/hardware"
	"tea.chunkbyte.com/kato/drive-health/lib/svc"
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
				d.GetTemperature(true)
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Disk List",
			"disks":   disks,
		})
	})

	api.GET("/snapshots", func(ctx *gin.Context) {
		snapshots := svc.GetHardwareSnapshot()

		ctx.JSON(http.StatusOK, snapshots)
	})

}
