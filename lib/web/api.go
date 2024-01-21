package web

import (
	"net/http"
	"strconv"
	"time"

	"github.com/JustKato/drive-health/lib/hardware"
	"github.com/JustKato/drive-health/lib/svc"
	"github.com/gin-gonic/gin"
)

func setupApi(r *gin.Engine) {
	api := r.Group("/api/v1")

	// Fetch the chart image for the disk's temperature
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

		var olderThan, newerThan *time.Time

		if ot := ctx.Query("older"); ot != "" {
			if otInt, err := strconv.ParseInt(ot, 10, 64); err == nil {
				otTime := time.UnixMilli(otInt)
				olderThan = &otTime
			}
		}

		if nt := ctx.Query("newer"); nt != "" {
			if ntInt, err := strconv.ParseInt(nt, 10, 64); err == nil {
				ntTime := time.UnixMilli(ntInt)
				newerThan = &ntTime
			}
		}

		graphData, err := svc.GetDiskGraphImage(diskId, newerThan, olderThan)
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

	// Get a list of all the disks
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
