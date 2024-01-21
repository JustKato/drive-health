package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"tea.chunkbyte.com/kato/drive-health/lib/hardware"
	"tea.chunkbyte.com/kato/drive-health/lib/svc"
)

func setupApi(r *gin.Engine) {
	api := r.Group("/api/v1")

	api.GET("/disks/:diskid/chart", func(ctx *gin.Context) {
		diskIDString := ctx.Param("diskid")
		diskId, err := strconv.Atoi(diskIDString)
		if err != nil {
			ctx.AbortWithStatusJSON(400, gin.H{
				"error":   err.Error(),
				"message": "Invalid Disk ID",
			})

			return
		}

		graphData, err := svc.GetDiskGraphImage(diskId, nil, nil)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"error":   err.Error(),
				"message": "Graph generation issue",
			})

			return
		}

		// Set the content type header
		ctx.Writer.Header().Set("Content-Type", "image/png")

		// Write the image data to the response
		ctx.Writer.WriteHeader(http.StatusOK)
		_, err = graphData.WriteTo(ctx.Writer)
		if err != nil {
			ctx.AbortWithStatusJSON(500, gin.H{
				"error":   err.Error(),
				"message": "Write error",
			})

			return
		}
	})

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
