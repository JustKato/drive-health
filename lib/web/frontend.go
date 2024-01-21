package web

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/JustKato/drive-health/lib/hardware"
	"github.com/JustKato/drive-health/lib/svc"
	"github.com/gin-gonic/gin"
)

func setupFrontend(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./static")

	// Set up a route for the root URL
	r.GET("/", func(ctx *gin.Context) {
		hardDrives, err := hardware.GetSystemHardDrives(svc.GetDatabaseRef(), nil, nil)
		if err != nil {
			ctx.AbortWithStatus(500)
		}

		for _, hdd := range hardDrives {
			hdd.GetTemperature()
		}

		var olderThan, newerThan *time.Time

		if ot := ctx.Query("older"); ot != "" {
			fmt.Printf("ot = %s\n", ot)
			if otInt, err := strconv.ParseInt(ot, 10, 64); err == nil {
				otTime := time.UnixMilli(otInt)
				olderThan = &otTime
			}
		}

		if nt := ctx.Query("newer"); nt != "" {
			fmt.Printf("nt = %s\n", nt)
			if ntInt, err := strconv.ParseInt(nt, 10, 64); err == nil {
				ntTime := time.UnixMilli(ntInt)
				newerThan = &ntTime
			}
		}

		if olderThan == nil {
			genTime := time.Now().Add(time.Hour * -1)

			olderThan = &genTime
		}

		if newerThan == nil {
			genTime := time.Now()

			newerThan = &genTime
		}

		// Render the HTML template
		ctx.HTML(http.StatusOK, "index.html", gin.H{
			"drives": hardDrives,
			"older":  olderThan.UnixMilli(),
			"newer":  newerThan.UnixMilli(),
		})
	})
}
