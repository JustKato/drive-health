package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func setupHealth(r *gin.Engine) {

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "Pong",
		})
	})
}
