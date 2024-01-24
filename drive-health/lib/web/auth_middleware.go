package web

import "github.com/gin-gonic/gin"

func BasicAuthMiddleware(username, password string) gin.HandlerFunc {
	authorized := gin.Accounts{
		username: password,
	}

	return gin.BasicAuth(authorized)
}
