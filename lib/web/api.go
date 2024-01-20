package web

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"tea.chunkbyte.com/kato/drive-health/lib/hardware"
	"tea.chunkbyte.com/kato/drive-health/lib/svc"
)

func setupApi(r *gin.Engine) {
	api := r.Group("/api/v1")

	api.GET("/disks", func(ctx *gin.Context) {

		olderThan := time.Now().Add(time.Minute * time.Duration(10) * -1)
		newerThan := time.Now()

		// Fetch the disk list
		disks, err := hardware.GetSystemHardDrives(svc.GetDatabaseRef(), &olderThan, &newerThan)
		if err != nil {
			ctx.Error(err)
		}

		if ctx.Request.URL.Query().Get("temp") != "" {
			for _, d := range disks {
				d.GetTemperature()
			}
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Disk List",
			"disks":   disks,
		})
	})

}
