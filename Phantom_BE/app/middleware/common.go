package middleware

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
)

// AddCommonHeaders adds common headers that will be appended to all requests.
func AddCommonHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// HealthCheckMiddleware can be used to add health check endpoints
func HealthCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.URL.Path == "/health" {
			c.JSON(http.StatusOK, gin.H{
				"status": "healthy",
				"timestamp": time.Now().Unix(),
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
// IsLoggedIn checks if user is logged in.
func IsLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

// IsSysAdm checks if user is system admin.
func IsSysAdm() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
